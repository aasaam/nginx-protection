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
		defer config.getLogger().Error().Str("ip", ip).Str("rid", requestID).Str("err", parseError.Error()).Msg(errorMessage)
		return errors.New(errorMessage)
	}

	challenge, challengeParseErr := newChallengeFromString(chReq.ChallengeToken, config.clientSecret)
	if challengeParseErr != nil {
		errorMessage := "invalid challenge token for parse"
		defer config.getLogger().Warn().Str("ip", ip).Str("rid", requestID).Str("err", challengeParseErr.Error()).Msg(errorMessage)
		return errors.New(errorMessage)
	}

	if !challenge.verify(temporaryChecksum) {
		errorMessage := "token invalid, timeout or expired"
		config.getLogger().Warn().Str("ip", ip).Str("rid", requestID).Msg(errorMessage)
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
			config.getLogger().Error().Str("ip", ip).Str("rid", requestID).Str("err", err.Error()).Msg(errorMessage)
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
		config.getLogger().Error().Str("ip", ip).Str("rid", requestID).Str("err", challengeErr.Error()).Msg(errorMessage)
		return errors.New(errorMessage)
	}

	c.Set(httpResponseChallengeToken, chResp.ChallengeToken)

	return c.JSON(chResp)
}
