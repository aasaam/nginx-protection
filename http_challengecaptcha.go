package main

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/gofiber/fiber"
)

// RouterCaptcha captcha challenge
func RouterCaptcha(
	config *Config,
	router fiber.Router,
	rsa *rsa.PrivateKey,
) {

	router.Patch("/challenge/captcha", func(c *fiber.Ctx) {
		challengeType := ChallengeTypeCaptcha

		ttlSeconds, err := GetChallengeTTL(c)
		if err != nil {
			return
		}

		challengeToken := new(ChallengeToken)
		if err := c.BodyParser(challengeToken); err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("invalid_client_challenge_json")
			protectionStatus.AddToResponse(c)

			message := "Invalid client challenge json"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		jsonChallenge, err := VerifyToken(challengeToken.Token, config.Salt, rsa)
		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("invalid_challenge_json")
			protectionStatus.AddToResponse(c)

			message := "Invalid challenge token"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		challengeChecksum := GetClientTokenCheckSum(c)

		challenge, err := ChallengeFromJSON(jsonChallenge)

		if err != nil || challenge.Validate(challengeType, challengeChecksum) == false {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("challenge_token_not_valid_or_expired")
			protectionStatus.AddToResponse(c)

			message := "Invalid challenge token or it's expired"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		success := challenge.Result == challengeToken.Content

		if success {
			cookieName, err := GetCookieName(c)
			if err != nil {
				return
			}

			clientCheckSum := GetClientCheckSum(c)

			authToken := AuthToken{}
			authToken.Checksum = clientCheckSum
			authToken.Type = challenge.Type

			tokenValue, err := GenerateAuthToken(authToken, ttlSeconds, config.Salt, rsa)

			if err != nil {
				protectionStatus := NewProtectionStatus()
				protectionStatus.SetStatus("error")
				protectionStatus.SetType("cannot_generate_token")
				protectionStatus.AddToResponse(c)

				message := "Cannot generate token"
				code := 500
				c.Status(code)
				c.Send(message)
				return
			}

			Expiry := time.Now().UTC()
			Expiry = Expiry.Add(time.Second * time.Duration(ttlSeconds))

			c.Cookie(&fiber.Cookie{
				Name:     cookieName,
				Value:    tokenValue,
				Expires:  Expiry,
				Path:     "/",
				HTTPOnly: true,
			})
			code := 200
			c.Status(code).Send(http.StatusText(code))
			return
		}

		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("client_error")
		protectionStatus.SetType("challenge_is_not_valid")
		protectionStatus.AddToResponse(c)

		message := "Challenge is not valid"
		code := 400
		c.Status(code)
		c.Send(message)
	})

	router.Post("/challenge/captcha", func(c *fiber.Ctx) {
		challengeType := ChallengeTypeCaptcha

		ttlSeconds, err := GetChallengeTTL(c)
		if err != nil {
			return
		}

		challengeToken := new(ChallengeToken)
		if err := c.BodyParser(challengeToken); err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("invalid_client_challenge_json")
			protectionStatus.AddToResponse(c)

			message := "Invalid client challenge json"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		jsonChallenge, err := VerifyToken(challengeToken.Token, config.Salt, rsa)
		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("invalid_challenge_token")
			protectionStatus.AddToResponse(c)

			message := "Invalid challenge token"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		challengeChecksum := GetClientTokenCheckSum(c)

		challenge, err := ChallengeFromJSON(jsonChallenge)

		if err != nil || challenge.Validate(challengeType, challengeChecksum) == false {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("client_error")
			protectionStatus.SetType("challenge_token_not_valid_or_expired")
			protectionStatus.AddToResponse(c)

			message := "Invalid challenge token or it's expired"
			code := 400
			c.Status(code)
			c.Send(message)
			return
		}

		isFarsiCaptcha, err := GetFarsiCaptcha(c)
		if err != nil {
			return
		}
		challenge.GenerateCaptchaChallenge(isFarsiCaptcha)

		updatedChallengeToken, err := GenerateToken(challenge.JSON(), 0, ttlSeconds, config.Salt, rsa)

		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("error")
			protectionStatus.SetType("cannot_update_token")
			protectionStatus.AddToResponse(c)

			message := "Cannot update the challenge"
			code := 500
			c.Status(code)
			c.Send(message)
			return
		}

		newClientToken := ChallengeToken{}
		newClientToken.Token = updatedChallengeToken
		newClientToken.Type = challenge.Type
		newClientToken.Content = challenge.Content

		if config.Testing {
			newClientToken.Value = challenge.Result
		}

		if err := c.JSON(newClientToken); err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("error")
			protectionStatus.SetType("cannot_json_token")
			protectionStatus.AddToResponse(c)

			message := "Cannot json the token and content"
			code := 500
			c.Status(code)
			c.Send(message)
			return
		}
	})
}
