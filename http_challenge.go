package main

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

func httpChallenge(c *fiber.Ctx, config *config) error {
	ip := c.Locals(localVarIP).(string)
	requestID := c.Locals(localVarRequestID).(string)
	challengeType := c.Locals(localVarChallengeType).(string)
	persistChecksum := c.Locals(localVarClientPersistChecksum).(string)
	temporaryChecksum := c.Locals(localVarClientTemporaryChecksum).(string)
	lang := getLanguage(c, config)
	challengeEmoji := "⌛️"
	waitSeconds := getConfigWaitSeconds(c)
	timeoutSeconds := getConfigTimeoutSeconds(c)

	defer prometheusRequestChallenge.WithLabelValues(challengeType).Inc()

	switch challengeType {
	case challengeTypeBlock:
		challengeEmoji = "⛔️"
	}

	var challenge *challenge

	challengeToken := ""

	if challengeType != challengeTypeBlock {
		ttl := getConfigTTLSeconds(c)

		challenge = newChallenge(lang, challengeType, temporaryChecksum, persistChecksum, waitSeconds, timeoutSeconds, ttl)
		challengeToken, _ = challenge.getChallengeToken(config.clientSecret)
	}

	if challengeType == challengeTypeTOTP {
		challenge.setTOTPSecret(c.Locals(localVarTOTPSecret).(string))
	}

	supportedLanguages, _ := json.Marshal(config.supportedLangauges)
	supportInfo, _ := json.Marshal(getConfigSupportInfo(c))
	ipData, _ := json.Marshal(getClientProperties(c))
	unixTime, _ := json.Marshal(time.Now().Unix())

	defer config.getLogger().Info().Str("challenge_type", challengeType).Str("ip", ip).Str("rid", requestID).Send()

	// set header
	c.Set(httpResponseChallengeToken, challengeToken)

	// return nil
	return c.Render("templates/"+challengeType, fiber.Map{

		// html
		"title":                 translateData[lang][challengeType],
		"dir":                   getLanguageDirection(lang),
		"staticURL":             config.staticURL,
		"i18n":                  translateData[lang],
		"supportedLanguages":    string(supportedLanguages),
		"multiLanguage":         len(config.supportedLangauges) > 1,
		"languageData":          languagesData(config.supportedLangauges, lang),
		"challengeEmoji":        challengeEmoji,
		"organizationTitle":     getConfigI18nOrganizationTitle(c, config),
		"organizationBrandIcon": getConfigI18nOrganizationBrandIcon(c),
		"challengeType":         challengeType,
		"persistChecksum":       persistChecksum,
		"cdnStatic":             config.cdnStatic,
		"aasaamWebServer":       config.aasaamWebServer,

		// js variables
		"lang":           lang,
		"unixTime":       string(unixTime),
		"challengeToken": challengeToken,
		"ipData":         string(ipData),
		"protectedPath":  getProtectedPath(c),
		"supportInfo":    string(supportInfo),
		"waitSeconds":    waitSeconds,
		"baseURL":        config.baseURL,
	}, "templates/layouts/main")
}
