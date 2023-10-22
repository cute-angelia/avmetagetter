package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"log"
	"strings"
)

type javLibrary struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewJavLibrary(no string, useragent, cookies, proxy string) *javLibrary {
	return &javLibrary{
		no:        strings.ToUpper(no),
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "http://www.javlibrary.com/",
	}
}

// GetPageUri 获取页面地址
func (that *javLibrary) GetPageUri() []string {
	return []string{fmt.Sprintf("http://www.javlibrary.com/cn/vl_searchbyid.php?keyword=%s", that.no)}
}

func (that *javLibrary) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	for _, uri := range uris {

		if !strings.Contains(uri, "http") {
			return resp, errors.New("error url address:" + uri)
		}
		var htmlBody string

		// get
		if htmlBody, err = utils.GetBody(uri, "#content"); err != nil {
			return resp, err
		}

		if root, err2 := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err2 != nil {
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {

			// 查找是否获取到
			resp.No = that.no
			resp.WebSite = uri
			resp.Title = root.Find("h3.post-title").Text()
			if len(resp.Title) == 0 {
				err = errors.New("404 Not Found")
				continue
			}

			resp.Intro = ""
			resp.Director = root.Find(`#video_director span.director`).First().Text()

			txt := strings.Trim(root.Find(`#video_date`).Text(), "\n")
			txt = strings.Replace(txt, "\n", "", -1)
			txt = strings.Replace(txt, "\t", "", -1)
			txt = strings.Replace(txt, " ", "", -1)
			txt = strings.Replace(txt, "发行日期:", "", -1)
			resp.ReleaseDate = txt

			txt = strings.Trim(root.Find(`#video_length`).Text(), "分钟")
			txt = strings.Replace(txt, "\n", "", -1)
			txt = strings.Replace(txt, "\t", "", -1)
			txt = strings.Replace(txt, " ", "", -1)
			txt = strings.Replace(txt, "长度:", "", -1)
			resp.Runtime = txt

			resp.Studio = strings.Trim(root.Find(`#video_maker span.maker`).Text(), " ")
			resp.Series = strings.Trim(root.Find(`#video_label span.label a`).Text(), " ")

			// 类别数组
			var tags []string
			// 循环获取
			root.Find(`#video_genres span.genre`).Each(func(i int, item *goquery.Selection) {
				tags = append(tags, strings.TrimSpace(item.Text()))
			})
			resp.Tags = tags

			// 获取cover图片
			// 获取图片
			fanart, _ := root.Find(`#video_jacket img`).Attr("src")
			if !strings.Contains(fanart, "http") {
				fanart = "http:" + fanart
			}
			resp.Cover = fanart

			// 获取sample图片
			sample := []string{}
			root.Find(`.previewthumbs img`).Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Attr("src"); ok {
					sample = append(sample, v)
				}
			})
			resp.SampleImg = sample

			// 演员数组
			// 演员数组
			actors := make(map[string]string)
			// 循环获取
			root.Find(`.star a`).Each(func(i int, item *goquery.Selection) {
				// 获取演员图片
				img := ""
				// 获取演员名字
				name := item.Text()

				// 加入列表
				actors[strings.TrimSpace(name)] = strings.TrimSpace(img)
			})
			resp.Actors = actors
		}
	}
	// log.Println(htmlBody)

	return resp, err
}
