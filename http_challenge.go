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
	challengeEmoji := "⌛️"
	waitSeconds := getConfigWaitSeconds(c)
	timeoutSeconds := getConfigTimeoutSeconds(c)

	supportedLangauges := getSupportedLangauges(c, config)
	lang := getLanguage(c, config)

	if !isSupportedLangaugeConfig(lang, supportedLangauges) {
		lang = config.defaultLanguage
		supportedLangauges = []string{lang}
	}

	defer prometheusRequestChallenge.WithLabelValues(challengeType).Inc()

	switch challengeType {
	case challengeTypeBlock:
		waitSeconds = 0
		timeoutSeconds = 0
		challengeEmoji = "⛔️"
	case challengeTypeCaptcha:
		challengeEmoji = getHumanEmoji()
	case challengeTypeTOTP:
		challengeEmoji = "🔐"
	case challengeTypeLDAP:
		challengeEmoji = "🛂"
	}

	var challenge *challenge
	challengeToken := ""

	if challengeType != challengeTypeBlock {
		ttl := getConfigTTLSeconds(c)

		challenge = newChallenge(lang, challengeType, temporaryChecksum, persistChecksum, waitSeconds, timeoutSeconds, ttl)
	}

	if challengeType == challengeTypeTOTP {
		challenge.setTOTPSecret(c.Locals(localVarTOTPSecret).(string))
	}

	languageData := languagesData(supportedLangauges, lang)

	supportedLanguages, _ := json.Marshal(supportedLangauges)
	supportInfo, _ := json.Marshal(getConfigSupportInfo(c))
	ipData, _ := json.Marshal(getClientProperties(c))
	unixTime, _ := json.Marshal(time.Now().Unix())

	defer config.getLogger().
		Info().
		Str(logType, logTypeChallengeGenerate).
		Str(logPropertyIP, ip).
		Str(logPropertyRequestID, requestID).
		Str(logPropertyChallengeType, challengeType).
		Send()

	if challenge != nil {
		challengeToken, _ = challenge.getChallengeToken(config.clientSecret)
	}

	// set header
	c.Set(httpResponseChallengeTemporary, challengeToken)

	// return nil
	return c.Render("templates/"+challengeType, fiber.Map{

		// html
		"title":                 translateData[lang][challengeType],
		"dir":                   getLanguageDirection(lang),
		"staticURL":             config.staticURL,
		"i18n":                  translateData[lang],
		"supportedLanguages":    string(supportedLanguages),
		"multiLanguage":         len(config.supportedLangauges) > 1,
		"languageData":          languageData,
		"challengeEmoji":        challengeEmoji,
		"organizationTitle":     getConfigI18nOrganizationTitle(c, config, lang),
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
		"timeoutSeconds": timeoutSeconds,
		"baseURL":        config.baseURL,
	}, "templates/layouts/main")
}
