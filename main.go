package main

import (
	"avmetagetter/cmd"
	_ "embed"
	"log"
	"os"

	"github.com/cute-angelia/go-xutils/components/loggers/loggerV3"
	"github.com/urfave/cli/v2"
)

func main() {
	// 日志
	loggerV3.New(loggerV3.WithIsOnline(false))

	app := cli.NewApp()
	app.Name = "AV工具"
	app.Description = "查询/nfo"

	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true

	app.Commands = []*cli.Command{
		cmd.NewCmdInfo(),
		cmd.NewCmdMove(),
		cmd.NewCmdSite(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
