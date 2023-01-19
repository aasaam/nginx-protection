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

	}
	return true
}

func checkAuth(c *fiber.Ctx, config *config, realCheck bool) bool {
	ttl := getConfigTTLSeconds(c)
	persistChecksum := c.Locals(localVarClientPersistChecksum).(string)
	ip := c.Locals(localVarIP).(string)

	defer config.getLogger().
		Info().
		Str(logType, logTypeAuthCheck).
		Str(logPropertyIP, ip).
		Send()

	// api keys
	success, apiClientName := aclCheckAPIKeys(c)
	if success {
		return successResponse(config, persistChecksum, aclRuleAPI, apiClientName, "", ttl, ip, realCheck)
	}

	// country
	success, countryCode := aclCheckCountries(c)
	if success {
		return successResponse(config, persistChecksum, aclRuleCountry, countryCode, "", ttl, ip, realCheck)
	}

	// cidr
	success, cidr := aclCheckCIDRs(c)
	if success {
		return successResponse(config, persistChecksum, aclRuleCIDR, cidr, "", ttl, ip, realCheck)
	}

	// asn
	success, asn := aclCheckASNs(c)
	if success {
		return successResponse(config, persistChecksum, aclRuleASN, asn, "", ttl, ip, realCheck)
	}

	// cookie check
	cookieVar := c.Cookies(c.Get(httpRequestHeaderConfigCookie, defaultCookieName), "")
	if cookieVar != "" {
		cookieToken, cookieErr := newPersistTokenFromString(cookieVar, config.tokenSecret)
		if cookieErr == nil {
			return successResponse(config, persistChecksum, aclRuleChallenge, cookieToken.Type, cookieToken.Username, ttl, ip, realCheck)
		}
	}

	return failedResponse(config, ip, realCheck)
}
