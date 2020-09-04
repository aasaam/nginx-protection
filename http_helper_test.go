package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"
)

const pv1 = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDHZ+SKx2uhfrai
EmQI8SNNzlWnTBSXMJzvMHBnogcNYTMsEUOa9hYWRrKvWXQMubfVW41vxhMblH9b
ozChsJKYhIXPmXxLMmhfTUbh97UsOxBBTu54ZP9bwz1SxHdpOjIDfuAetiblWLiL
omYuh9DUpOm3F+l1OStPDfTdQNYhXt5cRE//5y7DLuvSlingHsprTvO84ia/nIBP
5ost5hLQ80Jyaj3SBOOAjYEw/R0fIt9tl5e5kfSe5eMjzyX/v55HlAkDv8NK1gb4
bpURHLsC9UUHUxeKIRQsfXhzzTuL8YCyyV7KG3NgvW9ZbVnCZZ3nmELMsgad5pEx
Z4pUENL7AgMBAAECggEBALaZ7vEe+PLkRH5Z9P0zRK8FWe5ffyOMQsnOQ8DC4U5h
Sij6jjwjScqQZySn99uHXk6lDfnjGrBQ5eeWovwN49CC2r5mwSljOay76UMYQPIG
DDah/0KEykrPmSJoAyl7Pz1wO/Ajwa6X9jb4OjY17QgtFFC0Nvc/qOc10puhufTH
ePwhknHTOH+Y5XAyiMNp8gJoAXQ8NHh/6eJW3W0BIHf/HXKfsucLwVkz2412y/Jt
GrFZbRUUnHYSyJ/VjOfUwHAsmWR77YNh3EL12GKEF4u5OGA4zhwWeACIzj13wJa6
V0QgkpmYvXiTxtHwlLjAE3EznwRbs5YkR12h4/5UMmECgYEA++YFNkbOoTslH5F7
Y8NuzfjYzAxdZIKk1XsOOVpqkeas0euEvx539Qkgnt2RK6svtrFr5PCi19nf0Olc
csVDrzLnnv2S2kpI/JOwopxa0FFvdhHoe+wno0Rf8lnzY62jt61s66bRXJpk98IH
s7twYe2USwj8UD7ORLb8zpbbcIMCgYEAyqcRqPvhdMC9ewgBIhJV/c+FeComweTW
SFPUsd9dPsIg4OCXEX0Mq2xSJvRzAne95tZGrmFsr+pRQ3OXGaZdScYExD86OBUF
MpKY9T6J7ODVkUm31zjF631Bvp61py2OWbgTjguD73lLL7l9qVTnyb66BpeQyfon
E77Bk1M/mikCgYA3TEiqoKKtzGEa7AINZZLGjrFxIenCrddnsgruVkX834niz3Ql
zJeC6E0L8xHyZzMjRRGtgZIOFptGrmQIIfv40xD72yjI2PPq1rU5DV/2SVpRrh6+
TZpqAhGaD1sZ7714Dg9SMB3X2WD+7s5oC2bhaJlcW42gRBleBlm7NGzZ5wKBgD5R
G8wkEJNvhZTkxDxu+QSAoSFvjNWJAh/hr4E3F5xp4+RjC/Fzy8aXG7gg6ZDzs3Dd
qYSMLvj1jCG61NctYniCLQsQCl4ekKeZjvGzVoSCKwpvadoD+lDNBr+QXHnZN3H9
ef3vKpYkbWtyleLRWimeveOzDfIeO5AF087zBZbpAoGABFeS5KjjYI6C2Jsdhh+r
Z4qVKLTGcyeLL30lDmCGVZZvbSN2/3rW3/vaf5/Bm42wVHTPnBN2IBFTH5gv99qc
xUORO3LUFlDzX/WSDDpUF5RbT3eE8IFrLAH0ouazyfR9KhV9w5qzhy8cGGEE+3gw
J6msr7d0ebh9hhGZpdGzifA=
-----END PRIVATE KEY-----`

func TestGetHTTPServer(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	req1 := httptest.NewRequest("GET", "/not-found", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 404 {
		t.Errorf("Status must be 404")
	}

	req2 := httptest.NewRequest("GET", "/metrics", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
}
func TestHTTPHelpersGetRequestIP(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/ip", func(c *fiber.Ctx) {
		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("ok")
		protectionStatus.SetType("type")
		protectionStatus.SetUser("usr")
		protectionStatus.AddExtra("extra", "ok")
		protectionStatus.AddToResponse(c)

		GetClientCheckSum(c)
		GetClientTokenCheckSum(c)

		c.Send(ip)
	})

	req1 := httptest.NewRequest("GET", "/ip", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/ip", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 500")
	}
	req3 := httptest.NewRequest("GET", "/ip", nil)
	req3.Header.Set(HTTPRequestXForwardedFor, "192.168.1.1,127.0.0.1")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("Status must be 500")
	}
	req4 := httptest.NewRequest("GET", "/ip", nil)
	req4.Header.Set(HTTPRequestXForwardedFor, "192.168.1.2,  127.0.0.2")
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 200 {
		t.Errorf("Status must be 500")
	}
	req5 := httptest.NewRequest("GET", "/ip", nil)
	req5.Header.Set(HTTPRequestXForwardedFor, "foo-bar")
	resp5, _ := app.Test(req5)
	if resp5.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}
func TestHTTPHelpersGetFarsiCaptcha(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/farsicaptcha", func(c *fiber.Ctx) {
		isFarsiCaptcha, err := GetFarsiCaptcha(c)
		if err != nil {
			return
		}

		protectionStatus := NewProtectionStatus()
		protectionStatus.SetStatus("ok")
		protectionStatus.SetType("type")
		protectionStatus.AddExtra("extra", "ok")
		protectionStatus.AddToResponse(c)

		c.Send(isFarsiCaptcha)
	})

	req1 := httptest.NewRequest("GET", "/farsicaptcha", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/farsicaptcha", nil)
	req2.Header.Set(HTTPRequestHeaderConfigFarsiCaptcha, "true")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req3 := httptest.NewRequest("GET", "/farsicaptcha", nil)
	req3.Header.Set(HTTPRequestHeaderConfigFarsiCaptcha, "0")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
}

func TestHTTPHelpersGetChallengeType(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/challenge-type", func(c *fiber.Ctx) {
		challengeType, err := GetChallengeType(c)
		if err != nil {
			return
		}

		c.Send(challengeType)
	})

	req1 := httptest.NewRequest("GET", "/challenge-type", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/challenge-type", nil)
	req2.Header.Set(HTTPRequestHeaderConfigChallenge, ChallengeTypeJS)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req3 := httptest.NewRequest("GET", "/challenge-type", nil)
	req3.Header.Set(HTTPRequestHeaderConfigChallenge, "0")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}

func TestHTTPHelpersGetLanguage(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/lang", func(c *fiber.Ctx) {
		lang, err := GetLanguage(c)
		if err != nil {
			return
		}

		c.Send(lang)
	})

	req1 := httptest.NewRequest("GET", "/lang", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}

	req2 := httptest.NewRequest("GET", "/lang", nil)
	req2.Header.Set(HTTPRequestHeaderConfigLang, "fa")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req22 := httptest.NewRequest("GET", "/lang?lang=fa", nil)
	req22.Header.Set(HTTPRequestHeaderConfigLang, "fa")
	resp22, _ := app.Test(req22)
	if resp22.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/lang", nil)
	req3.Header.Set(HTTPRequestHeaderConfigLang, "zz")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}
func TestHTTPHelpersGetWaitTime(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/wait", func(c *fiber.Ctx) {
		waitSeconds, err := GetWaitTime(c)
		if err != nil {
			return
		}

		c.Send(waitSeconds)
	})

	req1 := httptest.NewRequest("GET", "/wait", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/wait", nil)
	req2.Header.Set(HTTPRequestHeaderConfigWait, "120")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req3 := httptest.NewRequest("GET", "/wait", nil)
	req3.Header.Set(HTTPRequestHeaderConfigWait, "zz")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req4 := httptest.NewRequest("GET", "/wait", nil)
	req4.Header.Set(HTTPRequestHeaderConfigWait, "1")
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}

func TestHTTPHelpersGetTimeout(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/timeout", func(c *fiber.Ctx) {
		timeoutSeconds, err := GetTimeout(c)
		if err != nil {
			return
		}

		c.Send(timeoutSeconds)
	})

	req1 := httptest.NewRequest("GET", "/timeout", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/timeout", nil)
	req2.Header.Set(HTTPRequestHeaderConfigTimeout, "120")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req3 := httptest.NewRequest("GET", "/timeout", nil)
	req3.Header.Set(HTTPRequestHeaderConfigTimeout, "zz")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req4 := httptest.NewRequest("GET", "/timeout", nil)
	req4.Header.Set(HTTPRequestHeaderConfigTimeout, "1")
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}

func TestHTTPHelpersGetChallengeTTL(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/ttl", func(c *fiber.Ctx) {
		ttlSeconds, err := GetChallengeTTL(c)
		if err != nil {
			return
		}

		c.Send(ttlSeconds)
	})

	req1 := httptest.NewRequest("GET", "/ttl", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/ttl", nil)
	req2.Header.Set(HTTPRequestHeaderConfigTTL, "120")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
	req3 := httptest.NewRequest("GET", "/ttl", nil)
	req3.Header.Set(HTTPRequestHeaderConfigTTL, "zz")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req4 := httptest.NewRequest("GET", "/ttl", nil)
	req4.Header.Set(HTTPRequestHeaderConfigTTL, "1")
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}
func TestHTTPHelpersGetCookieName(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/cookie-name", func(c *fiber.Ctx) {
		cookieName, err := GetCookieName(c)
		if err != nil {
			return
		}

		c.Send(cookieName)
	})

	req1 := httptest.NewRequest("GET", "/cookie-name", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
	req2 := httptest.NewRequest("GET", "/cookie-name", nil)
	req2.Header.Set(HTTPRequestHeaderConfigCookie, "prt")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
}
func TestHTTPHelpersACLCheckAPIKeys(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/acl", func(c *fiber.Ctx) {
		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		aclAPIKeyPassed, err := ACLCheckAPIKeys(c, ip)
		if err != nil || aclAPIKeyPassed {
			return
		}

		code := 403
		c.Status(code).Send(http.StatusText(code))
	})

	req1 := httptest.NewRequest("GET", "/acl", nil)
	req1.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}

	req2 := httptest.NewRequest("GET", "/acl", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req2.Header.Set(HTTPRequestHeaderClientAPIKey, "foobar")
	req2.Header.Set(HTTPRequestHeaderACLAPIKeys, "{\"client_a\":\"foobar\"}")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/acl", nil)
	req3.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req3.Header.Set(HTTPRequestHeaderClientAPIKey, "foobarnot")
	req3.Header.Set(HTTPRequestHeaderACLAPIKeys, "{\"client_a\":\"foobar\"}")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}

	req4 := httptest.NewRequest("GET", "/acl", nil)
	req4.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req4.Header.Set(HTTPRequestHeaderClientAPIKey, "foobarnot")
	req4.Header.Set(HTTPRequestHeaderACLAPIKeys, "__ * (((")
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}
func TestHTTPHelpersACLCheckCountry(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/country", func(c *fiber.Ctx) {
		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		aclCountryPass, err := ACLCheckCountry(c, ip)
		if err != nil || aclCountryPass {
			return
		}

		code := 403
		c.Status(code).Send(http.StatusText(code))
	})

	req1 := httptest.NewRequest("GET", "/country", nil)
	req1.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}

	req2 := httptest.NewRequest("GET", "/country", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req2.Header.Set(HTTPRequestHeaderClientCountry, "IR")
	req2.Header.Set(HTTPRequestHeaderACLCountries, "IR,US")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/country", nil)
	req3.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req3.Header.Set(HTTPRequestHeaderClientCountry, "DE")
	req3.Header.Set(HTTPRequestHeaderACLCountries, "IR,US")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}
}
func TestHTTPHelpersACLCheckASN(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/asn", func(c *fiber.Ctx) {
		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		aclASNPass, err := ACLCheckASN(c, ip)
		if err != nil || aclASNPass {
			return
		}

		code := 403
		c.Status(code).Send(http.StatusText(code))
	})

	req1 := httptest.NewRequest("GET", "/asn", nil)
	req1.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}

	req2 := httptest.NewRequest("GET", "/asn", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req2.Header.Set(HTTPRequestHeaderClientASNNumber, "36040")
	req2.Header.Set(HTTPRequestHeaderClientASNOrganization, "Google LLC.")
	req2.Header.Set(HTTPRequestHeaderACLASNs, "36040,36041")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/asn", nil)
	req3.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req3.Header.Set(HTTPRequestHeaderClientASNNumber, "36040")
	req3.Header.Set(HTTPRequestHeaderClientASNOrganization, "Google LLC.")
	req3.Header.Set(HTTPRequestHeaderACLASNs, "36041,36042")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}
}

func TestHTTPHelpersACLCheckCIDR(t *testing.T) {
	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv1))
	app := GetHTTPServer(&config, rsa)

	app.Get("/cidr", func(c *fiber.Ctx) {
		ip, err := GetRequestIP(c)
		if err != nil {
			return
		}

		aclCIDRPass, err := ACLCheckCIDR(c, ip)
		if err != nil || aclCIDRPass {
			return
		}

		code := 403
		c.Status(code).Send(http.StatusText(code))
	})

	req1 := httptest.NewRequest("GET", "/cidr", nil)
	req1.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}

	req2 := httptest.NewRequest("GET", "/cidr", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req2.Header.Set(HTTPRequestHeaderACLCIDRs, "127.0.0.0/8")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/cidr", nil)
	req3.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req3.Header.Set(HTTPRequestHeaderACLCIDRs, "192.0.0.0/8")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 403 {
		t.Errorf("Status must be 403")
	}
}
