package main

import (
	"bytes"
	"crypto/rc4"
	"encoding/base64"
	"encoding/json"
	"image/color"
	"image/jpeg"
	"math/rand"
	"time"

	"github.com/afocus/captcha"
)

// Challenge is challenge will be encrypted
type Challenge struct {
	Type      string `json:"type"`
	Checksum  string `json:"checksum"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"start_end"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	Mobile    string `json:"mobile"`
	Country   string `json:"country"`
	Result    string `json:"result"`
	Content   string `json:"content"`
}

// NewChallenge is create new challenge
func NewChallenge(challengeType string, checksum string, wait int64, ttl int64) Challenge {

	challenge := Challenge{}
	challenge.Type = challengeType
	challenge.Checksum = checksum
	challenge.Result = ""
	challenge.Content = ""

	t := time.Now().UTC()
	StartTime := t.Add(time.Second * time.Duration(wait))

	EndTime := t.Add(time.Second * time.Duration(wait+ttl))

	challenge.StartTime = StartTime.Unix()
	challenge.EndTime = EndTime.Unix()

	return challenge
}

// HasResult return is challenge has result
func (c *Challenge) HasResult() bool {
	if c.Result != "" {
		return true
	}
	return false
}

// Validate the challenge
func (c *Challenge) Validate(challengeType string, checksum string) bool {
	if c.Type != challengeType {
		return false
	}

	if c.Checksum != checksum {
		return false
	}

	now := time.Now().UTC().Unix()

	if c.StartTime <= now && c.EndTime >= now {
		return true
	}

	return false
}

// JSON is convert challenge to json
func (c *Challenge) JSON() string {
	co := Challenge{}
	co.Type = c.Type
	co.Checksum = c.Checksum
	co.Result = c.Result
	co.StartTime = c.StartTime
	co.EndTime = c.EndTime

	bytes, _ := json.Marshal(co)

	return string(bytes)
}

// ChallengeFromJSON is convert json to challenge
func ChallengeFromJSON(jsonString string) (Challenge, error) {
	var challenge Challenge

	if err := json.Unmarshal([]byte(jsonString), &challenge); err != nil {
		return Challenge{}, err
	}

	return challenge, nil
}

// GenerateJSChallenge will generate eval for javascript challenge
func (c *Challenge) GenerateJSChallenge() {
	c.Result = RandomHex()

	src := []byte(c.Result)

	key := []byte(RandomHex())
	cipher, _ := rc4.NewCipher(key)
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)

	c.Content = base64.StdEncoding.EncodeToString(key) + ":" + base64.StdEncoding.EncodeToString(dst)
}

// GenerateCaptchaChallenge will generate captcha image for human detection
func (c *Challenge) GenerateCaptchaChallenge(farsi bool) {
	rand.Seed(time.Now().UnixNano())

	frontColor := color.RGBA{32, 32, 32, 0xff}
	backColor := color.RGBA{192, 192, 192, 0xff}

	var font []byte
	if farsi {
		max := len(FarsiCaptchaFonts)
		font = FarsiCaptchaFonts[rand.Intn(max-0)+0]
	} else {
		max := len(CaptchaFonts)
		font = CaptchaFonts[rand.Intn(max-0)+0]
	}

	cap := captcha.New()
	cap.SetSize(320, 128)
	cap.AddFontFromBytes(font)
	cap.SetFrontColor(frontColor)
	cap.SetBkgColor(backColor)
	cap.SetDisturbance(64)

	c.Result = RandomCaptchaLetters()
	img := cap.CreateCustom(c.Result)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, &jpeg.Options{
		Quality: 50,
	})

	c.Content = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}
