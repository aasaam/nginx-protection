package main

import (
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type clientProperties struct {
	IP              string `json:"ip"`
	Country         string `json:"country"`
	ASNNumber       string `json:"asn"`
	ASNOrganization string `json:"asn_org"`
	NodeID          string `json:"nodeID"`
}

var asnRangeRegexp = regexp.MustCompile(`^([0-9]{1,6})-([0-9]{1,6})$`)

func checkConfiguration(c *fiber.Ctx) error {
	ip, err := getRequestIP(c)
	if err != nil {
		return err
	}

	challengeType := c.Get(httpRequestHeaderConfigChallenge, "")
	if !isValidChallenge(challengeType) {
		return errors.New("invalid challenge type")
	}

	clientTemporaryChecksum := c.Get(httpRequestHeaderClientTemporaryChecksum, "")
	if len(clientTemporaryChecksum) < 16 {
		return errors.New("invalid client temporary checksum")
	}

	clientPersistChecksum := c.Get(httpRequestHeaderClientPersistChecksum, "")
	if len(clientPersistChecksum) < 16 {
		return errors.New("invalid client persist checksum")
	}

	totpSecret := ""
	if challengeType == challengeTypeTOTP {
		totpSecret = c.Get(httpRequestHeaderConfigTOTPSecret, "")
		if !totpSecretRegex.MatchString(totpSecret) {
			return errors.New("invalid totp secret")
		}
	}

	ldapURL := ""
	ldapReadonlyUsername := ""
	ldapReadonlyPassword := ""
	ldapBaseDN := ""
	ldapFilter := ""
	ldapAttributes := ""
	if challengeType == challengeTypeLDAP {
		ldapURL = c.Get(httpRequestHeaderConfigLDAPURL, "")
		if ldapURL == "" {
			return errors.New("invalid ldap url")
		}
		ldapReadonlyUsername = c.Get(httpRequestHeaderConfigLDAPReadonlyUsername, "")
		if ldapReadonlyUsername == "" {
			return errors.New("invalid ldap readonly username")
		}
		ldapReadonlyPassword = c.Get(httpRequestHeaderConfigLDAPReadonlyPassword, "")
		if ldapReadonlyPassword == "" {
			return errors.New("invalid ldap readonly password")
		}
		ldapBaseDN = c.Get(httpRequestHeaderConfigLDAPBaseDN, "")
		if ldapBaseDN == "" {
			return errors.New("invalid ldap base dn")
		}
		ldapFilter = c.Get(httpRequestHeaderConfigLDAPFilter, "")
		if ldapFilter == "" {
			return errors.New("invalid ldap filter")
		}
		ldapAttributes = c.Get(httpRequestHeaderConfigLDAPAttributes, "")
		if ldapAttributes == "" {
			return errors.New("invalid ldap attributes")
		}
	}

	requestID := c.Get(httpRequestHeaderRequestID, "")
	if len(requestID) < 16 {
		return errors.New("invalid request id")
	}

	c.Locals(localVarIP, ip)
	c.Locals(localVarRequestID, requestID)
	c.Locals(localVarTOTPSecret, totpSecret)
	c.Locals(localVarChallengeType, challengeType)

	// ldap
	c.Locals(localVarLDAPURL, ldapURL)
	c.Locals(localVarLDAPBaseDN, ldapBaseDN)
	c.Locals(localVarLDAPReadonlyUsername, ldapReadonlyUsername)
	c.Locals(localVarLDAPReadonlyPassword, ldapReadonlyPassword)
	c.Locals(localVarLDAPFilter, ldapFilter)
	c.Locals(localVarLDAPAttributes, ldapAttributes)

	c.Locals(localVarClientTemporaryChecksum, base64Hash(clientTemporaryChecksum))
	c.Locals(localVarClientPersistChecksum, base64Hash(clientPersistChecksum))
	return nil
}

func getRequestIP(c *fiber.Ctx) (string, error) {
	ip := c.Get(httpRequestHeaderXForwardedFor, "")
	if ip == "" {
		return "", errors.New("ip address not found")
	}
	ips := strings.Split(strings.TrimSpace(ip), ",")
	if len(ips) > 1 {
		ip = strings.TrimSpace(ips[0])
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("invalid ip address")
}

func getClientProperties(c *fiber.Ctx) *clientProperties {
	ip := c.Locals(localVarIP).(string)
	p := clientProperties{
		IP:              ip,
		NodeID:          c.Get(httpRequestHeaderConfigNodeID, ""),
		Country:         c.Get(httpRequestHeaderClientCountry, ""),
		ASNNumber:       c.Get(httpRequestHeaderClientASNNumber, ""),
		ASNOrganization: c.Get(httpRequestHeaderClientASNOrganization, ""),
	}

	return &p
}

// getConfigWaitSeconds rage is 1-180 seconds
func getConfigWaitSeconds(c *fiber.Ctx) int64 {
	wait := c.Get(httpRequestHeaderConfigWait, "0")
	waitSeconds, err := strconv.ParseInt(wait, 10, 64)
	if err != nil {
		return 3
	}
	return minMaxDefault64(waitSeconds, 2, 180)
}

// getConfigTimeoutSeconds rage is 300-1800 seconds
func getConfigTimeoutSeconds(c *fiber.Ctx) int64 {
	timeout := c.Get(httpRequestHeaderConfigTimeout, "0")
	timeoutSeconds, err := strconv.ParseInt(timeout, 10, 64)
	if err != nil {
		return 600
	}
	return minMaxDefault64(timeoutSeconds, 300, 1800)
}

// getConfigTTLSeconds rage is 3600-604800 seconds (1 hour to 1 week)
func getConfigTTLSeconds(c *fiber.Ctx) int64 {
	timeout := c.Get(httpRequestHeaderConfigTTL, "0")
	timeoutSeconds, err := strconv.ParseInt(timeout, 10, 64)
	if err != nil {
		return 28800
	}
	return minMaxDefault64(timeoutSeconds, 3600, 604800)
}

func getConfigUnauthorizedStatus(c *fiber.Ctx) int {
	status := c.Get(httpRequestHeaderConfigUnauthorizedStatus, "0")
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		return 403
	}
	return minMaxDefault(statusCode, 400, 499)
}

func getSupportedLangauges(c *fiber.Ctx, config *config) []string {
	supportedLanguages := []string{}
	headerValue := strings.Split(c.Get(httpRequestHeaderConfigSupportedLanguages, ""), ",")
	for _, lang := range headerValue {
		if isSupportedLangauge(lang) {
			supportedLanguages = append(supportedLanguages, lang)
		}
	}
	return supportedLanguages
}

func getLanguage(c *fiber.Ctx, config *config) string {
	queryLang := c.Query("lang", "")
	if isSupportedLangauge(queryLang) {
		return queryLang
	}
	headerLang := c.Get(httpRequestHeaderConfigLang, "")
	if isSupportedLangauge(headerLang) {
		return headerLang
	}
	return config.defaultLanguage
}

func getConfigSupportInfo(c *fiber.Ctx) *supportInfo {
	supportInfo := supportInfo{
		Email: c.Get(httpRequestHeaderConfigSupportEmail, ""),
		Tel:   c.Get(httpRequestHeaderConfigSupportTel, ""),
		URL:   c.Get(httpRequestHeaderConfigSupportURL, ""),
	}
	return &supportInfo
}

func getConfigI18nOrganizationTitle(c *fiber.Ctx, config *config, lang string) string {
	title := c.Get(httpRequestHeaderConfigOrganizationTitle, defaultOrganizationName)

	v := c.Get(httpRequestHeaderConfigI18nOrganizationTitle, "_")
	var titleMaps map[string]string
	e := json.Unmarshal([]byte(v), &titleMaps)
	if e == nil {
		if i18nTitle, ok := titleMaps[lang]; ok {
			title = i18nTitle
		}
	}

	return title
}

func getConfigI18nOrganizationBrandIcon(c *fiber.Ctx) string {
	return c.Get(httpRequestHeaderConfigOrganizationBrandIcon, "")
}

func getProtectedPath(c *fiber.Ctx) string {
	queryParam := c.Query("u", "")
	if queryParam != "" {
		urlFull, errFull := url.ParseRequestURI(queryParam)
		if errFull == nil {
			return getURLPath(urlFull)
		}
		url, err := url.ParseRequestURI("http://localhost/" + strings.TrimLeft(queryParam, "/"))
		if err == nil {
			return getURLPath(url)
		}
	}
	// default is back to home
	return "/"
}

func aclCheckCIDRs(c *fiber.Ctx) (bool, string) {
	ipString := c.Locals(localVarIP).(string)
	ip := net.ParseIP(ipString)
	for _, cidr := range strings.Split(c.Get(httpRequestHeaderACLCIDRs, ""), ",") {
		cleanCIDR := strings.TrimSpace(cidr)
		_, ipv4Net, err := net.ParseCIDR(cleanCIDR)
		if err == nil && ipv4Net.Contains(ip) {
			return true, cleanCIDR
		}
	}
	return false, ""
}

func aclCheckASNs(c *fiber.Ctx) (bool, string) {
	clientASN := c.Get(httpRequestHeaderClientASNNumber, "")
	if clientASN == "" {
		return false, ""
	}

	clientASNInt, err := strconv.Atoi(clientASN)
	if err != nil {
		return false, ""
	}

	asnString := strconv.Itoa(clientASNInt)

	for _, asnNum := range strings.Split(c.Get(httpRequestHeaderACLASNs, ""), ",") {
		v := strings.TrimSpace(asnNum)
		vInt, err := strconv.Atoi(v)
		if err == nil && clientASNInt == vInt {
			return true, asnString
		}
	}

	for _, asnRange := range strings.Split(c.Get(httpRequestHeaderACLASNRanges, ""), ",") {
		v := strings.TrimSpace(asnRange)

		if asnRangeRegexp.MatchString(v) {
			matched := asnRangeRegexp.FindStringSubmatch(v)
			fromNumber, err1 := strconv.Atoi(matched[1])
			toNumber, err2 := strconv.Atoi(matched[2])
			if err1 != nil || err2 != nil {
				continue
			}
			if clientASNInt >= fromNumber && clientASNInt <= toNumber {
				return true, asnString
			}
		}
	}

	return false, ""
}

func aclCheckCountries(c *fiber.Ctx) (bool, string) {
	clientCountry := c.Get(httpRequestHeaderClientCountry, "")
	if clientCountry == "" {
		return false, ""
	}

	for _, isoCode := range strings.Split(c.Get(httpRequestHeaderACLCountries, ""), ",") {
		code := strings.TrimSpace(isoCode)
		if len(code) == 2 && clientCountry == code {
			return true, code
		}
	}

	return false, ""
}

func aclCheckAPIKeys(c *fiber.Ctx) (bool, string) {
	clientAPIKey := c.Get(httpRequestHeaderClientAPIKeyHeaderName, "")
	if clientAPIKey == "" {
		return false, ""
	}

	v := c.Get(httpRequestHeaderACLAPIKeys, "_")
	var clientKeyMap map[string]string
	e := json.Unmarshal([]byte(v), &clientKeyMap)
	if e == nil {
		for clientName, clientKey := range clientKeyMap {
			if clientKey == clientAPIKey {
				return true, clientName
			}
		}
	}

	return false, ""
}
