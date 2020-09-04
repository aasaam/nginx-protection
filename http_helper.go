package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
)

// GetClientCheckSum get client checksum for auth identifier
func GetClientCheckSum(c *fiber.Ctx) string {
	return MD5(c.Get(HTTPRequestHeaderClientChecksum, "_"))
}

// GetClientTokenCheckSum get client temporary token checksum for challenge
func GetClientTokenCheckSum(c *fiber.Ctx) string {
	return MD5(c.Get(HTTPRequestHeaderClientTokenChecksum, "_"))
}

// GetRequestIP get client ip or return error
func GetRequestIP(c *fiber.Ctx) (string, error) {
	ip := c.Get(HTTPRequestXForwardedFor, "")
	if ip == "" {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("no_ip")
		protectionStatus.AddToResponse(c)

		message := "Client IP not found"
		code := 500
		c.Status(code)
		c.Send(message)
		return "", errors.New(message)
	}
	ips := strings.Split(ip, ",")
	if len(ips) > 1 {
		ip = strings.TrimSpace(ips[0])
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	message := "Client IP not found"
	code := 500
	c.Status(code)
	c.Send(message)
	return "", errors.New(message)
}

// GetFarsiCaptcha is farsi captcha or not
func GetFarsiCaptcha(c *fiber.Ctx) (bool, error) {
	farsiCaptcha := c.Get(HTTPRequestHeaderConfigFarsiCaptcha, "")
	isFarsiCaptcha, err := strconv.ParseBool(farsiCaptcha)
	if err != nil {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_farsi_captcha")
		protectionStatus.AddToResponse(c)

		message := "Invalid Farsi Captcha"
		code := 500
		c.Status(code)
		c.Send(message)
		return true, errors.New(message)
	}
	return isFarsiCaptcha, nil
}

// GetChallengeType get challenge type or return error
func GetChallengeType(c *fiber.Ctx) (string, error) {
	challengeType := c.Get(HTTPRequestHeaderConfigChallenge, "")
	if IsValidChallenge(challengeType) == false {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_challenge")
		protectionStatus.AddToResponse(c)

		message := "Invalid Challenge"
		code := 500
		c.Status(code)
		c.Send(message)
		return "", errors.New(message)
	}
	return challengeType, nil
}

// GetLanguage get language or return error
func GetLanguage(c *fiber.Ctx) (string, error) {
	lang := c.Get(HTTPRequestHeaderConfigLang, "")
	queryLang := c.Query("lang", "")
	if IsValidLanguage(lang) == false && IsValidLanguage(queryLang) == false {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_language")
		protectionStatus.AddToResponse(c)

		message := "Invalid Language"
		code := 500
		c.Status(code)
		c.Send(message)
		return "", errors.New(message)
	}
	if queryLang != "" {
		return queryLang, nil
	}
	return lang, nil
}

// GetWaitTime get wait time or return error
func GetWaitTime(c *fiber.Ctx) (int64, error) {
	wait := c.Get(HTTPRequestHeaderConfigWait, "0")
	waitSeconds, err := strconv.ParseInt(wait, 10, 64)
	if err != nil || waitSeconds < 3 || waitSeconds > 180 {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_wait")
		protectionStatus.AddToResponse(c)

		message := "Invalid Wait time"
		code := 500
		c.Status(code)
		c.Send(message)
		return 0, errors.New(message)
	}
	return waitSeconds, nil
}

// GetTimeout get timeout or return error
func GetTimeout(c *fiber.Ctx) (int64, error) {
	timeout := c.Get(HTTPRequestHeaderConfigTimeout, "0")
	timeoutSeconds, err := strconv.ParseInt(timeout, 10, 64)
	if err != nil || timeoutSeconds < 3 || timeoutSeconds > 300 {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_timeout")
		protectionStatus.AddToResponse(c)

		message := "Invalid Timeout"
		code := 500
		c.Status(code)
		c.Send(message)
		return 0, errors.New(message)
	}
	return timeoutSeconds, nil
}

// GetChallengeTTL get challenge time to live
func GetChallengeTTL(c *fiber.Ctx) (int64, error) {
	ttl := c.Get(HTTPRequestHeaderConfigTTL, "0")
	ttlSeconds, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil || ttlSeconds < 60 || ttlSeconds > 604800 {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("invalid_challenge_ttl")
		protectionStatus.AddToResponse(c)

		message := "Invalid Challenge TTL"
		code := 500
		c.Status(code)
		c.Send(message)
		return 0, errors.New(message)
	}
	return ttlSeconds, nil
}

// GetCookieName get cookie name
func GetCookieName(c *fiber.Ctx) (string, error) {
	cookieName := c.Get(HTTPRequestHeaderConfigCookie, "")
	if cookieName == "" {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("no_cookie")
		protectionStatus.AddToResponse(c)

		message := "Client IP not found"
		code := 500
		c.Status(code)
		c.Send(message)
		return "", errors.New(message)
	}
	return cookieName, nil
}

// GetOTPSecret get otp secret
func GetOTPSecret(c *fiber.Ctx) (OTPConfig, error) {
	secret := c.Get(HTTPRequestHeaderConfigOTPSecret, "")
	otpTime := c.Get(HTTPRequestHeaderConfigOTPTime, "")

	if ok := OtpSecretRegex.MatchString(secret); ok == false {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("otp_secret")
		protectionStatus.AddToResponse(c)

		message := "OTP Secret is invalid"
		code := 500
		c.Status(code)
		c.Send(message)
		return OTPConfig{}, errors.New(message)
	}

	u64, err := strconv.ParseUint(otpTime, 10, 32)
	if err != nil {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("otp_invalid_period")
		protectionStatus.AddToResponse(c)

		message := "Invalid period time"
		code := 500
		c.Status(code)
		c.Send(message)
		return OTPConfig{}, errors.New(message)
	}

	period := uint(u64)
	if period < 30 || period > 300 {
		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("error")
		protectionStatus.SetType("otp_period_time_range")
		protectionStatus.AddToResponse(c)

		message := "Invalid period time range"
		code := 500
		c.Status(code)
		c.Send(message)
		return OTPConfig{}, errors.New(message)
	}

	cnf := OTPConfig{}
	cnf.Period = period
	cnf.Secret = secret

	return cnf, nil
}

// ACLCheckAPIKeys is acl for check api keys
func ACLCheckAPIKeys(c *fiber.Ctx, ip string) (bool, error) {
	apiKeys := c.Get(HTTPRequestHeaderACLAPIKeys, "")
	clientAPIKey := c.Get(HTTPRequestHeaderClientAPIKey, "")
	if apiKeys != "" && clientAPIKey != "" {
		apiKeyObj := make(map[string]string)
		err := json.Unmarshal([]byte(apiKeys), &apiKeyObj)
		if err != nil {
			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("error")
			protectionStatus.SetType("invalid_api_keys")
			protectionStatus.AddToResponse(c)

			message := "Invalid API Key configuration"
			code := 500
			c.Status(code)
			c.Send(message)
			return false, errors.New(message)
		}
		for clientName, clientKey := range apiKeyObj {
			if clientKey == clientAPIKey {
				// prometheus
				PrometheusAccpetTotal.With(prometheus.Labels{"type": "api_key"})

				protectionStatus := NewProtectionStatus()
				protectionStatus.SetStatus("accept")
				protectionStatus.SetType("api_key")
				protectionStatus.AddExtra("client", clientName)
				protectionStatus.AddExtra("ip", ip)
				protectionStatus.AddToResponse(c)

				code := 200
				c.Status(code).Send(http.StatusText(code))
				return true, nil
			}
		}
	}
	return false, nil
}

// ACLCheckCountry is acl for check client country
func ACLCheckCountry(c *fiber.Ctx, ip string) (bool, error) {
	aclCountries := strings.Split(c.Get(HTTPRequestHeaderACLCountries, ""), ",")
	country := c.Get(HTTPRequestHeaderClientCountry, "")
	if country != "" && len(aclCountries) >= 1 {
		for _, code := range aclCountries {
			if country == code {
				// prometheus
				PrometheusAccpetTotal.With(prometheus.Labels{"type": "country"})

				protectionStatus := NewProtectionStatus()
				protectionStatus.SetStatus("accept")
				protectionStatus.SetType("country")
				protectionStatus.AddExtra("country", country)
				protectionStatus.AddExtra("ip", ip)
				protectionStatus.AddToResponse(c)

				code := 200
				c.Status(code).Send(http.StatusText(code))
				return true, nil
			}
		}
	}
	return false, nil
}

// ACLCheckASN is acl for check client asn
func ACLCheckASN(c *fiber.Ctx, ip string) (bool, error) {
	aclASNs := strings.Split(c.Get(HTTPRequestHeaderACLASNs, ""), ",")
	asn := c.Get(HTTPRequestHeaderClientASNNumber, "")
	asnOrg := c.Get(HTTPRequestHeaderClientASNOrganization, "")
	if asn != "" && asnOrg != "" && len(aclASNs) >= 1 {
		for _, num := range aclASNs {
			if asn == num {
				// prometheus
				PrometheusAccpetTotal.With(prometheus.Labels{"type": "asn"})

				protectionStatus := NewProtectionStatus()
				protectionStatus.SetStatus("accept")
				protectionStatus.SetType("asn")
				protectionStatus.AddExtra("asn", asn)
				protectionStatus.AddExtra("asnOrg", asnOrg)
				protectionStatus.AddToResponse(c)

				code := 200
				c.Status(code).Send(http.StatusText(code))
				return true, nil
			}
		}
	}
	return false, nil
}

// ACLCheckCIDR is acl for check client cidr
func ACLCheckCIDR(c *fiber.Ctx, ip string) (bool, error) {
	allowCIDRs := c.Get(HTTPRequestHeaderACLCIDRs, "")
	if len(allowCIDRs) >= 7 {
		allowCIDRsList := strings.Split(allowCIDRs, ",")
		ipB := net.ParseIP(ip)
		if ipB != nil {
			for _, CIDR := range allowCIDRsList {
				_, ipNetA, err := net.ParseCIDR(CIDR)
				if err == nil && ipNetA.Contains(ipB) {
					// prometheus
					PrometheusAccpetTotal.With(prometheus.Labels{"type": "cidr"})

					protectionStatus := NewProtectionStatus()
					protectionStatus.SetStatus("accept")
					protectionStatus.SetType("cidr")
					protectionStatus.AddExtra("ip", ip)
					protectionStatus.AddExtra("cidr", CIDR)
					protectionStatus.AddToResponse(c)

					code := 200
					c.Status(code).Send(http.StatusText(code))
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// NewProtectionStatus make new ProtectionStatus
func NewProtectionStatus() ProtectionStatus {
	p := ProtectionStatus{}
	p.User = ""
	p.Extras = make(map[string]string)

	return p
}

// SetStatus is set status of protection
func (p *ProtectionStatus) SetStatus(status string) {
	p.Status = status
}

// SetType is set type of protection
func (p *ProtectionStatus) SetType(typeOf string) {
	p.Type = typeOf
}

// SetUser is set user for protection status
func (p *ProtectionStatus) SetUser(user string) {
	p.User = user
}

// AddToResponse add status lines to response
func (p *ProtectionStatus) AddToResponse(c *fiber.Ctx) {
	c.Set(HTTPResponseHeaderProtectionStatus, p.Status)
	c.Set(HTTPResponseHeaderProtectionType, p.Type)
	if p.User != "" {
		c.Set(HTTPResponseHeaderProtectionUser, p.User)
	}

	extra := []string{}
	for k, v := range p.Extras {
		extra = append(extra, fmt.Sprintf("%s=%s", k, v))
	}
	c.Set(HTTPResponseHeaderProtectionExtra, strings.Join(extra[:], " "))
}

// AddExtra add extra data to protection status
func (p *ProtectionStatus) AddExtra(key string, value string) {
	p.Extras[key] = value
}
