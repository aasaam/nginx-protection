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
	supportedLangauges []string
	localePath         string

	logger *zerolog.Logger
}

func newConfig(
	logLevel string,
	aasaamWebServer bool,
	defaultLanguage string,
	supportedLangauges string,
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
		supportedLangauges: []string{"en"},
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

	if supportedLangauges != "" {
		langs := strings.Split(strings.TrimSpace(supportedLangauges), ",")
		c.supportedLangauges = []string{}
		for _, l := range langs {
			if isSupportedLangauge(l) {
				c.supportedLangauges = append(c.supportedLangauges, l)
			}
		}
	}

	if isSupportedLangauge(defaultLanguage) {
		c.defaultLanguage = defaultLanguage
	}

	if len(c.supportedLangauges) == 0 {
		c.supportedLangauges = []string{c.defaultLanguage}
	}

	return &c
}

func (c *config) getLogger() *zerolog.Logger {
	return c.logger
}
