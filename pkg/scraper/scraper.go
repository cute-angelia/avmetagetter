package scraper

import (
	"github.com/cute-angelia/go-utils/components/loggers/loggerV3"
	"github.com/spf13/viper"
	"metagetter/pkg/sites"
	"regexp"
)

type scraper struct {
	no    string
	proxy string
}

// 刮削对象
type captures struct {
	Name string
	S    sites.IScraper
	R    *regexp.Regexp
}

func NewScraper(no string, proxy string) *scraper {
	return &scraper{
		no:    no,
		proxy: proxy,
	}
}

func (that *scraper) initCapture() []captures {
	// 定义一个拥有正则匹配的刮削对象数组
	return []captures{
		//{
		//	Name: "CaribBeanCom",
		//	S:    scraper.NewCaribBeanComScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`^\d{6}-\d{3}$`),
		//},
		//{
		//	Name: "TokyoHot",
		//	S:    scraper.NewTokyoHotScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`(^red-\d{3}|n\d{4})`),
		//},
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
		//	Name: "FC2",
		//	S:    scraper.NewFC2Scraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`^fc2-[0-9]{6,7}`),
		//},
		//{
		//	Name: "FC2Club",
		//	S:    scraper.NewFC2ClubScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`^fc2-[0-9]{6,7}`),
		//},
		{
			Name: "JavBus",
			S:    sites.NewJavBus(that.no, viper.GetString("javbus.useragent"), viper.GetString("javbus.cookies"), that.proxy),
			R:    regexp.MustCompile(`^[a-zA-Z]+-\d{2,10}$`),
		},
		//{
		//	Name: "Javlibrary",
		//	S:    scraper.NewJavLibraryScraper(cfg.Base.Socket),
		//	R:    regexp.MustCompile(`^[a-zA-Z]+-\d{2,10}$`),
		//},
		//{
		//	Name: "Siro",
		//	S:    scraper.NewSiroScraper(cfg.Base.Proxy),
		//	//R:    regexp.MustCompile(`^(siro|abw|abp|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`),
		//	R: regexp.MustCompile(`^([a-zA-Z]{2,6}|[0-9]{3,5}[a-zA-Z]{2,6})-[0-9]{3,4}`),
		//},
		//{
		//	Name: "DMM",
		//	S:    scraper.NewDMMScraper(cfg.Base.Proxy),
		//	R:    regexp.MustCompile(`[a-zA-Z]{2,5}[-|\s\S][0-9]{3,4}`),
		//},
	}
}

func (that *scraper) Search() (resp sites.ScraperResp, err error) {
	icaptures := that.initCapture()
	for _, item := range icaptures {
		if item.R.MatchString(that.no) {
			loggerV3.GetLogger().Info().Str("Matched", that.no).Str("site", item.Name).Send()
			if resp, err = item.S.Fetch(); err != nil {
				continue
			} else {
				return
			}
		}
	}
	return
}
