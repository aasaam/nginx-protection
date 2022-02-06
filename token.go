package main

import (
	"encoding/json"
	"errors"

	aesGo "github.com/aasaam/aes-go"
)

type persistToken struct {
	Type        string `json:"t"`
	ClientIdent string `json:"c"`
	TTL         int64  `json:"ttl"`
}

func newPersistToken(typeOfToken string, clientIdent string, ttl int64) *persistToken {
	pt := persistToken{
		Type:        typeOfToken,
		ClientIdent: clientIdent,
		TTL:         ttl,
	}

	return &pt
}

func newPersistTokenFromString(tokenString string, secret string) (*persistToken, error) {
	aes := aesGo.NewAasaamAES(secret)
	jsonString := aes.DecryptTTL(tokenString)
	if jsonString == "" {
		return nil, errors.New("invalid token or expired")
	}

	var pt persistToken
	err := json.Unmarshal([]byte(jsonString), &pt)
	if err != nil {
		return nil, errors.New("invalid token structure")
	}

	return &pt, nil
}

func (pt *persistToken) generate(secret string) string {
	aes := aesGo.NewAasaamAES(secret)
	byteData, err := json.Marshal(&pt)
	if err != nil {
		return ""
	}
	return aes.EncryptTTL(string(byteData), pt.TTL)
}
