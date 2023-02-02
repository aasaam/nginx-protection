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

func generateSecret(_ *cli.Context) error {
	fmt.Println(aes.GenerateKey())
	return nil
}

func generateTOTPSecret(c *cli.Context) error {
	fmt.Println(totpGenerate())
	return nil
}

func generateTOTP(c *cli.Context) error {
	otpURL := fmt.Sprintf("otpauth://totp/protection@%s", c.String("host"))
	secret := c.String("secret")
	u, _ := url.Parse(otpURL)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("secret", secret)
	u.RawQuery = q.Encode()

	t1 := time.Now()
	passCode1, err := totp.GenerateCodeCustom(secret, t1, totp.ValidateOpts{
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		panic(err)
	}

	t2 := time.Now().Add(time.Second * time.Duration(30))
	passCode2, err := totp.GenerateCodeCustom(secret, t2, totp.ValidateOpts{
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Current password code:")
	fmt.Println(t1.Format("15:04:05 MST"))
	fmt.Println(passCode1)
	fmt.Println("Next password code:")
	fmt.Println(t2.Format("15:04:05 MST"))
	fmt.Println(passCode2)
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
	fmt.Println(passCode)
	fmt.Println(valid)
	return nil
}

func ldapCheck(c *cli.Context) error {
	success, ldapCheckError, ldapConfigError := ldapLogin(
		c.String("url"),
		c.String("read-only-username"),
		c.String("read-only-password"),
		c.String("base-dn"),
		c.String("filter"),
		c.String("attributes-json"),
		c.String("username"),
		c.String("password"),
	)

	if ldapCheckError != nil {
		fmt.Println("LDAP check error")
		return ldapCheckError
	}

	if ldapConfigError != nil {
		fmt.Println("LDAP config error")
		return ldapConfigError
	}

	if success {
		fmt.Println("Login successfull")
	} else {
		fmt.Println("Login failed")
	}

	return nil
}

func runServer(c *cli.Context) error {

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

	loadLocales(config)
	app := newHTTPServer(config)
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
			Name:   "generate-totp-secret",
			Usage:  "Generate TOTP secret",
			Action: generateTOTPSecret,
		},
		{
			Name:   "generate-totp",
			Usage:  "Generate QR code for otp",
			Action: generateTOTP,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Usage:    "Hostname",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "secret",
					Usage:    "Secret of TOTP",
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
			Name:   "check-ldap",
			Usage:  "Check ldap login",
			Action: ldapCheck,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Usage:    "LDAP server URL",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "read-only-username",
					Usage:    "Read only username",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "read-only-password",
					Usage:    "Read only password",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "base-dn",
					Usage:    "Base distinguished name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "filter",
					Usage:    "Filter",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "attributes-json",
					Usage: "Attributes list in JSON format",
					Value: `["dn"]`,
				},
				&cli.StringFlag{
					Name:     "username",
					Usage:    "Username",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
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
					EnvVars:  []string{"ASM_NGINX_PROTECTION_TOKEN_SECRET"},
				},
				&cli.StringFlag{
					Name:     "client-secret",
					Usage:    "Client Secret",
					Required: true,
					EnvVars:  []string{"ASM_NGINX_PROTECTION_CLIENT_SECRET"},
				},
				&cli.StringFlag{
					Name:    "listen",
					Usage:   "Listen HTTP address for application",
					Value:   "127.0.0.1:9121",
					EnvVars: []string{"ASM_NGINX_PROTECTION_LISTEN"},
				},
				&cli.StringFlag{
					Name:    "rest-captcha-server-url",
					Usage:   "REST Captcha server URL",
					Value:   "http://127.0.0.1:4000",
					EnvVars: []string{"ASM_NGINX_PROTECTION_REST_CAPTCHA_URL"},
				},
				&cli.StringFlag{
					Name:    "static-url",
					Usage:   "Static data base URL for CDN usage like https://cdn.jsdelivr.net/gh/aasaam/nginx-protection@static",
					Value:   "",
					EnvVars: []string{"ASM_NGINX_PROTECTION_STATIC_URL"},
				},
				&cli.StringFlag{
					Name:    "base-url",
					Usage:   "Base URL of application",
					Value:   "/.well-known/protection",
					EnvVars: []string{"ASM_NGINX_PROTECTION_HTTP_BASE_URL"},
				},
				&cli.StringFlag{
					Name:    "default-language",
					Usage:   "Default language",
					EnvVars: []string{"ASM_NGINX_PROTECTION_DEFAULT_LANG"},
				},
				&cli.StringFlag{
					Name:    "supported-languages",
					Usage:   "Comma separeted list of supported languages",
					Value:   "en,fa",
					EnvVars: []string{"ASM_NGINX_PROTECTION_SUPPORTED_LANGUAGES"},
				},
				&cli.StringFlag{
					Name:    "locale-path",
					Usage:   "Path of merged locale data",
					Value:   "/etc/nginx-protection/locale",
					EnvVars: []string{"ASM_NGINX_PROTECTION_LOCALE_PATH"},
				},
				&cli.StringFlag{
					Name:    "log-level",
					Usage:   "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:   "warn",
					EnvVars: []string{"ASM_NGINX_PROTECTION_LOG_LEVEL"},
				},
				&cli.BoolFlag{
					Name:    "aasaam-web-server",
					Usage:   "Run via aasaam web server",
					EnvVars: []string{"ASM_NGINX_PROTECTION_VIA_AASAAM_WEB_SERVER"},
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
