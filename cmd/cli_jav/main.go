package main

import (
	_ "embed"
	"github.com/cute-angelia/go-utils/components/loggers/loggerV3"
	"github.com/cute-angelia/go-utils/utils/conf"
	"github.com/urfave/cli/v2"
	"log"
	"metagetter/pkg/sites"
	"os"
)

func main() {
	// config.toml
	conf.MustLoadConfigFile("config.toml")

	// logger
	loggerV3.New(loggerV3.WithIsOnline(false))

	var no string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "no",
				Value:       "",
				Usage:       "no for movie",
				Destination: &no,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if len(no) == 0 {
				if cCtx.NArg() > 0 {
					no = cCtx.Args().Get(0)
				}
			}
			if len(no) > 0 {
				sites.NewJavBus(no, "", "", "")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
