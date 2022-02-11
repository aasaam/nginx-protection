package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func httpChallengePatch(c *fiber.Ctx, config *config) error {
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
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logPropertyIP, ip).
			Str(logPropertyError, challengeParseErr.Error()).
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

	chResp := challengeResponse{}

	switch challenge.ChallengeType {
	case challengeTypeJS:
		chResp.JSProblem = challenge.setJSValue()
	case challengeTypeCaptcha:
		image, err := challenge.setCaptchaValue(config.restCaptchaURL, c.Get(httpRequestHeaderConfigCaptchaDifficulty, "hard"))
		if err != nil {
			errorMessage := "cannot get image from captcha server"

			defer config.getLogger().
				Error().
				Str(logPropertyChallengeType, challenge.ChallengeType).
				Str(logType, logTypeChallengeFailed).
				Str(logPropertyIP, ip).
				Str(logPropertyRequestID, requestID).
				Msg(errorMessage)

			return errors.New(errorMessage)
		}
		chResp.JSProblem = challenge.setJSValue()
		chResp.CaptchaProblem = image
	case challengeTypeLDAP:
		image, err := challenge.setCaptchaValue(config.restCaptchaURL, c.Get(httpRequestHeaderConfigCaptchaDifficulty, "hard"))
		if err != nil {
			errorMessage := "cannot get image from captcha server"

			defer config.getLogger().
				Error().
				Str(logPropertyChallengeType, challenge.ChallengeType).
				Str(logType, logTypeChallengeFailed).
				Str(logPropertyIP, ip).
				Str(logPropertyRequestID, requestID).
				Msg(errorMessage)

			return errors.New(errorMessage)
		}
		chResp.JSProblem = challenge.setJSValue()
		chResp.CaptchaProblem = image
	case challengeTypeTOTP:
		chResp.JSProblem = challenge.setJSValue()
	}

	var challengeErr error
	chResp.ChallengeToken, challengeErr = challenge.getChallengeToken(config.clientSecret)
	if challengeErr != nil {
		errorMessage := "cannot generate updated challenge"

		defer config.getLogger().
			Warn().
			Str(logPropertyChallengeType, challenge.ChallengeType).
			Str(logType, logTypeChallengeFailed).
			Str(logPropertyIP, ip).
			Str(logPropertyError, challengeErr.Error()).
			Str(logPropertyRequestID, requestID).
			Msg(errorMessage)

		return errors.New(errorMessage)
	}

	c.Set(httpResponseChallengeTemporary, chResp.ChallengeToken)
	return c.JSON(chResp)
}
