package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/cute-angelia/avmetagetter/config"
	"github.com/cute-angelia/avmetagetter/pkg/media"
	"github.com/cute-angelia/avmetagetter/pkg/scraper"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/cute-angelia/go-utils/components/idownload"
	"github.com/cute-angelia/go-utils/components/loggers/loggerV3"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/cute-angelia/go-utils/utils/conf"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// logger
	loggerV3.New(loggerV3.WithIsOnline(false))

	var dir string
	var dest string
	var envstr string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Value:       "local",
				Usage:       "环境变量",
				Destination: &envstr,
			},
			&cli.StringFlag{
				Name:        "dir",
				Value:       "./",
				Usage:       "扫描文件夹",
				Destination: &dir,
			},
			&cli.StringFlag{
				Name:        "dest",
				Value:       "./jav",
				Usage:       "目标文件夹",
				Destination: &dest,
			},
		},
		Action: func(cCtx *cli.Context) error {
			config.InitConfig(envstr)
			conf.MergeConfigWithPath("./")

			cdir := viper.GetString("avnas.dir")
			if len(cdir) > 0 {
				dir = cdir
			}

			cdest := viper.GetString("avnas.dest")
			if len(cdest) > 0 {
				dest = cdest
			}

			return fire(dir, dest)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fire(dir string, dest string) error {
	exts := []string{
		".mp4",
		".m4p",
		".mkv",
		".avi",
		".mpeg",
	}
	if _, files, err := ifile.GetDepthOnePathsAndFilesIncludeExt(dir, exts...); err != nil {
		return err
	} else {
		for _, avfile := range files {
			no := utils.CleanNo(ifile.NameNoExt(avfile))
			log.Println("处理：", avfile, "-->", no)
			// 抓取信息
			iscraper := scraper.NewScraper(no, viper.GetString("common.socks5"), []string{})
			if resp, err := iscraper.Search(); err != nil {
				loggerV3.GetLogger().Err(err).Str("抓取失败", no).Send()
				continue
			} else {

				destdir := ""
				title := ""
				if !strings.Contains(resp.Title, no) {
					title = no + " "
				} else {
					title = resp.Title
				}

				if len(resp.Title) == 0 {
					return errors.New("抓取失败 " + no)
				}

				// nfo
				nfo := media.NewNfoJav()
				nfo.ParseMedia(resp)
				nfo.SetPoster("poster.jpg")
				nfo.SetFanArt("fanart.jpg")

				actorName := ""
				if len(nfo.Actor) > 0 {
					actorName = nfo.Actor[0].Name
				} else {
					actorName = "未知"
				}

				// 生成目标文件夹
				// 规则 [2019] STARS-065 ナマ派 初中出し解禁 本庄鈴
				destdir = fmt.Sprintf("%s/[%s] %s", actorName, nfo.Year, title)

				// 生成 inf
				nfoPath := filepath.Join(dest, destdir, fmt.Sprintf("%s.nfo", no))
				nfoFile, _ := ifile.OpenLocalFile(nfoPath)
				os.Truncate(nfoPath, 0)
				nfobyte, _ := nfo.Marshal()
				nfoFile.Write(nfobyte)
				nfoFile.Close()

				// 生成图片

				idown := idownload.New(
					idownload.WithProxySocks5(viper.GetString("common.socks5")),
					idownload.WithTimeout(time.Minute),
					idownload.WithReferer(viper.GetString("javbus.site")),
					idownload.WithCookie(viper.GetString("javbus.cookies")),
					idownload.WithUserAgent(viper.GetString("javbus.useragent")),
				)

				fanart := filepath.Join(dest, destdir, "fanart.jpg")
				thumb := filepath.Join(dest, destdir, "poster.jpg")
				if _, err := idown.Download(resp.Cover, fanart); err != nil {
					log.Println(err)
				}

				utils.MakeThumbCover(fanart, thumb)

				// 移动资源到目标文件夹
				dst := filepath.Join(dest, destdir, ifile.Name(avfile))
				os.Rename(avfile, dst)

				loggerV3.GetLogger().Info().Str("目标路径", dst).Send()
			}
		}
		return nil
	}
}
