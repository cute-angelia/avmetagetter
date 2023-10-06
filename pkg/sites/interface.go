package sites

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
