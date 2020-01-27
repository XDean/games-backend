package main

import (
	"games-backend/app"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	configPath string
	debug      bool
	secretKey  string
)

func main() {

	cliApp := cli.NewApp()

	cliApp.Name = "Wechat BG"
	cliApp.Usage = "Run XDean Wechat BG Server"

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "setting,s",
			Usage:       "Setting file path",
			Destination: &configPath,
		},
		cli.StringFlag{
			Name:        "key,k",
			Usage:       "Secret key",
			Destination: &secretKey,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "Debug mode",
			Destination: &debug,
		},
	}

	cliApp.Action = func(c *cli.Context) error {
		if debug {
			app.Debug()
		}
		app.App.Config.SecretKey = secretKey
		if configPath != "" {
			app.App.RegisterConfigPath(configPath)
		}
		app.App.Run()
		return nil
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
