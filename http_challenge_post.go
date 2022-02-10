package main

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func httpChallengePost(c *fiber.Ctx, config *config, challengeStorage *challengeStorage) error {
	ip := c.Locals(localVarIP).(string)
	requestID := c.Locals(localVarRequestID).(string)
	temporaryChecksum := c.Locals(localVarClientTemporaryChecksum).(string)

	var chReq challengeRequest
	if parseError := c.BodyParser(&chReq); parseError != nil {
		errorMessage := "invalid request"

		defer config.getLogger().
			Warn().
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyIP, ip).
			Str(logPropertyError, parseError.Error()).
			Str(logPropertyRequestID, requestID).
			Msg(errorMessage)

		return errors.New(errorMessage)
	}

	challenge, challengeParseErr := newChallengeFromString(chReq.ChallengeToken, config.clientSecret)
	if challengeParseErr != nil {
		errorMessage := "invalid challenge token for parse"

		defer config.getLogger().
			Warn().
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyIP, ip).
			Str(logPropertyError, challengeParseErr.Error()).
			Str(logPropertyRequestID, requestID).
			Msg(errorMessage)

		return errors.New(errorMessage)
	}

	if challengeStorage.exist(challenge.ID) {
		errorMessage := "duplicate try for solve"

		defer config.getLogger().
			Warn().
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyIP, ip).
			Str(logPropertyRequestID, requestID).
			Msg(errorMessage)

		return errors.New(errorMessage)
	}

	if !challenge.verify(temporaryChecksum) {
		errorMessage := "token invalid, timeout or expired"

		defer config.getLogger().
			Info().
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyIP, ip).
			Str(logPropertyRequestID, requestID).
			Msg(errorMessage)

		return errors.New(errorMessage)
	}

	challengeStorage.set(challenge.ID, challenge.TTL)

	valid := false

	var ldapCookie *fiber.Cookie = nil
	userName := ""

	switch challenge.ChallengeType {
	case challengeTypeJS:
		valid = challenge.verifyJSValue(chReq.JSValue)
	case challengeTypeCaptcha:
		valid = challenge.verifyJSValue(chReq.JSValue) && challenge.verifyCaptchaValue(chReq.CaptchaValue)
	case challengeTypeLDAP:
		preValidation := challenge.verifyJSValue(chReq.JSValue) && challenge.verifyCaptchaValue(chReq.CaptchaValue)
		if preValidation {
			ldapURL := c.Locals(localVarLDAPURL).(string)
			ldapReadonlyUsername := c.Locals(localVarLDAPReadonlyUsername).(string)
			ldapReadonlyPassword := c.Locals(localVarLDAPReadonlyPassword).(string)
			ldapBaseDN := c.Locals(localVarLDAPBaseDN).(string)
			ldapFilter := c.Locals(localVarLDAPFilter).(string)
			ldapAttributes := c.Locals(localVarLDAPAttributes).(string)

			success, loginErr, ldapConfigError := ldapLogin(
				ldapURL,
				ldapReadonlyUsername,
				ldapReadonlyPassword,
				ldapBaseDN,
				ldapFilter,
				ldapAttributes,
				chReq.LDAPUsername,
				chReq.LDAPPassword,
			)

			if ldapConfigError != nil {
				defer config.getLogger().
					Error().
					Str(logType, logTypeLDAPError).
					Str(logPropertyError, ldapConfigError.Error()).
					Str(logPropertyChallengeType, challenge.ChallengeType).
					Str(logPropertyIP, ip).
					Str(logPropertyRequestID, requestID).
					Send()
			} else if loginErr != nil {
				defer config.getLogger().
					Warn().
					Str(logType, logTypeLDAPError).
					Str(logPropertyError, loginErr.Error()).
					Str(logPropertyChallengeType, challenge.ChallengeType).
					Str(logPropertyIP, ip).
					Str(logPropertyRequestID, requestID).
					Send()
			} else {
				valid = success
				userName = chReq.LDAPUsername
				ldapCookie = &fiber.Cookie{
					Name:     defaultUsernameCookieName,
					Value:    chReq.LDAPUsername,
					HTTPOnly: true,
					Path:     "/",
					Expires:  time.Now().Add(time.Second * time.Duration(challenge.TTL)),
				}
			}
		}
	case challengeTypeTOTP:
		valid = challenge.verifyJSValue(chReq.JSValue) && challenge.verifyTOTP(chReq.TOTPPassword)
	}

	if valid {
		defer config.getLogger().
			Info().
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logType, logTypeChallengeSuccess).
			Str(logPropertyIP, ip).
			Str(logPropertyRequestID, requestID).
			Send()

		defer prometheusRequestChallengeSuccess.WithLabelValues(challenge.ChallengeType).Inc()

		persistToken := newPersistToken(challenge.ChallengeType, challenge.ClientPersistChecksum, userName, challenge.TTL)
		tokenString := persistToken.generate(config.tokenSecret)

		tokenCookie := &fiber.Cookie{
			Name:     c.Get(httpRequestHeaderConfigCookie, defaultCookieName),
			Value:    tokenString,
			HTTPOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Second * time.Duration(challenge.TTL)),
		}

		c.Set(httpResponseChallengeResult, tokenString)

		c.Cookie(tokenCookie)
		if ldapCookie != nil {
			c.Cookie(ldapCookie)
		}
		return c.JSON(tokenString)
	}

	defer prometheusRequestChallengeFailed.WithLabelValues(challenge.ChallengeType).Inc()

	defer config.getLogger().
		Warn().
		Str(logType, logTypeChallengeFailed).
		Str(logPropertyChallengeType, challenge.ChallengeType).
		Str(logPropertyIP, ip).
		Str(logPropertyRequestID, requestID).
		Send()

	c.SendStatus(403)
	return c.JSON(false)
}
