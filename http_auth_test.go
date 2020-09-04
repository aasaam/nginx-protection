package main

import (
	"net/http/httptest"
	"testing"
)

const pv10 = `-----BEGIN PRIVATE KEY-----
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

func TestGetHTTPServerAuth(t *testing.T) {
	cookieName := "prt"
	asnOrg := "Google LLC."

	aclAPIKeys := "{\"foo\":\"bar\"}"
	aclCountries := "IR,DE"
	aclCIDR := "127.0.0.1/8"
	aclASN := "36040,36041"

	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv10))
	app := GetHTTPServer(&config, rsa)

	req1 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req1.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req1.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req1.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req1.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req1.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req1.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req1.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req1.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req1.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 401 {
		t.Errorf("Status must be 401")
	}

	req2 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req2.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req2.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req2.Header.Set(HTTPRequestHeaderClientAPIKey, "bar")
	req2.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req2.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req2.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req2.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req2.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req2.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req2.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req3 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req3.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req3.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req3.Header.Set(HTTPRequestHeaderClientAPIKey, "bar")
	req3.Header.Set(HTTPRequestHeaderClientCountry, "IR")
	req3.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req3.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req3.Header.Set(HTTPRequestHeaderACLAPIKeys, "")
	req3.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req3.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req3.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req4 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req4.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req4.Header.Set(HTTPRequestXForwardedFor, "127.0.0.1")
	req4.Header.Set(HTTPRequestHeaderClientAPIKey, "bar")
	req4.Header.Set(HTTPRequestHeaderClientCountry, "CB")
	req4.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req4.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req4.Header.Set(HTTPRequestHeaderACLAPIKeys, "")
	req4.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req4.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req4.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp4, _ := app.Test(req4)
	if resp4.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	req5 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req5.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req5.Header.Set(HTTPRequestXForwardedFor, "192.168.1.1")
	req5.Header.Set(HTTPRequestHeaderClientAPIKey, "bar")
	req5.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req5.Header.Set(HTTPRequestHeaderClientASNNumber, "36040")
	req5.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req5.Header.Set(HTTPRequestHeaderACLAPIKeys, "")
	req5.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req5.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req5.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp5, _ := app.Test(req5)
	if resp5.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
}
func TestGetHTTPServerAuthErrors(t *testing.T) {
	cookieName := "prt"
	asnOrg := "Google LLC."

	aclAPIKeys := "{\"foo\":\"bar\"}"
	aclCountries := "IR,DE"
	aclCIDR := "127.0.0.1/8"
	aclASN := "36040,36041"

	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	rsa, _ := LoadPrivateKey([]byte(pv10))
	app := GetHTTPServer(&config, rsa)

	req1 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req1.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req1.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req1.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req1.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req1.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req1.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req1.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req1.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}

	req2 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req2.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req2.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req2.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req2.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req2.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req2.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req2.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 500 {
		t.Errorf("Status must be 500")
	}
}

func TestGetHTTPServerAuthToken(t *testing.T) {
	cookieName := "prt"
	asnOrg := "Google LLC."

	aclAPIKeys := "{\"foo\":\"bar\"}"
	aclCountries := "IR,DE"
	aclCIDR := "127.0.0.1/8"
	aclASN := "36040,36041"

	config := Config{}
	config.BaseURL = "/.well-known/protection"
	config.Salt = "0123456789"

	clientCheckSum := MD5("qw123456ew")

	challenge := ChallengeToken{}
	challenge.Type = ChallengeTypeJS

	authToken := AuthToken{}
	authToken.Checksum = clientCheckSum
	authToken.Type = challenge.Type

	rsa, _ := LoadPrivateKey([]byte(pv10))
	app := GetHTTPServer(&config, rsa)

	tokenValue, _ := GenerateAuthToken(authToken, 100, config.Salt, rsa)

	req1 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req1.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req1.Header.Set("Cookie", cookieName+"="+tokenValue)
	req1.Header.Set(HTTPRequestHeaderClientChecksum, "qw123456ew")
	req1.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req1.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req1.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req1.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req1.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req1.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req1.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req1.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}

	authToken2 := AuthToken{}
	authToken2.Checksum = clientCheckSum
	authToken2.Type = challenge.Type
	authToken2.User = "usr"

	tokenValue2, _ := GenerateAuthToken(authToken2, 100, config.Salt, rsa)

	req2 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req2.Header.Set(HTTPRequestXForwardedFor, "1.1.1.1")
	req2.Header.Set("Cookie", cookieName+"="+tokenValue2)
	req2.Header.Set(HTTPRequestHeaderClientChecksum, "qw123456ew")
	req2.Header.Set(HTTPRequestHeaderConfigCookie, cookieName)
	req2.Header.Set(HTTPRequestHeaderClientCountry, "CN")
	req2.Header.Set(HTTPRequestHeaderClientASNNumber, "11223")
	req2.Header.Set(HTTPRequestHeaderClientASNOrganization, asnOrg)
	req2.Header.Set(HTTPRequestHeaderACLAPIKeys, aclAPIKeys)
	req2.Header.Set(HTTPRequestHeaderACLCountries, aclCountries)
	req2.Header.Set(HTTPRequestHeaderACLASNs, aclASN)
	req2.Header.Set(HTTPRequestHeaderACLCIDRs, aclCIDR)

	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("Status must be 200")
	}
}
