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

	switch challenge.ChallengeType {
	case challengeTypeJS:
		valid = challenge.verifyJSValue(chReq.JSValue)
	case challengeTypeCaptcha:
		valid = challenge.verifyJSValue(chReq.JSValue) && challenge.verifyCaptchaValue(chReq.CaptchaValue)
	case challengeTypeTOTP:
		valid = challenge.verifyJSValue(chReq.JSValue) && challenge.verifyTOTP(chReq.TOTPCode)
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

		persistToken := newPersistToken(challenge.ChallengeType, challenge.ClientPersistChecksum, challenge.TTL)
		tokenString := persistToken.generate(config.tokenSecret)

		cookie := new(fiber.Cookie)
		cookie.Name = c.Get(httpRequestHeaderConfigCookie, defaultCookieName)
		cookie.Value = tokenString
		cookie.HTTPOnly = true
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(time.Second * time.Duration(challenge.TTL))
		c.Set(httpResponseChallengeResult, tokenString)
		c.Cookie(cookie)
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
