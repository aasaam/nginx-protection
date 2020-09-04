package main

import (
	"crypto/rsa"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
)

// RouterChallenge generate challenge for solve
func RouterChallenge(
	config *Config,
	router fiber.Router,
	rsa *rsa.PrivateKey,
) {

	router.Get("/challenge", func(c *fiber.Ctx) {
		challengeType, err := GetChallengeType(c)
		if err != nil {
			return
		}

		waitSeconds, err := GetWaitTime(c)
		if err != nil {
			return
		}

		timeoutSeconds, err := GetTimeout(c)
		if err != nil {
			return
		}

		lang, err := GetLanguage(c)
		if err != nil {
			return
		}

		challengeChecksum := GetClientTokenCheckSum(c)

		challenge := NewChallenge(challengeType, challengeChecksum, waitSeconds, timeoutSeconds)

		challengeToken, err := GenerateToken(challenge.JSON(), waitSeconds, timeoutSeconds, config.Salt, rsa)

		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("error")
			protectionStatus.SetType("challenge_token_generate_fail")
			protectionStatus.AddToResponse(c)

			message := "Cannot generate token"
			code := 500
			c.Status(code)
			c.Send(message)
			return
		}

		template := "challenge-" + challengeType + "." + lang

		// prometheus
		PrometheusGeneratedChallenge.With(prometheus.Labels{"type": challengeType})

		dataMap := fiber.Map{
			"Token":   challengeToken,
			"Wait":    waitSeconds,
			"Timeout": timeoutSeconds,
			"BaseURL": config.BaseURL,
		}

		c.Render(template, dataMap)
	})
}
