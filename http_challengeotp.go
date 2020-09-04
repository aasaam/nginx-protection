package main

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/gofiber/fiber"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// RouterOTP otp challenge
func RouterOTP(
	config *Config,
	router fiber.Router,
	rsa *rsa.PrivateKey,
) {

	router.Patch("/challenge/otp", func(c *fiber.Ctx) {
		challengeType := ChallengeTypeOTP

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

		otpConfig, err := GetOTPSecret(c)
		if err != nil {
			return
		}

		passcode, err := totp.GenerateCodeCustom(otpConfig.Secret, time.Now(), totp.ValidateOpts{
			Period:    otpConfig.Period,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA512,
		})

		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("error")
			protectionStatus.SetType("cannot_generate_otp_passcode")
			protectionStatus.AddToResponse(c)

			message := "Cannot generate OTP passcode"
			code := 500
			c.Status(code)
			c.Send(message)
			return
		}

		if challengeToken.Content == passcode {
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
}
