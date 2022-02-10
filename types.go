package main

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

const defaultOrganizationName = "aasaam software development group"
const defaultCookieName = "aasaam_protection"
const defaultUsernameCookieName = "protection_auth_user_name"

const misconfigureStatus = fiber.StatusNotImplemented

const httpRequestHeaderClientAPIKeyHeaderName = "x-protection-api-key"

const ldapUsernamePlaceHolder = "__USERNAME__"

const (
	localVarIP                      = "ip"
	localVarLang                    = "lang"
	localVarRequestID               = "request_id"
	localVarTOTPSecret              = "totp_secret"
	localVarChallengeType           = "challenge_type"
	localVarClientTemporaryChecksum = "client_temporary_checksum"
	localVarClientPersistChecksum   = "client_persist_checksum"
	localVarLDAPURL                 = "ldap_url"
	localVarLDAPReadonlyUsername    = "ldap_readonly_username"
	localVarLDAPReadonlyPassword    = "ldap_readonly_password"
	localVarLDAPBaseDN              = "ldap_base_dn"
	localVarLDAPFilter              = "ldap_filter"
	localVarLDAPAttributes          = "ldap_attributes"
)

const (
	aclRuleChallenge = "challenge"
	aclRuleCountry   = "country"
	aclRuleCIDR      = "cidr"
	aclRuleASN       = "asn"
	aclRuleAPI       = "api"
)

const (
	challengeTypeBlock   = "block"
	challengeTypeJS      = "js"
	challengeTypeCaptcha = "captcha"
	challengeTypeTOTP    = "totp"
	challengeTypeLDAP    = "ldap"
)

const (
	logType                  = "type"
	logTypeACL               = "acl"
	logTypeLDAPError         = "ldap_error"
	logTypeChallengeGenerate = "challenge_new"
	logTypeChallengeFailed   = "challenge_failed"
	logTypeChallengeSuccess  = "challenge_success"
	logTypeAuthFailed        = "auth_failed"
	logTypeAuthSuccess       = "auth_success"
	logTypeAuthCheck         = "auth_check"
	logTypeAuthCache         = "auth_cache"
	logTypeHTTPError         = "http_error"
	logTypeHTTPRequest       = "http_request"
	logTypeApp               = "app"
)

const (
	logPropertyIP            = "ip"
	logPropertyURL           = "url"
	logPropertyChallengeType = "challenge_type"
	logPropertyMethod        = "method"
	logPropertyStatusCode    = "status"
	logPropertyAuth          = "auth"
	logPropertyError         = "err"
	logPropertyACL           = "acl"
	logPropertyValue         = "value"
	logPropertyRequestID     = "rid"
	logPropertyUsername      = "username"
)

const (
	// Configurations:
	httpRequestHeaderXForwardedFor               = "X-Forwarded-For"
	httpRequestHeaderRequestID                   = "X-Request-ID"
	httpRequestHeaderConfigNodeID                = "X-Protection-Config-Node-ID"
	httpRequestHeaderConfigSupportEmail          = "X-Protection-Config-Support-Email"
	httpRequestHeaderConfigSupportTel            = "X-Protection-Config-Support-Tel"
	httpRequestHeaderConfigSupportURL            = "X-Protection-Config-Support-URL"
	httpRequestHeaderConfigChallenge             = "X-Protection-Config-Challenge"
	httpRequestHeaderConfigLang                  = "X-Protection-Config-Lang"
	httpRequestHeaderConfigSupportedLanguages    = "X-Protection-Config-Supported-Languages"
	httpRequestHeaderConfigCaptchaDifficulty     = "X-Protection-Config-Captcha-Difficulty"
	httpRequestHeaderConfigTTL                   = "X-Protection-Config-TTL"
	httpRequestHeaderConfigTimeout               = "X-Protection-Config-Timeout"
	httpRequestHeaderConfigTOTPSecret            = "X-Protection-Config-TOTP-Secret"
	httpRequestHeaderConfigWait                  = "X-Protection-Config-Wait"
	httpRequestHeaderConfigCookie                = "X-Protection-Config-Cookie"
	httpRequestHeaderConfigI18nOrganizationTitle = "X-Protection-Config-I18n-Organization-Title"
	httpRequestHeaderConfigOrganizationTitle     = "X-Protection-Config-Organization-Title"
	httpRequestHeaderConfigOrganizationBrandIcon = "X-Protection-Config-Organization-Brand-Icon"
	httpRequestHeaderConfigUnauthorizedStatus    = "X-Protection-Config-Unauthorized-Status"
	httpRequestHeaderConfigLDAPURL               = "X-Protection-Config-LDAP-URL"
	httpRequestHeaderConfigLDAPReadonlyUsername  = "X-Protection-Config-LDAP-Readonly-Username"
	httpRequestHeaderConfigLDAPReadonlyPassword  = "X-Protection-Config-LDAP-Readonly-Password"
	httpRequestHeaderConfigLDAPBaseDN            = "X-Protection-Config-LDAP-Base-DN"
	httpRequestHeaderConfigLDAPFilter            = "X-Protection-Config-LDAP-Filter"
	httpRequestHeaderConfigLDAPAttributes        = "X-Protection-Config-LDAP-Attributes"
	httpRequestHeaderConfigLDAPLoginPattern      = "X-Protection-Config-LDAP-Login-Pattern"

	// Client:
	httpRequestHeaderClientTemporaryChecksum = "X-Protection-Client-Temporary-Checksum"
	httpRequestHeaderClientPersistChecksum   = "X-Protection-Client-Persist-Checksum"
	httpRequestHeaderClientCountry           = "X-Protection-Client-Country"
	httpRequestHeaderClientASNNumber         = "X-Protection-Client-ASN-Number"
	httpRequestHeaderClientASNOrganization   = "X-Protection-Client-ASN-Organization"

	// ACL:
	httpRequestHeaderACLCountries = "X-Protection-ACL-Countries"
	httpRequestHeaderACLCIDRs     = "X-Protection-ACL-CIDRs"
	httpRequestHeaderACLASNs      = "X-Protection-ACL-ASNs"
	httpRequestHeaderACLASNRanges = "X-Protection-ACL-ASN-Ranges"
	httpRequestHeaderACLAPIKeys   = "X-Protection-ACL-API-Keys"

	httpResponseChallengeToken  = "X-Challenge-Token"
	httpResponseChallengeResult = "X-Challenge-Result"
)

type supportInfo struct {
	Email string `json:"email"`
	Tel   string `json:"tel"`
	URL   string `json:"url"`
}

type challengeRequest struct {
	ChallengeToken string `json:"t"`
	LDAPUsername   string `json:"lu,omitempty"`
	LDAPPassword   string `json:"lp,omitempty"`
	JSValue        string `json:"j,omitempty"`
	CaptchaValue   int    `json:"c,omitempty"`
	TOTPPassword   string `json:"totp,omitempty"`
}

type challengeResponse struct {
	ChallengeToken string `json:"t"`
	JSProblem      string `json:"js,omitempty"`
	CaptchaProblem string `json:"captcha,omitempty"`
}

var humanEmojis = []string{"🧑‍💻", "🧑🏻‍💻", "🧑🏼‍💻", "🧑🏽‍💻", "🧑🏾‍💻", "🧑🏿‍💻"}

var rtlLanguages = []string{"ar", "dv", "fa", "he", "ps", "ur", "yi"}
var rtlLanguagesMap map[string]bool

var supportedChallenges = []string{challengeTypeBlock, challengeTypeJS, challengeTypeCaptcha, challengeTypeTOTP, challengeTypeLDAP}
var supportedChallengesMap map[string]bool

var supportedLangauges = []string{"fa", "en"}
var supportedLangaugesMap map[string]bool

var totpSecretRegex = regexp.MustCompile(`^([0-9A-Z]{16})$`)

func init() {
	supportedLangaugesMap = make(map[string]bool)
	for _, v := range supportedLangauges {
		supportedLangaugesMap[v] = true
	}
	supportedChallengesMap = make(map[string]bool)
	for _, v := range supportedChallenges {
		supportedChallengesMap[v] = true
	}
	rtlLanguagesMap = make(map[string]bool)
	for _, v := range rtlLanguages {
		rtlLanguagesMap[v] = true
	}
}
