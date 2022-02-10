package main

import (
	"testing"

	aesGo "github.com/aasaam/aes-go"
)

func TestPersistToken(t *testing.T) {
	secret := aesGo.GenerateKey()
	tk1 := newPersistToken(challengeTypeJS, "c1", "", 1)
	str1 := tk1.generate(secret)
	_, err := newPersistTokenFromString(str1, secret)
	if err != nil {
		t.Error(err)
	}

	_, e1 := newPersistTokenFromString("a", secret)
	if e1 == nil {
		t.Errorf("must empty")
	}
}
