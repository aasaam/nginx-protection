package main

import (
	"crypto/rsa"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
)

// RouterAuth check for validation
func RouterAuth(
	config *Config,
	router fiber.Router,
	rsa *rsa.PrivateKey,
) {
	router.All("/auth", func(c *fiber.Ctx) {
		// prometheus
		PrometheusRequestTotal.Inc()

		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		cookieName, err := GetCookieName(c)
		if err != nil {
			return
		}

		aclAPIKeyPassed, err := ACLCheckAPIKeys(c, ip)
		if err != nil || aclAPIKeyPassed {
			return
		}

		aclCountryPass, err := ACLCheckCountry(c, ip)
		if err != nil || aclCountryPass {
			return
		}

		aclCIDRPass, err := ACLCheckCIDR(c, ip)
		if err != nil || aclCIDRPass {
			return
		}

		aclASNPass, err := ACLCheckASN(c, ip)
		if err != nil || aclASNPass {
			return
		}

		clientCheckSum := GetClientCheckSum(c)

		authToken, err := ValidateAuthToken(c.Cookies(cookieName), config.Salt, rsa)

		if err == nil && authToken.Checksum == clientCheckSum {
			// prometheus
			PrometheusAccpetTotal.With(prometheus.Labels{"type": authToken.Type}).Inc()

			protectionStatus := NewProtectionStatus()
			protectionStatus.SetStatus("accept")
			protectionStatus.SetType("token")
			if authToken.User != "" {
				protectionStatus.SetUser(authToken.User)
			}
			protectionStatus.AddToResponse(c)

			code := 200
			c.Status(code).Send(http.StatusText(code))
			return
		}

		statusCodeHeader := c.Get(HTTPRequestHeaderConfigUnauthorizedStatus, "401")
		statusCode, _ := strconv.ParseInt(statusCodeHeader, 10, 32)

		// prometheus
		PrometheusRejectTotal.Inc()
		code := int(statusCode)
		c.Status(code).Send(http.StatusText(401))
	})
}
