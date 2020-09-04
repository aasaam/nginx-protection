package main

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AuthToken is token for auth validation
type AuthToken struct {
	Checksum string `json:"checksum"`
	Type     string `json:"type"`
	User     string `json:"user"`
}

// ChallengeToken is encrypted challenge data for transport layer
type ChallengeToken struct {
	Token   string `json:"token"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Value   string `json:"value"`
}

// ProtectionStatus is status response for header
type ProtectionStatus struct {
	Status string            `json:"status"`
	Type   string            `json:"type"`
	User   string            `json:"user"`
	Extras map[string]string `json:"extra"`
}

// Config the application
type Config struct {
	Salt    string
	BaseURL string
	Testing bool
}

// OTPConfig otp config
type OTPConfig struct {
	Secret string
	Period uint
}

// SupportedLanguage is list of language
var SupportedLanguage = []string{"fa", "en"}

// CaptchaLetters is letters for generation captcha code
var CaptchaLetters = []rune("123456789")

// OtpSecretRegex is regex for check OTP secret is good
var OtpSecretRegex, _ = regexp.Compile(`^[0-9A-Z]{16}$`)

// SupportedChallenges is slice of supported challenges
var SupportedChallenges = []string{ChallengeTypeJS, ChallengeTypeCaptcha, ChallengeTypeOTP, ChallengeTypeUserPass, ChallengeTypeSMS}

var (
	// PrometheusRequestTotal prometheus counter
	PrometheusRequestTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "protection_request_total",
		Help: "The total number of request for protection",
	})

	// PrometheusRejectTotal prometheus counter
	PrometheusRejectTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "protection_reject_total",
		Help: "The total number of reject for protection",
	})

	// PrometheusAccpetTotal prometheus counter
	PrometheusAccpetTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_accpet_total",
		Help: "The total number of accpet for protection",
	}, []string{"type"})

	// PrometheusGeneratedChallenge prometheus counter
	PrometheusGeneratedChallenge = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_generated_challenge_total",
		Help: "The total number of generated challenge",
	}, []string{"type"})
)

const (
	// ChallengeTypeJS type js
	ChallengeTypeJS = "js"
	// ChallengeTypeCaptcha type captcha
	ChallengeTypeCaptcha = "captcha"
	// ChallengeTypeOTP type otp
	ChallengeTypeOTP = "otp"
	// ChallengeTypeUserPass type user-pass
	ChallengeTypeUserPass = "user-pass"
	// ChallengeTypeSMS type sms
	ChallengeTypeSMS = "sms"
)

const (
	// Standards:

	// HTTPRequestXForwardedFor is standard header for client ip
	HTTPRequestXForwardedFor = "X-Forwarded-For"

	// Configurations:

	// HTTPRequestHeaderConfigChallenge config for challenge type
	HTTPRequestHeaderConfigChallenge = "X-Protection-Config-Challenge"
	// HTTPRequestHeaderConfigTTL config for time to live of generated token
	HTTPRequestHeaderConfigTTL = "X-Protection-Config-TTL"
	// HTTPRequestHeaderConfigTimeout how long client token for solving need time
	HTTPRequestHeaderConfigTimeout = "X-Protection-Config-Timeout"
	// HTTPRequestHeaderConfigLang config theme language
	HTTPRequestHeaderConfigLang = "X-Protection-Config-Lang"
	// HTTPRequestHeaderConfigOTPSecret is secret for OTP challenge
	HTTPRequestHeaderConfigOTPSecret = "X-Protection-Config-OTP-Secret"
	// HTTPRequestHeaderConfigOTPTime is time for OTP challenge
	HTTPRequestHeaderConfigOTPTime = "X-Protection-Config-OTP-Time"
	// HTTPRequestHeaderConfigFarsiCaptcha is boolean type for farsi captcha
	HTTPRequestHeaderConfigFarsiCaptcha = "X-Protection-Config-FarsiCaptcha"
	// HTTPRequestHeaderConfigWait is number of seconds to wait for start challenge
	HTTPRequestHeaderConfigWait = "X-Protection-Config-Wait"
	// HTTPRequestHeaderConfigCookie name of cookie for set token
	HTTPRequestHeaderConfigCookie = "X-Protection-Config-Cookie"

	// Client:

	// HTTPRequestHeaderClientChecksum is client checksum for identifier
	HTTPRequestHeaderClientChecksum = "X-Protection-Client-Checksum"
	// HTTPRequestHeaderClientTokenChecksum is client checksum for token challenges
	HTTPRequestHeaderClientTokenChecksum = "X-Protection-Client-Token-Checksum"
	// HTTPRequestHeaderClientAPIKey is client api key for authentication
	HTTPRequestHeaderClientAPIKey = "X-Protection-Client-API-Key"
	// HTTPRequestHeaderClientCountry is client country iso code
	HTTPRequestHeaderClientCountry = "X-Protection-Client-Country"
	// HTTPRequestHeaderClientASNNumber is client network ASN number
	HTTPRequestHeaderClientASNNumber = "X-Protection-Client-ASN-Number"
	// HTTPRequestHeaderClientASNOrganization is client network ASN name
	HTTPRequestHeaderClientASNOrganization = "X-Protection-Client-ASN-Organization"

	// ACL:

	// HTTPRequestHeaderACLCountries is ACL for countries
	HTTPRequestHeaderACLCountries = "X-Protection-ACL-Countries"
	// HTTPRequestHeaderACLCIDRs is ACL for CIDRs
	HTTPRequestHeaderACLCIDRs = "X-Protection-ACL-CIDRs"
	// HTTPRequestHeaderACLASNs is ACL for ASNs
	HTTPRequestHeaderACLASNs = "X-Protection-ACL-ASNs"
	// HTTPRequestHeaderACLAPIKeys list of map[string]string of possible api keys
	HTTPRequestHeaderACLAPIKeys = "X-Protection-ACL-API-Keys"

	// Responses:

	// HTTPResponseHeaderProtectionStatus is response header for status
	HTTPResponseHeaderProtectionStatus = "X-Protection-Status"
	// HTTPResponseHeaderProtectionType is response header for type
	HTTPResponseHeaderProtectionType = "X-Protection-Type"
	// HTTPResponseHeaderProtectionUser is response header for special user identifier
	HTTPResponseHeaderProtectionUser = "X-Protection-User"
	// HTTPResponseHeaderProtectionExtra is response header for extra data
	HTTPResponseHeaderProtectionExtra = "X-Protection-Extra"
)
