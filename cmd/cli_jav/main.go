package main

import (
	_ "embed"
	"errors"
	"github.com/cute-angelia/go-utils/components/loggers/loggerV3"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"github.com/cute-angelia/go-utils/utils/conf"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"metagetter/pkg/media"
	"metagetter/pkg/scraper"
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
			&cli.BoolFlag{
				Name:  "nfo",
				Usage: "数据以nfo返回",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if len(no) == 0 {
				if cCtx.NArg() > 0 {
					no = cCtx.Args().Get(0)
				}
			}
			if len(no) > 0 {
				iscraper := scraper.NewScraper(no, viper.GetString("common.socks5"))
				if resp, err := iscraper.Search(); err != nil {
					return err
				} else {
					// nfo
					if cCtx.Bool("nfo") {
						nfo := media.NewNfoJav()
						nfo.ParseMedia(resp)
						log.Println(ijson.Pretty(nfo))
					} else {
						log.Println(ijson.Pretty(resp))
					}
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
