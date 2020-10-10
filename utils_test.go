package main

import (
	"testing"
)

func TestMD5(t *testing.T) {
	hash := MD5("1")
	if hash != "c4ca4238a0b923820dcc509a6f75849b" {
		t.Errorf("Invalid md5")
	}
}

func TestGenerateOTPSecret(t *testing.T) {
	GenerateOTPSecret()
}

func TestRandomHex(t *testing.T) {
	v1 := RandomHex()
	v2 := RandomHex()
	if v1 == v2 {
		t.Errorf("RandomHex must be uniq")
	}
}
func TestRandomCaptchaLetters(t *testing.T) {
	v1 := RandomCaptchaLetters()
	v2 := RandomCaptchaLetters()
	if v1 == v2 {
		t.Errorf("RandomCaptchaLetters must be uniq")
	}
}

func TestIsValidLanguage(t *testing.T) {
	true1 := IsValidLanguage("en")
	true2 := IsValidLanguage("fa")
	if true1 != true || true2 != true {
		t.Errorf("Valid langauge must be true")
	}
	false1 := IsValidLanguage("xx")
	false2 := IsValidLanguage("00")
	if false1 != false || false2 != false {
		t.Errorf("Invalid langauge must be false")
	}
}
func TestIsValidChallenge(t *testing.T) {
	true1 := IsValidChallenge(ChallengeTypeJS)
	true2 := IsValidChallenge(ChallengeTypeCaptcha)
	if true1 != true || true2 != true {
		t.Errorf("Valid challenge must be true")
	}
	false1 := IsValidChallenge(ChallengeTypeJS + "xx")
	false2 := IsValidChallenge(ChallengeTypeCaptcha + "00")
	if false1 != false || false2 != false {
		t.Errorf("Invalid challenge must be false")
	}
}
