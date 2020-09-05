package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "nginx-protection"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:  "otp-secret",
			Usage: "Generate OTP Secret for usage",
			Action: func(c *cli.Context) error {
				secret := GenerateOTPSecret()
				fmt.Println("Time configure is depend on your configuration")
				fmt.Println("Algorithmic is SHA512")
				fmt.Println("And your secret could be: ")
				fmt.Println(string(secret))
				return nil
			},
		},
		{
			Name:  "webserver",
			Usage: "HTTP Server for REST API",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "listen", Aliases: []string{"l"}, EnvVars: []string{"AASAAM_WEBSERVER_LISTEN"}, Value: "127.0.0.1:19000", DefaultText: "127.0.0.1:19000", Usage: "HTTP address:port to listen"},
				&cli.StringFlag{Name: "base-url", Aliases: []string{"b"}, EnvVars: []string{"AASAAM_WEBSERVER_BASEURL"}, Value: "/.well-known/protection", DefaultText: "/.well-known/protection", Usage: "Base URL to serve server"},
				&cli.StringFlag{Name: "enable-log", Aliases: []string{"log"}, EnvVars: []string{"AASAAM_ENABLE_LOG"}, Value: "0", DefaultText: "0", Usage: "HTTP address:port to listen"},
				&cli.StringFlag{Name: "salt", Aliases: []string{"s"}, EnvVars: []string{"AASAAM_SALT"}, Usage: "Salt for encryption"},
				&cli.StringFlag{Name: "private-key", Aliases: []string{"key"}, EnvVars: []string{"AASAAM_PRIVATE_KEY"}, Usage: "Path of RSA private key"},
			},
			Action: func(c *cli.Context) error {
				config := Config{}
				config.BaseURL = c.String("base-url")
				config.Salt = c.String("salt")
				config.Testing = false
				config.Log = c.Bool("enable-log")

				pemBytes, err := ioutil.ReadFile(c.String("private-key"))
				if err != nil {
					return cli.Exit("RSA Key file not reachable", 128)
				}
				rsa, err := LoadPrivateKey(pemBytes)
				if err != nil {
					return cli.Exit("Invalid RSA Key file format", 128)
				}

				app := GetHTTPServer(&config, rsa)

				app.Listen(c.String("listen"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
