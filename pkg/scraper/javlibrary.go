package scraper

import (
	"fmt"
	"strings"

	"github.com/cute-angelia/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// JavBusScraper javbus网站刮削器
type JavlibraryScraper struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// site 字符串参数，传入免翻地址，
// proxy 字符串参数，传入代理信息
func NewJavLibraryScraper(proxy string) *JavlibraryScraper {
	return &JavlibraryScraper{Site: "http://www.javlibrary.com/cn/vl_searchbyid.php?keyword=", Proxy: proxy}
}

// Fetch 刮削
func (s *JavlibraryScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 获取信息
	err := s.detail()
	// 检查错误
	if err != nil {
		// 设置番号
		s.number = strings.ReplaceAll(s.number, "-", "_")
		// 使用 _ 方式
		err = s.detail()
		// 检查错误
		if err != nil {
			// 设置番号
			s.number = strings.ReplaceAll(strings.ReplaceAll(s.number, "-", ""), "_", "")
			// 去除符号
			err = s.detail()
			// 检查错误
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 获取获取
func (s *JavlibraryScraper) detail() error {
	// 组合uri
	uri := fmt.Sprintf("%s%s", util.CheckDomainPrefix(s.Site), s.number)
	// 获取节点
	root, err := util.GetRootNewGout(uri, s.Proxy, nil, false)

	// 检查错误
	if err != nil {
		return err
	}
	s.root = root

	// 查找是否获取到
	if len(s.GetTitle()) == 0 {
		// log.Println(root.Html())
		return fmt.Errorf("404 Not Found")
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *JavlibraryScraper) GetTitle() string {
	return s.root.Find("h3.post-title").Text()
}

// GetIntro 获取简介
func (s *JavlibraryScraper) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavlibraryScraper) GetDirector() string {
	return s.root.Find(`#video_director span.director`).First().Text()
}

// GetRelease 发行时间
func (s *JavlibraryScraper) GetRelease() string {
	txt := strings.Trim(s.root.Find(`#video_date`).Text(), "\n")
	txt = strings.Replace(txt, "\n", "", -1)
	txt = strings.Replace(txt, "\t", "", -1)
	txt = strings.Replace(txt, " ", "", -1)
	txt = strings.Replace(txt, "发行日期:", "", -1)

	return txt
}

// GetRuntime 获取时长
func (s *JavlibraryScraper) GetRuntime() string {
	txt := strings.Trim(s.root.Find(`#video_length`).Text(), "分钟")
	txt = strings.Replace(txt, "\n", "", -1)
	txt = strings.Replace(txt, "\t", "", -1)
	txt = strings.Replace(txt, " ", "", -1)
	txt = strings.Replace(txt, "长度:", "", -1)

	return txt
}

// GetStudio 获取厂商
func (s *JavlibraryScraper) GetStudio() string {
	txt := strings.Trim(s.root.Find(`#video_maker span.maker`).Text(), " ")
	return txt
}

// GetSeries 获取系列
func (s *JavlibraryScraper) GetSeries() string {
	txt := strings.Trim(s.root.Find(`#video_label span.label a`).Text(), " ")
	return txt
}

// GetTags 获取标签
func (s *JavlibraryScraper) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`#video_genres span.genre`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *JavlibraryScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`#video_jacket img`).Attr("src")

	if !strings.Contains(fanart, "http") {
		fanart = "http:" + fanart
	}

	return fanart
}

// GetActors 获取演员
func (s *JavlibraryScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`.star a`).Each(func(i int, item *goquery.Selection) {
		// 获取演员图片
		img := ""
		// 获取演员名字
		name := item.Text()

		// 加入列表
		actors[strings.TrimSpace(name)] = strings.TrimSpace(img)
	})

	return actors
}

// GetURI 获取页面地址
func (s *JavlibraryScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavlibraryScraper) GetNumber() string {
	return s.number
}
