package main

import (
	"testing"
)

func TestUtils1(t *testing.T) {
	if isValidChallenge("i'm sure not exist") {
		t.Errorf("invalid")
	}

	if !isValidChallenge(challengeTypeJS) {
		t.Errorf("invalid")
	}

	base64Hash("")
	hashHex("")
	totpGenerate()
}

func TestMinMaxDefault64(t *testing.T) {
	v1 := minMaxDefault64(1, 10, 20)
	if v1 != 10 {
		t.Errorf("invalid min")
	}
	v2 := minMaxDefault64(30, 10, 20)
	if v2 != 20 {
		t.Errorf("invalid max")
	}
	v3 := minMaxDefault64(30, 10, 50)
	if v3 != 30 {
		t.Errorf("invalid value")
	}
}
func TestMinMaxDefault(t *testing.T) {
	v1 := minMaxDefault(1, 10, 20)
	if v1 != 10 {
		t.Errorf("invalid min")
	}
	v2 := minMaxDefault(30, 10, 20)
	if v2 != 20 {
		t.Errorf("invalid max")
	}
	v3 := minMaxDefault(30, 10, 50)
	if v3 != 30 {
		t.Errorf("invalid value")
	}
}
