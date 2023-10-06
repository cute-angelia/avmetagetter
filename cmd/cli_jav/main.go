package main

import (
	_ "embed"
	"errors"
	"fmt"
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
				Usage:       "番号",
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
				site := sites.NewSite("javbus", no)
				log.Println(site)
				if resp, err := site.Fetch(); err != nil {
					return err
				} else {
					log.Println(fmt.Sprintf("%#v", resp))
					return nil
				}
			} else {
				return errors.New("需要一个番号")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
