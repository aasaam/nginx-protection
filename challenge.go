package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	aesGo "github.com/aasaam/aes-go"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type captchaRequest struct {
	Lang            string `json:"lang,omitempty"`
	Quality         int    `json:"quality,omitempty"`
	DifficultyLevel string `json:"level,omitempty"`
}

type captchaResponse struct {
	Value int    `json:"value"`
	Image string `json:"image"`
}

type challenge struct {
	ID                      string `json:"i"`
	ClientTemporaryChecksum string `json:"c_t_ch"`
	ClientPersistChecksum   string `json:"c_p_ch"`
	ChallengeType           string `json:"ch_t"`
	Lang                    string `json:"l"`
	TimeInit                int64  `json:"t_i"`
	TimeWait                int64  `json:"t_w"`
	ExpireTime              int64  `json:"t_o"`
	TTL                     int64  `json:"ttl"`
	TOTPSecret              string `json:"totp_s"`
	JSChallengeValue        string `json:"js_ch_v"`
	CaptchaChallengeValue   int    `json:"ch_ch_v"`
}

func newChallenge(
	lang string,
	challengeType string,
	clientTemporaryChecksum string,
	clientPersistChecksum string,
	waitSeconds int64,
	timeoutSeconds int64,
	ttl int64,
) *challenge {
	ch := challenge{
		ID:                      randomBase64(),
		Lang:                    lang,
		ClientTemporaryChecksum: clientTemporaryChecksum,
		ClientPersistChecksum:   clientPersistChecksum,
		ChallengeType:           challengeType,
		TimeInit:                time.Now().Unix(),
		TimeWait:                time.Now().Unix() + waitSeconds,
		ExpireTime:              time.Now().Unix() + timeoutSeconds,
		TTL:                     ttl,
	}

	return &ch
}

func newChallengeFromString(tokenString string, secret string) (*challenge, error) {
	if tokenString == "" {
		return nil, errors.New("empty token data")
	}
	aes := aesGo.NewAasaamAES(secret)
	data := aes.DecryptTTL(tokenString)
	if data == "" {
		return nil, errors.New("invalid token data")
	}
	var ch challenge
	err := json.Unmarshal([]byte(data), &ch)
	if err != nil {
		return nil, errors.New("invalid token structure")
	}
	return &ch, nil
}

func (ch *challenge) verify(clientTemporaryChecksum string) bool {
	now := time.Now().Unix()
	return len(clientTemporaryChecksum) > 16 && ch.ClientTemporaryChecksum == clientTemporaryChecksum && now >= ch.TimeWait && now <= ch.ExpireTime
}

func (ch *challenge) setJSValue() string {
	value := randomBase64()
	jsEval := fmt.Sprintf(`window.jsv=function(){return '%s'};`, value)
	jsEvalEncode := base64.StdEncoding.EncodeToString([]byte(jsEval))
	ch.JSChallengeValue = value
	return jsEvalEncode
}

func (ch *challenge) setTOTPSecret(totpSecret string) {
	ch.TOTPSecret = totpSecret
}

func (ch *challenge) setCaptchaValue(restCaptchaURL string, difficultyLevel string) (string, error) {
	r := captchaRequest{
		Lang:            ch.Lang,
		Quality:         1,
		DifficultyLevel: difficultyLevel,
	}
	jsonBytes, _ := json.Marshal(r)
	req, err0 := http.NewRequest("POST", restCaptchaURL+"/new", bytes.NewBuffer(jsonBytes))
	if err0 != nil {
		return "", errors.New("failed to init request captcha server")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err1 := client.Do(req)
	if err1 != nil {
		return "", errors.New("cannot request to captcha server")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var captchaRes captchaResponse
	err2 := json.Unmarshal(body, &captchaRes)
	if err2 != nil {
		return "", errors.New("captcha server response invalid data")
	}

	if captchaRes.Value < 1 {
		return "", errors.New("captcha server value not return")
	}

	ch.CaptchaChallengeValue = captchaRes.Value

	return captchaRes.Image, nil
}

func (ch *challenge) verifyJSValue(v string) bool {
	return v != "" && len(v) >= 32 && v == ch.JSChallengeValue
}

func (ch *challenge) verifyCaptchaValue(v int) bool {
	return v > 9 && v < 9999999999 && v == ch.CaptchaChallengeValue
}

func (ch *challenge) verifyTOTP(v string) bool {
	passCode, err := totp.GenerateCodeCustom(ch.TOTPSecret, time.Now(), totp.ValidateOpts{
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && passCode == v
}

func (ch *challenge) getChallengeToken(secret string) (string, error) {
	aes := aesGo.NewAasaamAES(secret)
	byteData, err := json.Marshal(ch)
	if err != nil {
		return "", err
	}
	return aes.EncryptTTL(string(byteData), ch.TTL), nil
}
