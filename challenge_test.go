package main

import (
	"testing"
	"time"
)

func TestNewChallenge1(t *testing.T) {
	challengeChecksum := "random"
	challengeType := ChallengeTypeJS
	ch1 := NewChallenge(challengeType, challengeChecksum, 1, 3)
	if ch1.Validate(challengeType, challengeChecksum) == true {
		t.Errorf("Wait time not passed so it's must be false")
	}
	time.Sleep(time.Second * 2)
	if ch1.Validate(ChallengeTypeCaptcha, challengeChecksum) == true {
		t.Errorf("Type is not the same")
	}
	if ch1.Validate(challengeType, "other") == true {
		t.Errorf("Checksum is not the same")
	}
	if ch1.Validate(challengeType, challengeChecksum) == false {
		t.Errorf("Wait time passed so it's must be true")
	}
}

func TestChallengeJSON(t *testing.T) {
	challengeChecksum := "random"
	challengeType := ChallengeTypeJS
	ch1 := NewChallenge(challengeType, challengeChecksum, 1, 3)
	json := ch1.JSON()
	ch2, err := ChallengeFromJSON(json)
	if err != nil || ch1.Checksum != ch2.Checksum {
		t.Errorf("Challenge json error or not equal after serialize/deserialize")
	}

	_, err1 := ChallengeFromJSON("invali")
	if err1 == nil {
		t.Errorf("invalid json must throw an error")
	}
}
func TestGenerateJSChallenge(t *testing.T) {
	challengeChecksum := "random"
	challengeType := ChallengeTypeJS
	ch1 := NewChallenge(challengeType, challengeChecksum, 1, 3)
	if ch1.HasResult() == true {
		t.Errorf("No result")
	}
	ch1.GenerateJSChallenge()
	if ch1.HasResult() == false {
		t.Errorf("No result")
	}
	if ch1.Result == "" || ch1.Content == "" {
		t.Errorf("Result and content must be set")
	}
}
func TestGenerateCaptchaChallenge(t *testing.T) {
	challengeChecksum := "random"
	challengeType := ChallengeTypeJS
	ch1 := NewChallenge(challengeType, challengeChecksum, 1, 3)
	ch1.GenerateCaptchaChallenge(true)
	ch2 := NewChallenge(challengeType, challengeChecksum, 1, 3)
	ch2.GenerateCaptchaChallenge(false)
	if ch1.Result == "" || ch1.Content == "" {
		t.Errorf("Result and content must be set")
	}
	if ch2.Result == "" || ch2.Content == "" {
		t.Errorf("Result and content must be set")
	}
}
