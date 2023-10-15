package scraper

import (
	"github.com/cute-angelia/go-utils/components/loggers/loggerV3"
	"github.com/cute-angelia/metagetter/pkg/sites"
	"github.com/spf13/viper"
	"regexp"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
)

type scraper struct {
	no           string
	proxy        string
	captureNames []string
}

// 刮削对象
type captures struct {
	Name         string
	Scraper      sites.IScraper
	Reg          *regexp.Regexp
	Enable       bool
	NeedChromeDp bool // 需要安装chromeDp
}

func NewScraper(no string, proxy string, captureNames []string) *scraper {
	return &scraper{
		no:           no,
		proxy:        proxy,
		captureNames: captureNames,
	}
}

func (that *scraper) getCaptures() []captures {
	// 定义一个拥有正则匹配的刮削对象数组
	cs := []captures{
		{
			Name:    "JavBus",
			Scraper: sites.NewJavBus(that.no, viper.GetString("javbus.useragent"), viper.GetString("javbus.cookies"), that.proxy),
			Reg:     regexp.MustCompile(`^[a-zA-Z]+-\d{2,10}$`),
			Enable:  true,
		},
		{
			Name:         "JavLibrary",
			Scraper:      sites.NewJavLibrary(that.no, DefaultUserAgent, "", that.proxy),
			Reg:          regexp.MustCompile(`^[a-zA-Z]+-\d{2,10}$`),
			NeedChromeDp: true,
			Enable:       false, // 关闭
		},
		{
			Name:         "JavDb",
			Scraper:      sites.NewJavDb(that.no, DefaultUserAgent, "", that.proxy),
			Reg:          regexp.MustCompile(`^[a-zA-Z-0-9]{2,15}$`),
			Enable:       true,
			NeedChromeDp: true,
		},
		{
			Name:    "CaribBeanCom",
			Scraper: sites.NewCaribBeanCom(that.no, DefaultUserAgent, "", that.proxy),
			Reg:     regexp.MustCompile(`^\d{6}-\d{3}$`),
			Enable:  true,
		},
		{
			Name:    "DMM",
			Scraper: sites.NewDmm(that.no, DefaultUserAgent, "", that.proxy),
			Reg:     regexp.MustCompile(`[a-zA-Z]{2,5}[-|\s\S][0-9]{3,4}`),
			Enable:  false,
		},
		{
			Name:    "FC2",
			Scraper: sites.NewFc2(that.no, DefaultUserAgent, "", that.proxy),
			Reg:     regexp.MustCompile(`^(fc2|FC2)-[0-9]{6,8}`),
			Enable:  true,
		},
		{
			Name:    "TokyoHot",
			Scraper: sites.NewTokyohot(that.no, DefaultUserAgent, "", that.proxy),
			Reg:     regexp.MustCompile(`(^red-\d{3}|n\d{4})`),
			Enable:  true,
		},
		//{
		//	Name: "Heyzo",
		//	S:    scraper.NewHeyzoScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`^heyzo-[0-9]{4}`),
		//},
		//{
		//	Name: "Heydouga",
		//	S:    scraper.NewHeydougaScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`([0-9]{4}).+?([0-9]{3,4})$`),
		//},

		//{
		//	Name: "FC2Club",
		//	S:    scraper.NewFC2ClubScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`^fc2-[0-9]{6,7}`),
		//},

		//{
		//	Name: "Siro",
		//	S:    scraper.NewSiroScraper(cfg.Base.Proxy),
		//	//R:    regexp.MustCompile(`^(siro|abw|abp|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`),
		//	R: regexp.MustCompile(`^([a-zA-Z]{2,6}|[0-9]{3,5}[a-zA-Z]{2,6})-[0-9]{3,4}`),
		//},

	}
	if len(that.captureNames) > 0 && len(that.captureNames[0]) > 0 {
		var cs2 []captures
		for _, c := range cs {
			for _, name := range that.captureNames {
				if c.Name == name {
					cs2 = append(cs2, c)
				}
			}
		}
		return cs2
	} else {
		return cs
	}
}

func (that *scraper) Search() (resp sites.ScraperResp, err error) {
	icaptures := that.getCaptures()
	for _, item := range icaptures {
		if item.Enable && item.Reg.MatchString(that.no) {
			loggerV3.GetLogger().Info().Str("Matched", that.no).Str("site", item.Name).Send()
			if resp, err = item.Scraper.Fetch(); err != nil {
				continue
			} else {
				return
			}
		}
	}
	return
}
