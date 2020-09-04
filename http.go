package main

import (
	"crypto/rsa"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
	"github.com/gofiber/template/html"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	rice "github.com/GeertJohan/go.rice"
)

// GetHTTPServer return fiber server
func GetHTTPServer(config *Config, rsa *rsa.PrivateKey) *fiber.App {
	engine := html.NewFileSystem(rice.MustFindBox("./assets/templates").HTTPBox(), ".html")
	engine.Delims("[[", "]]")

	r := prometheus.NewRegistry()
	r.MustRegister(PrometheusRequestTotal)
	r.MustRegister(PrometheusRejectTotal)
	r.MustRegister(PrometheusAccpetTotal)
	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	srv := fiber.New(&fiber.Settings{
		Views:                 engine,
		Prefork:               false,
		DisableStartupMessage: true,
		CaseSensitive:         true,
		StrictRouting:         true,
		UnescapePath:          true,
		ServerHeader:          "protection",
	})

	srv.Get("/metrics", func(c *fiber.Ctx) {
		p(c.Fasthttp)
	})

	srv.Use(helmet.New(helmet.Config{
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "default-src 'self' data:; script-src 'self' data: 'unsafe-inline' 'unsafe-eval'; style-src 'self' data: 'unsafe-inline'",
	}))

	api := srv.Group(config.BaseURL)

	api.Use("/challenge/assets", middleware.FileSystem(rice.MustFindBox("./assets/public").HTTPBox()))

	RouterAuth(config, api, rsa)
	RouterChallenge(config, api, rsa)
	RouterChallengeJS(config, api, rsa)
	RouterChallengeCaptcha(config, api, rsa)
	RouterChallengeOTP(config, api, rsa)

	return srv
}
