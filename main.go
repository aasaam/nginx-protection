package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	aes "github.com/aasaam/aes-go"
	"github.com/mdp/qrterminal"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/urfave/cli/v2"
)

func generateSecret(c *cli.Context) error {
	fmt.Println(aes.GenerateKey())
	return nil
}

func generateTOTPSecret(c *cli.Context) error {
	secret := totpGenerate()
	otpURL := fmt.Sprintf("otpauth://totp/protection@%s", c.String("host"))

	u, _ := url.Parse(otpURL)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("secret", secret)
	u.RawQuery = q.Encode()
	fmt.Println("Nginx Variable:")

	nginxVariable1 := fmt.Sprintf("set $protection_config_totp_secret '%s';", secret)
	fmt.Println("")
	fmt.Println(nginxVariable1)
	fmt.Println("")
	fmt.Println("URI:")
	fmt.Println(u)
	fmt.Println("")
	fmt.Println("QR Code:")

	config := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	}

	qrterminal.GenerateWithConfig(u.String(), config)

	return nil
}

func totpCheck(c *cli.Context) error {
	passCode, err := totp.GenerateCodeCustom(c.String("secret"), time.Now(), totp.ValidateOpts{
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return err
	}
	valid := passCode == c.String("pass")
	fmt.Println(valid)
	return nil
}

func runServer(c *cli.Context) error {

	challengeStorage := newChallengeStorage()
	aclStorage := newACLStorage()

	config := newConfig(
		c.String("log-level"),
		c.Bool("aasaam-web-server"),
		c.String("default-language"),
		c.String("supported-languages"),
		c.String("token-secret"),
		c.String("client-secret"),
		c.String("base-url"),
		c.String("static-url"),
		c.String("rest-captcha-server-url"),
		c.String("locale-path"),
	)

	go func() {
		for {
			challengeStorageCount := challengeStorage.gc()
			aclStorageCount := aclStorage.gc()
			defer config.getLogger().
				Debug().
				Str(logType, logTypeApp).
				Int("challenge_storage_count", challengeStorageCount).
				Int("acl_storage_count", aclStorageCount).
				Send()

			time.Sleep(time.Second * 10)
		}
	}()

	loadLocales(config)
	app := newHTTPServer(config, challengeStorage, aclStorage)
	return app.Listen(c.String("listen"))
}

func main() {
	app := cli.NewApp()
	app.Usage = "nginx protection"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "generate-secret",
			Usage:  "Generate secret for security usage",
			Action: generateSecret,
		},
		{
			Name:   "generate-totp",
			Usage:  "Generate QR code for otp",
			Action: generateTOTPSecret,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Usage:    "Hostname",
					Required: true,
				},
			},
		},
		{
			Name:   "check-totp",
			Usage:  "Check otp secret",
			Action: totpCheck,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "secret",
					Usage:    "Secret",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "pass",
					Usage:    "Password",
					Required: true,
				},
			},
		},
		{
			Name:  "run",
			Usage: "Run protection server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "token-secret",
					Usage:    "Token Secret",
					Required: true,
					EnvVars:  []string{"ASM_TOKEN_SECRET"},
				},
				&cli.StringFlag{
					Name:     "client-secret",
					Usage:    "Client Secret",
					Required: true,
					EnvVars:  []string{"ASM_CLIENT_SECRET"},
				},
				&cli.StringFlag{
					Name:    "listen",
					Usage:   "Listen HTTP address for application",
					Value:   "127.0.0.1:9121",
					EnvVars: []string{"ASM_HTTP_LISTEN"},
				},
				&cli.StringFlag{
					Name:    "rest-captcha-server-url",
					Usage:   "REST captcha server URL",
					Value:   "http://127.0.0.1:4000",
					EnvVars: []string{"ASM_REST_CAPTCHA_URL"},
				},
				&cli.StringFlag{
					Name:    "static-url",
					Usage:   "Static data base URL for CDN usage like https://cdn.jsdelivr.net/gh/aasaam/nginx-protection@static",
					Value:   "",
					EnvVars: []string{"ASM_REST_CAPTCHA_URL"},
				},
				&cli.StringFlag{
					Name:    "base-url",
					Usage:   "Base URL of application",
					Value:   "/.well-known/protection",
					EnvVars: []string{"ASM_HTTP_BASE_URL"},
				},
				&cli.StringFlag{
					Name:    "default-language",
					Usage:   "Default language",
					EnvVars: []string{"ASM_DEFAULT_LANG"},
				},
				&cli.StringFlag{
					Name:    "supported-languages",
					Usage:   "Comma separeted list of supported languages",
					Value:   "en,fa",
					EnvVars: []string{"ASM_LANGUAGES"},
				},
				&cli.StringFlag{
					Name:    "locale-path",
					Usage:   "Path of merged locale data",
					Value:   "/etc/nginx-protection/locale",
					EnvVars: []string{"ASM_LOCALE_DIR"},
				},
				&cli.StringFlag{
					Name:    "log-level",
					Usage:   "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:   "trace",
					EnvVars: []string{"ASM_LOCALE_DIR"},
				},
				&cli.BoolFlag{
					Name:    "aasaam-web-server",
					Usage:   "Run via aasaam web server",
					EnvVars: []string{"ASM_AASAAM_WEB_SERVER"},
				},
			},
			Action: runServer,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
