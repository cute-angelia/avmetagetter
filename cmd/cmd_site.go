package cmd

import (
	"github.com/cute-angelia/avmetagetter/pkg/scraper"
	"log"

	"github.com/urfave/cli/v2"
)

func NewCmdSite() *cli.Command {
	return &cli.Command{
		Name:        "site",
		Usage:       "查询支持站点信息",
		Description: "查询支持站点信息",
		Action: func(c *cli.Context) error {
			iscraper := scraper.NewScraper("", "")
			captures := iscraper.GetCaptures([]string{})
			for _, capture := range captures {
				if capture.Enable {
					log.Println(capture.Name + " " + capture.Desc)
				}
			}
			return nil
		},
	}
}
