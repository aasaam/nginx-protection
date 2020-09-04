package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"
)

// IsValidLanguage check for supporeted challenge
func IsValidLanguage(lang string) bool {
	supportedLanguage := reflect.ValueOf(SupportedLanguage)

	for i := 0; i < supportedLanguage.Len(); i++ {
		if supportedLanguage.Index(i).Interface() == lang {
			return true
		}
	}

	return false
}

// GenerateOTPSecret return random generated OTP secret
func GenerateOTPSecret() string {
	data := []byte(RandomHex())

	str := base32.StdEncoding.EncodeToString(data)

	return str[0:16]
}

// IsValidChallenge check for supporeted challenge
func IsValidChallenge(challenge string) bool {
	supportedChallenge := reflect.ValueOf(SupportedChallenges)

	for i := 0; i < supportedChallenge.Len(); i++ {
		if supportedChallenge.Index(i).Interface() == challenge {
			return true
		}
	}

	return false
}

// RandomCaptchaLetters return random letters for captcha
func RandomCaptchaLetters() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 6)
	for i := range b {
		b[i] = CaptchaLetters[rand.Intn(len(CaptchaLetters))]
	}
	return string(b)
}

// RandomHex generate random hex
func RandomHex() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	hash := sha512.New()
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// MD5 return md5 of input string
func MD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}
