package main

import (
	"testing"
	"time"

	aesGo "github.com/aasaam/aes-go"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestChallenge(t *testing.T) {
	secret := aesGo.GenerateKey()
	ch1 := newChallenge("fa", challengeTypeBlock, secret, secret, 1, 2, 3)
	ch1.setJSValue()
	ch1.CaptchaChallengeValue = 10000
	totpSecret := totpGenerate()
	ch1.setTOTPSecret(totpSecret)
	str, _ := ch1.getChallengeToken(secret)
	ch2, _ := newChallengeFromString(str, secret)
	if ch2.JSChallengeValue != ch1.JSChallengeValue {
		t.Errorf("invalid token check")
	}
	if !ch2.verifyJSValue(ch1.JSChallengeValue) {
		t.Errorf("invalid token js")
	}
	if !ch2.verifyCaptchaValue(ch1.CaptchaChallengeValue) {
		t.Errorf("invalid token js")
	}

	passCode, _ := totp.GenerateCodeCustom(ch1.TOTPSecret, time.Now(), totp.ValidateOpts{
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if !ch2.verifyTOTP(passCode) {
		t.Errorf("invalid token totp")
	}
}
