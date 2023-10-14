package sites

type (
	IScraper interface {
		GetPageUri() string
		Fetch() (ScraperResp, error)
	}
	ScraperResp struct {
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
		WebSite     string
	}
	BuildInScraper IScraper
)

func (_ *ScraperResp) GetPageUri() (_ string) {
	return
}

func (_ *ScraperResp) Fetch() (resp ScraperResp, err error) {
	return
}
