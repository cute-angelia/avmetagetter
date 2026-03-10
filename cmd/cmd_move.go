package cmd

import (
	"errors"
	"fmt"
	"github.com/cute-angelia/avmetagetter/config"
	"github.com/cute-angelia/avmetagetter/pkg/media"
	"github.com/cute-angelia/avmetagetter/pkg/scraper"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cute-angelia/go-xutils/components/idownload"
	"github.com/cute-angelia/go-xutils/components/loggers/loggerV3"
	"github.com/cute-angelia/go-xutils/syntax/ifile"
	"github.com/cute-angelia/go-xutils/utils/conf"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func NewCmdMove() *cli.Command {
	return &cli.Command{
		Name:        "move",
		Usage:       "move -dir={dir} -dest={dest} -env={local} -scraper={scraper}",
		Description: "扫描文件夹，将nfo等信息移动到目标文件夹",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "local",
				Usage: "环境变量",
			},
			&cli.StringFlag{
				Name:  "dir",
				Value: "./",
				Usage: "扫描文件夹",
			},
			&cli.StringFlag{
				Name:  "dest",
				Value: "./jav",
				Usage: "目标文件夹",
			},
			&cli.StringFlag{
				Name:  "scraper",
				Value: "",
				Usage: "指定scraper：JavBus JavDb CaribBeanCom FC2 TokyoHot Heyzo Heydouga Siro memojav202508等",
			},
		},

		Action: func(c *cli.Context) error {

			envstr := c.String("env")
			dirIn := c.String("dir")
			destIn := c.String("dest")
			captureNames := c.String("scraper")

			// 加载 config
			config.InitConfig(envstr)
			if ifile.IsExist("./config.toml") {
				if err := conf.MergeConfigWithPath("./config.toml"); err != nil {
					log.Println(err)
				}
			}

			crps := strings.Split(captureNames, ",")

			if len(dirIn) == 0 {
				dirIn = viper.GetString("avnas.dir")
			}
			if len(destIn) == 0 {
				destIn = viper.GetString("avnas.dest")
			}
			return fire(dirIn, destIn, crps)
		},
	}
}

func fire(dir string, dest string, crps []string) error {
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
			iscraper := scraper.NewScraper(no, viper.GetString("common.socks5"))
			if resps, err := iscraper.Search(crps); err != nil {
				loggerV3.GetLogger().Err(err).Str("抓取失败", no).Send()
				continue
			} else {

				for _, resp := range resps {

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

					// linux 最大字符 255
					title = MaxLength(title, 80)

					// 生成目标文件夹
					// 规则 [2019] STARS-065 ナマ派 初中出し解禁 本庄鈴
					destdir = fmt.Sprintf("%s/[%s] %s", actorName, nfo.Year, title)

					loggerV3.GetLogger().Info().Str("生成目标文件夹", destdir).Send()

					// 生成 inf
					nfoPath := filepath.Join(dest, destdir, fmt.Sprintf("%s.nfo", ifile.NameNoExt(avfile)))
					nfoFile, _ := ifile.CreateFile(nfoPath)
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

					// 只保存一个
					break
				}
			}
		}
		return nil
	}
}

// MaxLength 文本最大长度 已经放入go-utils
func MaxLength(txt string, size int) string {
	//将字符串转为[]rune类型
	txtRune := []rune(txt)
	fLength := len(txtRune)
	diff := fLength - size

	if diff > 0 {
		return string(txtRune[diff:fLength])
	} else {
		return txt
	}
}
