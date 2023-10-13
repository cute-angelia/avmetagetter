package sites

import (
	"github.com/spf13/viper"
	"regexp"
)

type SiteResp struct {
	No          string
	Title       string            // 标题
	Intro       string            // 简介
	Director    string            // 导演
	ReleaseDate string            // 发行时间
	Runtime     string            // 时长
	Studio      string            // 获取厂商
	Series      string            // 系列
	Tags        []string          // 标签
	Cover       string            // 封面
	SampleImg   []string          // 获取样图
	Actors      map[string]string // 演员：name : 头像

}

type SiteCommon interface {
	Fetch() (SiteResp, error)
}

var regMap = map[string][]string{
	"javbus": {`[a-zA-Z]+-\d+`},
}

func GetSiteInfo(no string) (resp SiteResp, err error) {
	sites := []string{}
	for t, regs := range regMap {
		for _, reg := range regs {
			matched, _ := regexp.MatchString(reg, no)
			if matched {
				sites = append(sites, t)
			}
		}
	}

	for _, site := range sites {
		if resp, err = getSiteObj(site, no).Fetch(); err != nil {
			continue
		} else {
			return
		}
	}

	return
}

func getSiteObj(site string, no string) SiteCommon {
	var s SiteCommon
	switch site {
	case "javbus":
		s = NewJavBus(no, viper.GetString(site+".useragent"), viper.GetString(site+".cookies"), viper.GetString("common.socks5"))
		break
	default:
	}
	return s
}
