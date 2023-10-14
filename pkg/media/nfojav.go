package media

import (
	"encoding/xml"
	"fmt"
	"metagetter/pkg/sites"
	"regexp"
	"strings"
)

// NfoJav Nfo信息结构，
// 用以存储 nfo 文件所需各项信息。
type NfoJav struct {
	BuildMedia
	XMLName   xml.Name `xml:"movie"`
	Title     Inner    `xml:"title"`
	SortTitle string   `xml:"sorttitle"`
	Number    string   `xml:"num"`
	Studio    Inner    `xml:"studio"`
	Maker     Inner    `xml:"maker"`
	Director  Inner    `xml:"director"`
	Release   string   `xml:"release"`
	Premiered string   `xml:"premiered"`
	Year      string   `xml:"year"`
	Plot      Inner    `xml:"plot"`
	Outline   Inner    `xml:"outline"`
	RunTime   string   `xml:"runtime"`
	Mpaa      string   `xml:"mpaa"`
	Country   string   `xml:"country"`
	Poster    string   `xml:"poster"`
	Thumb     string   `xml:"thumb"`
	FanArt    string   `xml:"fanart"`
	Actor     []Actor  `xml:"actor"`
	Tag       []Inner  `xml:"tag"`
	Genre     []Inner  `xml:"genre"`
	Set       string   `xml:"set"`
	Label     string   `xml:"label"`
	Cover     string   `xml:"cover"`
	WebSite   string   `xml:"website"`
	Month     string   `xml:"-"`
	DirPath   string   `xml:"-"`
	Sample    []string `xml:"-"` // 样图
}

func NewNfoJav() *NfoJav {
	return &NfoJav{}
}

func (m *NfoJav) ParseMedia(resp sites.ScraperResp) {
	// 短标题
	m.SortTitle = strings.TrimSpace(resp.No)
	// 番号
	m.Number = m.SortTitle
	// 厂商
	m.Studio = Inner{Inner: strings.TrimSpace(resp.Studio)}
	// 厂商
	m.Maker = m.Studio
	// 导演
	m.Director = Inner{Inner: strings.TrimSpace(resp.Director)}
	// 发行时间
	m.Release = strings.TrimSpace(strings.ReplaceAll(resp.ReleaseDate, "/", "-"))
	// 发行时间
	m.Premiered = m.Release
	// 设置年份
	m.Year = strings.TrimSpace(GetYear(m.Release))
	// 简介
	m.Plot = Inner{Inner: resp.Intro}
	// 简介
	m.Outline = m.Plot
	// 时长
	m.RunTime = strings.TrimSpace(resp.Runtime)
	// 分级
	m.Mpaa = "XXX"
	// 国家
	m.Country = "JP"
	// 演员
	// 定义演员列表
	var actors []Actor
	// 获取演员并循环
	for name, thumb := range resp.Actors {
		// 加入列表
		actors = append(actors, Actor{
			Name:  name,
			Thumb: thumb,
		})
	}
	m.Actor = actors
	// 标签
	tags := resp.Tags
	// 循环标签
	for _, tag := range tags {
		m.Tag = append(m.Tag, Inner{Inner: tag})
	}
	// 类型
	m.Genre = m.Tag
	// 系列
	m.Set = strings.TrimSpace(resp.Series)
	// 图片
	m.Cover = strings.TrimSpace(resp.Cover)
	// 地址
	m.WebSite = strings.TrimSpace(resp.WebSite)
	// 设置月份
	m.Month = strings.TrimSpace(GetMonth(m.Release))

	// 获取标题
	title := strings.TrimSpace(resp.Title)
	// 替换原有番号
	title = strings.TrimSpace(strings.ReplaceAll(title, m.Number, ""))
	// 重新增加番号
	title = fmt.Sprintf("%s %s", m.Number, title)
	// 设置标题
	m.Title = Inner{Inner: title}

	m.Sample = resp.SampleImg
}

func (m *NfoJav) Marshal() ([]byte, error) {
	// 转换
	x, err := xml.MarshalIndent(m, "", "  ")
	// 检查
	if err != nil {
		return nil, err
	}
	// 转码为[]byte
	x = []byte(xml.Header + string(x))
	return x, nil
}

// GetYear 通过获取到的发行日期获取年份信息。
//
// date 字符串参数，传入发行日期。
func GetYear(date string) string {
	// 年份搜索正则
	re := regexp.MustCompile(`\d{4}`)

	return re.FindString(date)
}

// GetMonth 通过获取到的发行日期获取月份信息。
//
// date 字符串参数，传入发行日期。
func GetMonth(date string) string {
	// 月份搜索正则
	re := regexp.MustCompile(`\d{4}-([\d]{2})-\d{2}`)
	// 查找
	month := re.FindStringSubmatch(date)
	// 如果找到
	if len(month) > 0 {
		return month[1]
	}

	return ""
}
