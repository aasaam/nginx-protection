package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	mathRandom "math/rand"
	"net/url"
	"time"
)

var noPadding rune = -1

func isValidChallenge(challenge string) bool {
	return supportedChallengesMap[challenge]
}

func hashHex(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func base64Hash(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func randomBase64() string {
	b := make([]byte, 512)
	_, err := rand.Read(b)
	if err != nil {
		panic(err.Error())
	}
	hash := sha512.New()
	hash.Write(b)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func minMaxDefault64(value int64, min int64, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

func minMaxDefault(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

func totpGenerate() string {
	b := make([]byte, 512)
	_, err := rand.Read(b)
	if err != nil {
		panic(err.Error())
	}
	hash := sha512.New()
	hash.Write(b)
	enc := base32.StdEncoding.WithPadding(noPadding)
	return enc.EncodeToString(hash.Sum(nil))[0:16]
}

func getURLPath(u *url.URL) string {
	if u.RawQuery != "" {
		return fmt.Sprintf("%s?%s", u.Path, u.RawQuery)
	}
	return u.Path
}

func getHumanEmoji() string {
	mathRandom.Seed(time.Now().Unix())
	i := mathRandom.Int() % len(humanEmojis)
	return humanEmojis[i]
}
