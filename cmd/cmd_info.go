package cmd

import (
	"log"
	"strings"

	"github.com/cute-angelia/avmetagetter/config"
	"github.com/cute-angelia/avmetagetter/pkg/media"
	"github.com/cute-angelia/avmetagetter/pkg/scraper"

	"github.com/cute-angelia/go-xutils/syntax/ifile"
	"github.com/cute-angelia/go-xutils/syntax/ijson"
	"github.com/cute-angelia/go-xutils/utils/conf"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func NewCmdInfo() *cli.Command {
	return &cli.Command{
		Name:        "info",
		Usage:       "info -no=avb-123",
		Description: "查询no信息",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "local",
				Usage: "环境变量",
			},
			&cli.StringFlag{
				Name:     "no",
				Value:    "",
				Usage:    "番号",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "nfo",
				Value: false,
				Usage: "数据以nfo返回",
			},
			&cli.StringFlag{
				Name:  "scraper",
				Value: "",
				Usage: "指定scraper：JavBus JavDb CaribBeanCom FC2 TokyoHot Heyzo Heydouga Siro memojav202508等",
			},
		},

		Action: func(c *cli.Context) error {

			envstr := c.String("env")
			no := c.String("no")
			nfo := c.Bool("nfo")
			captureNames := c.String("scraper")

			// 加载 config
			config.InitConfig(envstr)
			if ifile.IsExist("./config.toml") {
				if err := conf.MergeConfigWithPath("./config.toml"); err != nil {
					log.Println(err)
				}
			}

			crps := strings.Split(captureNames, ",")
			iscraper := scraper.NewScraper(no, viper.GetString("common.socks5"))
			if resps, err := iscraper.Search(crps); err != nil {
				return err
			} else {
				if nfo {
					for scraperName, resp := range resps {
						nfo := media.NewNfoJav()
						nfo.ParseMedia(resp)
						log.Println(ijson.Pretty(nfo))
						log.Println("站名", scraperName)
					}
				} else {
					log.Println(ijson.Pretty(resps))
				}
				return nil
			}
		},
	}
}
