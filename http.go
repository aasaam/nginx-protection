package main

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

func newHTTPServer(config *config, challengeStorage *challengeStorage, aclStorage *aclStorage) *fiber.App {
	engine := html.NewFileSystem(http.FS(templates), ".html")
	engine.Delims("[[", "]]")
	engine.AddFunc(
		"unescapeJS", func(s string) template.JS {
			return template.JS(s)
		},
	)
	engine.AddFunc(
		"unescapeHTML", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Prefork:               false,
		Views:                 engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			defer prometheusRequestError.WithLabelValues(strconv.Itoa(code)).Inc()

			defer config.getLogger().Warn().
				Str("type", "http_error").
				Str("error", err.Error()).
				Int("status_code", code).
				Str("url", c.Request().URI().String()).
				Str("method", c.Method()).
				Str("rid", c.Get(httpRequestHeaderRequestID, "")).
				Send()

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).SendString("Internal Server Error")
		},
	})

	app.Use(recover.New())
	app.Use(helmet.New())

	handler := promhttp.HandlerFor(getPrometheusRegistry(), promhttp.HandlerOpts{})
	app.Get("/metrics", adaptor.HTTPHandler(handler))

	// middle ware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Robots-Tag", "noindex,nofollow")

		defer config.getLogger().
			Debug().
			Str("type", "http_request").
			Str("url", c.Request().URI().String()).
			Str("method", c.Method()).
			Str("rid", c.Get(httpRequestHeaderRequestID, "")).
			Send()

		// count expect static or metrics
		if !strings.Contains(c.Path(), "static") && !strings.Contains(c.Path(), "metrics") {
			defer prometheusRequestTotal.Inc()
		}

		return c.Next()
	})

	app.All(config.baseURL+"/auth", func(c *fiber.Ctx) error {
		configError := checkConfiguration(c)
		if configError != nil {
			return fiber.NewError(misconfigureStatus, "Configuration failed: "+configError.Error())
		}

		success := checkAuth(c, config, aclStorage)
		if success {
			return c.JSON("Authorized")
		}

		unAuthorizedStatus := getConfigUnauthorizedStatus(c)
		c.Status(unAuthorizedStatus)
		return c.JSON("Unauthorized")
	})

	app.Get(config.baseURL+"/challenge", func(c *fiber.Ctx) error {
		configError := checkConfiguration(c)
		if configError != nil {
			return fiber.NewError(misconfigureStatus, "Configuration failed: "+configError.Error())
		}

		success := checkAuth(c, config, aclStorage)
		if success {
			return c.Redirect(getProtectedPath(c))
		}

		return httpChallenge(c, config)
	})

	app.Patch(config.baseURL+"/challenge", func(c *fiber.Ctx) error {
		configError := checkConfiguration(c)
		if configError != nil {
			return fiber.NewError(misconfigureStatus, "Configuration failed: "+configError.Error())
		}

		return httpChallengePatch(c, config)
	})

	app.Post(config.baseURL+"/challenge", func(c *fiber.Ctx) error {
		configError := checkConfiguration(c)
		if configError != nil {
			return fiber.NewError(misconfigureStatus, "Configuration failed: "+configError.Error())
		}

		return httpChallengePost(c, config, challengeStorage)
	})

	// static serve
	if !config.cdnStatic {
		config.getLogger().Warn().Msg("disable cdn mode")
		appConfig := filesystem.Config{
			Next: func(c *fiber.Ctx) bool {
				c.Set("Cache-Control", "public, max-age=14400")
				return false
			},
			Root: http.FS(static),
		}

		app.Use(config.baseURL+"/challenge", filesystem.New(appConfig))
	} else {
		config.getLogger().Warn().Msg("enable cdn mode")
	}

	// 404
	app.Use(func(c *fiber.Ctx) error {
		defer prometheusRequestError.WithLabelValues("404").Inc()
		defer config.getLogger().
			Warn().
			Str("type", "http_404").
			Str("url", c.Request().URI().String()).
			Str("method", c.Method()).
			Str("rid", c.Locals(localVarRequestID).(string)).
			Send()

		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	return app
}
