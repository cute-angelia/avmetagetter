package media


import (
	"encoding/xml"
	"fmt"
	"github.com/cute-angelia/go-utils/utils/generator/hash"
)

// Inner 文字数据，为了避免某些内容被转义。
type Inner struct {
	Inner string `xml:",innerxml"`
}

// Actor 演员信息，保存演员姓名及头像地址。
type Actor struct {
	Name  string `xml:"name"`
	Thumb string `xml:"thumb"`
}

type Art struct {
	Poster string `xml:"poster"`
	Fanart string `xml:"fanart"`
}

type JellyfinMeta struct {
	Tmdbid int    `xml:"tmdbid"` // tmdbid
	Key    string `xml:"key"`    // cache key
	Update string `xml:"update"`
}

// Nfo 信息结构，
// 用以存储 nfo 文件所需各项信息。
// https://www.tinymediamanager.org/docs/movies/settings
type Nfo struct {
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
	Plot      string   `xml:"plot"`
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
	Art       Art      `xml:"art"`

	Imdbid         string `xml:"imdbid"`
	Tmdbid         int    `xml:"tmdbid"`         // tmdbid
	Opendoubanidid int    `xml:"opendoubanidid"` // 豆瓣

	JellyfinMeta JellyfinMeta `xml:"jellyfinmeta"` // 自建脚本
	NfoPath      string       `xml:"nfopath"`      // 自建路径，用于后期修改文件保存
}

func NewNfo() *Nfo {
	return &Nfo{}
}

func (n *Nfo) Marshal() ([]byte, error) {
	// 转换
	x, err := xml.MarshalIndent(n, "", "  ")
	// 检查
	if err != nil {
		return nil, err
	}
	// 转码为[]byte
	x = []byte(xml.Header + string(x))
	return x, nil
}

// CalcJellyfinMetaCacheKey 简单根据 tmbid 和 tmdbid 的变化计算
func (n *Nfo) CalcJellyfinMetaCacheKey() string {
	str := fmt.Sprintf("%d%s", n.Tmdbid, n.Imdbid)

	return hash.Hash(hash.AlgoSha1, str)
}
