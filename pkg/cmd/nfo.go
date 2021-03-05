package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/media"
	"github.com/ylqjgm/AVMeta/pkg/util"
	"log"
	"strings"
	"time"
)

var nfo_title string
var nfo_no string
var nfo_release string
var nfo_actor string
var nfo_tag string
var nfo_plot string
var nfo_mpaa string
var nfo_cover string

// nfo 生成
func (e *Executor) initNfo() {
	nfoCmd := &cobra.Command{
		Use:  "nfo",
		Long: `生成无法削刮的nfo`,
		Example: `  AVMeta nfo
  --title
  --no
  --actor
  --tag a,b
  --desc
  --mpaa
  --cover
  `,
		Run: e.nfoRunFunc,
	}

	// 添加参数
	// title no release
	nfoCmd.Flags().StringVar(&nfo_title, "title", "", "标题")
	nfoCmd.Flags().StringVar(&nfo_no, "no", "", "no")
	nfoCmd.Flags().StringVar(&nfo_actor, "actor", "", "演员")
	nfoCmd.Flags().StringVar(&nfo_tag, "tag", "", "tag")
	nfoCmd.Flags().StringVar(&nfo_plot, "desc", "", "简介")
	nfoCmd.Flags().StringVar(&nfo_mpaa, "mpaa", "xxx", "评级")
	nfoCmd.Flags().StringVar(&nfo_cover, "cover", "", "封面")

	re := time.Now().Format("2006-01-02 15:04:05")
	nfoCmd.Flags().StringVar(&nfo_release, "release", re, "no")

	e.rootCmd.AddCommand(nfoCmd)
}

// 头像执行命令
func (e *Executor) nfoRunFunc(cmd *cobra.Command, args []string) {
	// 定义参数变量
	//var arg string
	//down := false

	// 检测参数
	//if len(args) > 1 {
	//	// 输出帮助
	//	_ = cmd.Help()
	//	return
	//} else if len(args) > 0 {
	//	// 获取参数
	//	arg = args[0]
	//}

	mediainfo := media.Media{
		Title: media.Inner{
			Inner: nfo_title,
		},
		Plot: media.Inner{
			Inner: nfo_plot,
		},
		Number:    nfo_no,
		SortTitle: nfo_no,
		Release:   nfo_release,
	}

	// 简介
	mediainfo.Outline = mediainfo.Plot

	nfo_actors := strings.Split(nfo_actor, ",")
	for _, nfo_actor := range nfo_actors {
		mediainfo.Actor = append(mediainfo.Actor, media.Actor{Name: nfo_actor})
	}

	nfo_tags := strings.Split(nfo_tag, ",")
	for _, tag := range nfo_tags {
		mediainfo.Tag = append(mediainfo.Tag, media.Inner{Inner: tag})
	}

	// 类型
	mediainfo.Genre = mediainfo.Tag

	// 发行时间
	mediainfo.Premiered = mediainfo.Release
	// 设置年份
	mediainfo.Year = strings.TrimSpace(media.GetYear(mediainfo.Release))
	mediainfo.Month = strings.TrimSpace(media.GetMonth(mediainfo.Release))

	mediainfo.Poster = "poster.jpg"
	mediainfo.Thumb = "poster.jpg"
	mediainfo.FanArt = "fanart.jpg"

	if len(nfo_cover) > 0 {
		mediainfo.Cover = nfo_cover
	} else {
		mediainfo.Cover = "fanart.jpg"
	}

	// 写入nfo
	// 转换为XML
	if buff, err := media.MediaToXML(&mediainfo); err != nil {
		log.Println("buff:MediaToXML failed", err)
	} else {
		saveNfoPath := fmt.Sprintf("%s/%s.nfo", util.GetRunPath(), mediainfo.Number)
		if err := util.WriteFile(saveNfoPath, buff); err != nil {
			log.Println(err)
		}
	}
}
