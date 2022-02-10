package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	aesGo "github.com/aasaam/aes-go"
	"github.com/gofiber/fiber/v2"
)

func TestHTTPTest01(t *testing.T) {
	// variables
	tokenSecret := aesGo.GenerateKey()
	clientSecret := aesGo.GenerateKey()
	config := newConfig("fatal", false, "en", "en,fa", tokenSecret, clientSecret, "/.well-known/protection", "", "", "")
	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()
	clientPersistChecksum := aesGo.GenerateKey()
	clientTemporaryChecksum := aesGo.GenerateKey()
	ip := "1.1.1.1"

	// http server
	httpApp := newHTTPServer(config, challengeStorage, aclStorage)

	// misconfigure: X-Forwarded-For
	req00 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req00.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req00.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req00.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	resp00, _ := httpApp.Test(req00)
	if resp00.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: X-Forwarded-For
	req001 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req001.Header.Set("X-Forwarded-For", "a, c, d")
	req001.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req001.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req001.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	resp001, _ := httpApp.Test(req001)
	if resp001.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: clientPersistChecksum
	req01 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req01.Header.Set("X-Forwarded-For", ip)
	req01.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req01.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	resp01, _ := httpApp.Test(req01)
	if resp01.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: clientTemporaryChecksum
	req02 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req02.Header.Set("X-Forwarded-For", ip)
	req02.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req02.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	resp02, _ := httpApp.Test(req02)
	if resp02.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: challengeType
	req03 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req03.Header.Set("X-Forwarded-For", ip)
	req03.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req03.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	resp03, _ := httpApp.Test(req03)
	if resp03.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: totp secret
	req04 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req04.Header.Set("X-Forwarded-For", ip)
	req04.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req04.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req04.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeTOTP)
	resp04, _ := httpApp.Test(req04)
	if resp04.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

	// misconfigure: request id
	req05 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req05.Header.Set("X-Forwarded-For", ip)
	req05.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req05.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req05.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeTOTP)
	req05.Header.Set(httpRequestHeaderConfigTOTPSecret, totpGenerate())
	resp05, _ := httpApp.Test(req05)
	if resp05.StatusCode != misconfigureStatus {
		t.Errorf("must response")
	}

}

func TestHTTPTest02(t *testing.T) {
	// variables
	tokenSecret := aesGo.GenerateKey()
	clientSecret := aesGo.GenerateKey()
	config := newConfig("fatal", false, "en", "en,fa", tokenSecret, clientSecret, "/.well-known/protection", "", "", "")
	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()
	clientPersistChecksum := aesGo.GenerateKey()
	clientTemporaryChecksum := aesGo.GenerateKey()
	ip := "1.1.1.1"

	// http server
	httpApp := newHTTPServer(config, challengeStorage, aclStorage)

	// cidr
	req1 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req1.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req1.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req1.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req1.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	// configure
	req1.Header.Set("X-Forwarded-For", "192.168.1.1, 1.1.1.1")
	req1.Header.Set(httpRequestHeaderACLCIDRs, "127.0.0.0/8, invalid, 192.168.0.0/16")
	// for acl cache
	httpApp.Test(req1)
	resp1, _ := httpApp.Test(req1)
	if resp1.StatusCode != 200 {
		t.Errorf("must authorized cidr")
	}

	// country
	req2 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req2.Header.Set("X-Forwarded-For", ip)
	req2.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req2.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req2.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req2.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	// configure
	req2.Header.Set(httpRequestHeaderClientCountry, "CN")
	req2.Header.Set(httpRequestHeaderACLCountries, "DE, IR, CN, 678, AA")
	resp2, _ := httpApp.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("must authorized country")
	}

	// asn
	req3 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req3.Header.Set("X-Forwarded-For", ip)
	req3.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req3.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req3.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req3.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	// configure
	req3.Header.Set(httpRequestHeaderClientASNNumber, "1000")
	req3.Header.Set(httpRequestHeaderACLASNs, "2000,10,aa , 2000, 1000, 99")
	resp3, _ := httpApp.Test(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("must authorized country")
	}

	// asn ranges
	req4 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req4.Header.Set("X-Forwarded-For", ip)
	req4.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req4.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req4.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req4.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	// configure
	req4.Header.Set(httpRequestHeaderClientASNNumber, "1000")
	req4.Header.Set(httpRequestHeaderACLASNRanges, "100-200, aa-100, 200-cc, 1000-1100")
	resp4, _ := httpApp.Test(req4)
	if resp4.StatusCode != 200 {
		t.Errorf("must authorized country")
	}

	// api key
	req5 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req5.Header.Set("X-Forwarded-For", ip)
	req5.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req5.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req5.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req5.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	// configure
	req5.Header.Set(httpRequestHeaderClientAPIKeyHeaderName, "v2")
	req5.Header.Set(httpRequestHeaderACLAPIKeys, `{"c1":"v1", "c2":"v2"}`)
	resp5, _ := httpApp.Test(req5)
	if resp5.StatusCode != 200 {
		t.Errorf("must authorized country")
	}
}

func TestHTTPTest03(t *testing.T) {
	// variables
	tokenSecret := aesGo.GenerateKey()
	clientSecret := aesGo.GenerateKey()
	config := newConfig("fatal", false, "en", "en,fa", tokenSecret, clientSecret, "/.well-known/protection", "", "", "")
	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()
	clientPersistChecksum := aesGo.GenerateKey()
	clientTemporaryChecksum := aesGo.GenerateKey()
	ip := "1.1.1.1, aa , 8.8.8.8"

	// http server
	httpApp := newHTTPServer(config, challengeStorage, aclStorage)

	// block
	req1 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req1.Header.Set("X-Forwarded-For", ip)
	req1.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req1.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req1.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req1.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
	resp1, _ := httpApp.Test(req1)
	if resp1.StatusCode != 200 {
		t.Errorf("must response")
	}

	// js
	req2 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req2.Header.Set("X-Forwarded-For", ip)
	req2.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req2.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req2.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req2.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req2.Header.Set(httpRequestHeaderConfigWait, "0")
	resp2, _ := httpApp.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("must response")
	}
	tokenString2 := resp2.Header.Get(httpResponseChallengeToken)
	time.Sleep(time.Second * 3)
	chReq2 := challengeRequest{
		ChallengeToken: tokenString2,
	}
	chReq2JSON, _ := json.Marshal(chReq2)

	req22 := httptest.NewRequest("PATCH", "/.well-known/protection/challenge", bytes.NewBuffer(chReq2JSON))
	req22.Header.Set("X-Forwarded-For", ip)
	req22.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req22.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req22.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req22.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req22.Header.Set(httpRequestHeaderConfigWait, "0")
	req22.Header.Set(fiber.HeaderContentType, "application/json")
	resp22, _ := httpApp.Test(req22)
	if resp22.StatusCode != 200 {
		t.Errorf("must response")
		return
	}
	tokenString22 := resp22.Header.Get(httpResponseChallengeToken)
	challenge22, _ := newChallengeFromString(tokenString22, clientSecret)
	chReq22 := challengeRequest{
		ChallengeToken: tokenString22,
		JSValue:        challenge22.JSChallengeValue,
	}
	chReq22JSON, errChReq22JSON := json.Marshal(chReq22)
	if errChReq22JSON != nil {
		t.Error(errChReq22JSON)
	}
	req222 := httptest.NewRequest("POST", "/.well-known/protection/challenge", bytes.NewBuffer(chReq22JSON))
	req222.Header.Set("X-Forwarded-For", ip)
	req222.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req222.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req222.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req222.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req222.Header.Set(httpRequestHeaderConfigWait, "0")
	req222.Header.Set(fiber.HeaderContentType, "application/json")
	resp222, _ := httpApp.Test(req222)
	if resp222.StatusCode != 200 {
		t.Errorf("must response")
	}
	resultValue222 := resp222.Header.Get(httpResponseChallengeResult)

	// using cookie
	req4 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req4.Header.Set("X-Forwarded-For", ip)
	req4.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req4.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req4.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req4.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	cookie := http.Cookie{
		Name:  defaultCookieName,
		Value: resultValue222,
	}
	req4.AddCookie(&cookie)
	resp4, _ := httpApp.Test(req4)
	if resp4.StatusCode != 200 {
		t.Errorf("must authorized country")
	}
}
func TestHTTPTest04(t *testing.T) {
	// variables
	tokenSecret := aesGo.GenerateKey()
	clientSecret := aesGo.GenerateKey()
	config := newConfig("fatal", false, "en", "en,fa", tokenSecret, clientSecret, "/.well-known/protection", "", "", "")
	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()
	clientPersistChecksum := aesGo.GenerateKey()
	clientTemporaryChecksum := aesGo.GenerateKey()
	ip := "1.1.1.1, aa , 8.8.8.8"

	// http server
	httpApp := newHTTPServer(config, challengeStorage, aclStorage)

	// js
	req2 := httptest.NewRequest("GET", "/.well-known/protection/challenge", nil)
	req2.Header.Set("X-Forwarded-For", ip)
	req2.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req2.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req2.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req2.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req2.Header.Set(httpRequestHeaderConfigWait, "0")
	resp2, _ := httpApp.Test(req2)
	if resp2.StatusCode != 200 {
		t.Errorf("must response")
		return
	}
	tokenString2 := resp2.Header.Get(httpResponseChallengeToken)
	time.Sleep(time.Second * 3)
	chReq2 := challengeRequest{
		ChallengeToken: tokenString2,
	}
	chReq2JSON, _ := json.Marshal(chReq2)

	req22 := httptest.NewRequest("PATCH", "/.well-known/protection/challenge", bytes.NewBuffer(chReq2JSON))
	req22.Header.Set("X-Forwarded-For", ip)
	req22.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req22.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req22.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req22.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req22.Header.Set(httpRequestHeaderConfigWait, "0")
	req22.Header.Set(fiber.HeaderContentType, "application/json")
	resp22, _ := httpApp.Test(req22)
	if resp22.StatusCode != 200 {
		t.Errorf("must response")
		return
	}
	tokenString22 := resp22.Header.Get(httpResponseChallengeToken)
	challenge22, _ := newChallengeFromString(tokenString22, clientSecret)
	chReq22 := challengeRequest{
		ChallengeToken: tokenString22,
		JSValue:        challenge22.JSChallengeValue,
	}
	chReq22JSON, _ := json.Marshal(chReq22)
	req222 := httptest.NewRequest("POST", "/.well-known/protection/challenge", bytes.NewBuffer(chReq22JSON))
	req222.Header.Set("X-Forwarded-For", ip)
	req222.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req222.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req222.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req222.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	req222.Header.Set(httpRequestHeaderConfigWait, "0")
	req222.Header.Set(fiber.HeaderContentType, "application/json")
	resp222, _ := httpApp.Test(req222)
	if resp222.StatusCode != 200 {
		t.Errorf("must response")
	}
	resultValue222 := resp222.Header.Get(httpResponseChallengeResult)
	// using cookie
	req4 := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
	req4.Header.Set("X-Forwarded-For", ip)
	req4.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
	req4.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
	req4.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
	req4.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeJS)
	cookie := http.Cookie{
		Name:  defaultCookieName,
		Value: resultValue222,
	}
	req4.AddCookie(&cookie)
	resp4, _ := httpApp.Test(req4)
	if resp4.StatusCode != 200 {
		t.Errorf("must authorized country")
	}
}

func BenchmarkACLStorageOnHTTP(b *testing.B) {
	// variables
	tokenSecret := aesGo.GenerateKey()
	clientSecret := aesGo.GenerateKey()
	clientPersistChecksum := aesGo.GenerateKey()
	clientTemporaryChecksum := aesGo.GenerateKey()
	config := newConfig("fatal", false, "en", "en,fa", tokenSecret, clientSecret, "/.well-known/protection", "", "", "")
	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()
	ip := "192.168.1.1"

	// http server
	httpApp := newHTTPServer(config, challengeStorage, aclStorage)

	for i := 0; i < b.N; i++ {
		// asn ranges
		req := httptest.NewRequest("GET", "/.well-known/protection/auth", nil)
		req.Header.Set("X-Forwarded-For", ip)
		req.Header.Set(httpRequestHeaderRequestID, aesGo.GenerateKey())
		req.Header.Set(httpRequestHeaderClientPersistChecksum, clientPersistChecksum)
		req.Header.Set(httpRequestHeaderClientTemporaryChecksum, clientTemporaryChecksum)
		req.Header.Set(httpRequestHeaderConfigChallenge, challengeTypeBlock)
		req.Header.Set(httpRequestHeaderClientASNNumber, "1000")
		req.Header.Set(httpRequestHeaderACLASNRanges, "10-100,1000-1100")
		resp, _ := httpApp.Test(req)
		if resp.StatusCode != 200 {
			b.Error("must valid")
		}
	}
}
