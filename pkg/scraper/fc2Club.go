package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// FC2Scraper fc2网站刮削器
type Fc2ClubScraper struct {
	Proxy  string            // 代理设置
	uri    string            // 页面地址
	code   string            // 临时番号
	number string            // 最终番号
	root   *goquery.Document // fc2根节点
}

// NewFC2Scraper 返回一个被初始化的fc2刮削对象
//
// proxy 字符串参数，传入代理信息
func NewFC2ClubScraper(proxy string) *Fc2ClubScraper {
	return &Fc2ClubScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *Fc2ClubScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 过滤番号
	r := regexp.MustCompile(`[0-9]{6,7}`)
	// 获取临时番号
	s.code = r.FindString(code)

	fc2cluburi := fmt.Sprintf("https://fc2club.net/html/FC2-%s.html", s.code)

	// 打开fc2
	root, err := util.GetRoot(fc2cluburi, s.Proxy, nil)
	// 检查错误
	if err != nil {
		log.Println(err)
		return err
	}

	// 设置页面地址
	s.uri = fc2cluburi
	// 设置fc2club根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *Fc2ClubScraper) GetTitle() string {
	// 获取标题
	title := s.root.Find(`.main h3`).Text()

	return title
}

// GetIntro 获取简介
func (s *Fc2ClubScraper) GetIntro() string {
	return ""
}

// GetDirector 获取导演
func (s *Fc2ClubScraper) GetDirector() string {
	// 获取导演
	director := s.root.Find(`.main h5:nth-child(5) a:nth-child(2)`).Text()

	return director
}

// GetRelease 发行时间
func (s *Fc2ClubScraper) GetRelease() string {
	return strings.ReplaceAll(strings.ReplaceAll(s.root.Find(`.items_article_Releasedate p`).Text(), "上架时间 :", ""), "販売日 :", "")
}

// GetRuntime 获取时长
func (s *Fc2ClubScraper) GetRuntime() string {
	return "0"
}

// GetStudio 获取厂商
func (s *Fc2ClubScraper) GetStudio() string {
	return util.FC2
}

// GetSeries 获取系列
func (s *Fc2ClubScraper) GetSeries() string {
	return util.FC2
}

// GetTags 获取标签
func (s *Fc2ClubScraper) GetTags() []string {
	// 组合地址
	uri := fmt.Sprintf("http://adult.contents.fc2.com/api/v4/article/%s/tag?", s.code)

	// 读取远程数据
	data, err := util.GetResult(uri, s.Proxy, nil)
	// 检查
	if err != nil {
		return nil
	}

	// 读取内容
	body, err := ioutil.ReadAll(bytes.NewReader(data))
	// 检查错误
	if err != nil {
		return nil
	}

	// json
	var tagsJSON fc2tags

	// 解析json
	err = json.Unmarshal(body, &tagsJSON)
	// 检查
	if err != nil {
		return nil
	}

	// 定义数组
	var tags []string

	// 循环标签
	for _, tag := range tagsJSON.Tags {
		tags = append(tags, strings.TrimSpace(tag.Tag))
	}

	return tags
}

// GetCover 获取图片
func (s *Fc2ClubScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`.slides li:nth-child(1) img`).Attr("src")
	// 检查
	if fanart == "" {
		return ""
	}
	// 组合地址
	return fmt.Sprintf("https://fc2club.net%s", fanart)
}

// GetActors 获取演员
func (s *Fc2ClubScraper) GetActors() map[string]string {
	return nil
}

// GetURI 获取页面地址
func (s *Fc2ClubScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *Fc2ClubScraper) GetNumber() string {
	return s.number
}

// 获取样图
func (s *Fc2ClubScraper) GetSample() []string {
	// 获取图片
	sample := []string{}

	s.root.Find(`.items_article_SampleImages a`).Each(func(i int, selection *goquery.Selection) {
		if v, ok := selection.Attr("href"); ok {
			sample = append(sample, v)
		}
	})
	return sample
}
