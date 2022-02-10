package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func failedResponse(config *config, ip string, realCheck bool) bool {
	if realCheck {
		defer prometheusRequestAuthFailed.Inc()
		defer config.getLogger().
			Info().
			Str(logType, logTypeAuthFailed).
			Str(logPropertyIP, ip).
			Send()
	}
	return false
}

func successResponse(
	config *config,
	aclStorage *aclStorage,
	clientPersistChecksum string,
	aclRule string,
	value string,
	username string,
	ttl int64,
	ip string,
	realCheck bool,
) bool {
	if realCheck {
		defer prometheusRequestAuthSuccess.With(prometheus.Labels{"acl": aclRule, "value": value}).Inc()

		defer config.getLogger().
			Info().
			Str(logType, logTypeAuthSuccess).
			Str(logPropertyIP, ip).
			Str(logPropertyACL, aclRule).
			Str(logPropertyValue, value).
			Str(logPropertyUsername, username).
			Send()

		aclStorage.add(clientPersistChecksum, aclRule, value, username, minMaxDefault64(ttl, 60, 600))
	}
	return true
}

func checkAuth(c *fiber.Ctx, config *config, aclStorage *aclStorage, realCheck bool) bool {
	ttl := getConfigTTLSeconds(c)
	persistChecksum := c.Locals(localVarClientPersistChecksum).(string)
	ip := c.Locals(localVarIP).(string)

	storageItem := aclStorage.exist(persistChecksum)
	if storageItem != nil {
		defer config.getLogger().
			Info().
			Str(logType, logTypeAuthCache).
			Str(logPropertyIP, ip).
			Str(logPropertyUsername, storageItem.userName).
			Send()
		return successResponse(config, aclStorage, persistChecksum, storageItem.rule, storageItem.name, storageItem.userName, ttl, ip, realCheck)
	}

	defer config.getLogger().
		Info().
		Str(logType, logTypeAuthCheck).
		Str(logPropertyIP, ip).
		Send()

	// api keys
	success, apiClientName := aclCheckAPIKeys(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleAPI, apiClientName, "", ttl, ip, realCheck)
	}

	// country
	success, countryCode := aclCheckCountries(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleCountry, countryCode, "", ttl, ip, realCheck)
	}

	// cidr
	success, cidr := aclCheckCIDRs(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleCIDR, cidr, "", ttl, ip, realCheck)
	}

	// asn
	success, asn := aclCheckASNs(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleASN, asn, "", ttl, ip, realCheck)
	}

	// cookie check
	cookieVar := c.Cookies(c.Get(httpRequestHeaderConfigCookie, defaultCookieName), "")
	if cookieVar != "" {
		cookieToken, cookieErr := newPersistTokenFromString(cookieVar, config.tokenSecret)
		if cookieErr == nil {
			return successResponse(config, aclStorage, persistChecksum, aclRuleChallenge, cookieToken.Type, cookieToken.Username, ttl, ip, realCheck)
		}
	}

	return failedResponse(config, ip, realCheck)
}
