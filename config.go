package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type config struct {
	aasaamWebServer    bool
	tokenSecret        string
	clientSecret       string
	baseURL            string
	cdnStatic          bool
	staticURL          string
	restCaptchaURL     string
	defaultLanguage    string
	supportedLanguages []string
	localePath         string

	logger *zerolog.Logger
}

func newConfig(
	logLevel string,
	aasaamWebServer bool,
	defaultLanguage string,
	supportedLanguages string,
	tokenSecret string,
	clientSecret string,
	baseURL string,
	staticURL string,
	restCaptchaURL string,
	localePath string,
) *config {
	c := config{
		aasaamWebServer:    aasaamWebServer,
		tokenSecret:        tokenSecret,
		clientSecret:       clientSecret,
		baseURL:            strings.TrimRight(baseURL, "/"),
		staticURL:          "",
		cdnStatic:          false,
		restCaptchaURL:     strings.TrimRight(restCaptchaURL, "/"),
		defaultLanguage:    "en",
		supportedLanguages: []string{"en"},
		localePath:         localePath,
	}

	// logger config
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	logConfigLevel, errLogLevel := zerolog.ParseLevel(logLevel)
	if errLogLevel == nil {
		zerolog.SetGlobalLevel(logConfigLevel)
	}
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	c.logger = &logger

	_, err := url.ParseRequestURI(staticURL)
	if err == nil {
		c.staticURL = staticURL
		c.cdnStatic = true
	} else {
		c.staticURL = c.baseURL + "/challenge/static"
	}

	if supportedLanguages != "" {
		langs := strings.Split(strings.TrimSpace(supportedLanguages), ",")
		c.supportedLanguages = []string{}
		for _, l := range langs {
			if isSupportedLanguage(l) {
				c.supportedLanguages = append(c.supportedLanguages, l)
			}
		}
	}

	if isSupportedLanguage(defaultLanguage) {
		c.defaultLanguage = defaultLanguage
	}

	if len(c.supportedLanguages) == 0 {
		c.supportedLanguages = []string{c.defaultLanguage}
	}

	return &c
}

func (c *config) getLogger() *zerolog.Logger {
	return c.logger
}
