package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func failedResponse(config *config, ip string) bool {
	defer prometheusRequestAuthFailed.Inc()
	defer config.getLogger().Info().Str("ip", ip).Msg("auth blocked")
	return false
}

func successResponse(
	config *config,
	aclStorage *aclStorage,
	clientPersistChecksum string,
	aclRule string,
	name string,
	ttl int64,
) bool {
	defer prometheusRequestAuthSuccess.With(prometheus.Labels{"acl": aclRule, "value": name}).Inc()
	defer config.getLogger().Info().Str("acl", aclRule).Str("var", name).Msg("auth allowed")
	aclStorage.add(clientPersistChecksum, aclRule, name, minMaxDefault64(ttl, 60, 600))
	return true
}

func checkAuth(c *fiber.Ctx, config *config, aclStorage *aclStorage) bool {
	ttl := getConfigTTLSeconds(c)
	persistChecksum := c.Locals(localVarClientPersistChecksum).(string)
	ip := c.Locals(localVarIP).(string)

	storageItem := aclStorage.exist(persistChecksum)
	if storageItem != nil {
		defer config.getLogger().Info().Str("ip", ip).Str("checksum", persistChecksum).Msg("auth allowed using cache")
		return successResponse(config, aclStorage, persistChecksum, storageItem.rule, storageItem.name, ttl)
	}

	defer config.getLogger().Info().Str("ip", ip).Str("checksum", persistChecksum).Msg("try to process request")

	// api keys
	success, apiClientName := aclCheckAPIKeys(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleAPI, apiClientName, ttl)
	}

	// country
	success, countryCode := aclCheckCountries(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleCountry, countryCode, ttl)
	}

	// cidr
	success, cidr := aclCheckCIDRs(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleCIDR, cidr, ttl)
	}

	// asn
	success, asn := aclCheckASNs(c)
	if success {
		return successResponse(config, aclStorage, persistChecksum, aclRuleASN, asn, ttl)
	}

	// cookie check
	cookieVar := c.Cookies(c.Get(httpRequestHeaderConfigCookie, defaultCookieName), "")
	if cookieVar != "" {
		cookieToken, cookieErr := newPersistTokenFromString(cookieVar, config.tokenSecret)
		if cookieErr == nil {
			return successResponse(config, aclStorage, persistChecksum, aclRuleChallenge, cookieToken.Type, ttl)
		}
	}

	return failedResponse(config, ip)
}

// func httpAuth(c *fiber.Ctx, config *config, aclStorage *aclStorage) (bool, error) {

//
// }
