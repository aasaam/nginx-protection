package main

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

const defaultOrganizationName = "aasaam software development group"
const defaultCookieName = "aasaam_protection"

const misconfigureStatus = fiber.StatusNotImplemented

const httpRequestHeaderClientAPIKeyHeaderName = "x-protection-api-key"

const (
	localVarIP                      = "ip"
	localVarLang                    = "lang"
	localVarRequestID               = "request_id"
	localVarTOTPSecret              = "totp_secret"
	localVarChallengeType           = "challenge_type"
	localVarClientTemporaryChecksum = "client_temporary_checksum"
	localVarClientPersistChecksum   = "client_persist_checksum"
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
)

const (
	// Configurations:
	httpRequestHeaderRequestID                   = "X-Request-ID"
	httpRequestHeaderConfigNodeID                = "X-Protection-Config-Node-ID"
	httpRequestHeaderConfigSupportEmail          = "X-Protection-Config-Support-Email"
	httpRequestHeaderConfigSupportTel            = "X-Protection-Config-Support-Tel"
	httpRequestHeaderConfigSupportURL            = "X-Protection-Config-Support-URL"
	httpRequestHeaderConfigChallenge             = "X-Protection-Config-Challenge"
	httpRequestHeaderConfigLang                  = "X-Protection-Config-Lang"
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

	// Client:
	httpRequestHeaderClientTemporaryChecksum = "X-Protection-Client-Temporary-Checksum"
	httpRequestHeaderClientPersistChecksum   = "X-Protection-Client-Persist-Checksum"
	httpRequestHeaderClientCountry           = "X-Protection-Client-Country"
	httpRequestHeaderClientASNNumber         = "X-Protection-Client-ASN-Number"
	httpRequestHeaderClientASNOrganization   = "X-Protection-Client-ASN-Organization"

	// Acl:
	httpRequestHeaderAclCountries = "X-Protection-ACL-Countries"
	httpRequestHeaderAclCIDRs     = "X-Protection-ACL-CIDRs"
	httpRequestHeaderAclASNs      = "X-Protection-ACL-ASNs"
	httpRequestHeaderAclASNRanges = "X-Protection-ACL-ASN-Ranges"
	httpRequestHeaderAclAPIKeys   = "X-Protection-ACL-API-Keys"

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
	JSValue        string `json:"j_v,omitempty"`
	CaptchaValue   int    `json:"c_v,omitempty"`
	TOTPCode       string `json:"t_v,omitempty"`
}

type challengeResponse struct {
	ChallengeToken string `json:"t"`
	JSProblem      string `json:"js,omitempty"`
	CaptchaProblem string `json:"captcha,omitempty"`
}

var rtlLanguages = []string{"ar", "dv", "fa", "he", "ps", "ur", "yi"}
var rtlLanguagesMap map[string]bool

var supportedChallenges = []string{challengeTypeBlock, challengeTypeJS, challengeTypeCaptcha, challengeTypeTOTP}
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
