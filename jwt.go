package main

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// LoadPrivateKey from pem string
func LoadPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("Block Private Key is empty")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if key, ok := parseResult.(*rsa.PrivateKey); ok {
		return key, nil
	}
	return nil, errors.New("Key not parsable")
}

// ValidateAuthToken is generate token for auth
func ValidateAuthToken(encrypted string, key string, rsa *rsa.PrivateKey) (AuthToken, error) {
	authTokenString, err := VerifyToken(encrypted, key, rsa)
	if err != nil {
		return AuthToken{}, err
	}
	var authToken AuthToken
	if err := json.Unmarshal([]byte(authTokenString), &authToken); err != nil {
		return AuthToken{}, err
	}

	return authToken, nil
}

// GenerateAuthToken is generate token for auth
func GenerateAuthToken(authToken AuthToken, ttl int64, key string, rsa *rsa.PrivateKey) (string, error) {
	bytes, _ := json.Marshal(authToken)

	return GenerateToken(string(bytes), 0, ttl, key, rsa)
}

// GenerateToken creat token for special user
func GenerateToken(subject string, from int64, ttl int64, key string, rsa *rsa.PrivateKey) (string, error) {
	hash := md5.New()
	hash.Write([]byte(key))

	sharedKey := hash.Sum(nil)
	enc, err := jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       sharedKey,
		},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"))
	if err != nil {
		return "", err
	}

	Expiry := time.Now().UTC()
	Expiry = Expiry.Add(time.Second * time.Duration(ttl))

	NotBefore := time.Now().UTC()
	NotBefore = NotBefore.Add(time.Second * time.Duration(from))

	claim := jwt.Claims{
		Subject:   base58.Encode([]byte(subject)),
		NotBefore: jwt.NewNumericDate(NotBefore),
		Expiry:    jwt.NewNumericDate(Expiry),
	}
	rsaSigner, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: rsa}, nil)
	if err != nil {
		return "", err
	}
	token, err := jwt.SignedAndEncrypted(rsaSigner, enc).Claims(claim).CompactSerialize()
	if err != nil {
		return "", err
	}
	return base58.Encode([]byte(token)), nil
}

// VerifyToken will verify the generated token from client
func VerifyToken(token string, key string, rsa *rsa.PrivateKey) (string, error) {
	hash := md5.New()
	hash.Write([]byte(key))

	sharedKey := hash.Sum(nil)

	jwtToken := string(base58.Decode(token))

	tok, err := jwt.ParseSignedAndEncrypted(jwtToken)
	if err != nil {
		return "", err
	}

	nested, err := tok.Decrypt(sharedKey)
	if err != nil {
		return "", err
	}

	out := jwt.Claims{}

	if err := nested.Claims(&rsa.PublicKey, &out); err != nil {
		return "", err
	}

	if out.NotBefore.Time().Unix() <= time.Now().Unix() && out.Expiry.Time().Unix() >= time.Now().Unix() {
		sub := string(base58.Decode(out.Subject))
		return sub, nil
	}

	return "", errors.New("Invlid token NotBefore or Expiry")
}
