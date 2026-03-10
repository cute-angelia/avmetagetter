package sites

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"avmetagetter/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/guonaihong/gout"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type caribBeanComPr struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewCaribBeanComPr(no string, useragent, cookies, proxy string) *caribBeanComPr {
	return &caribBeanComPr{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://www.caribbeancompr.com",
	}
}

func (that *caribBeanComPr) GetPageUri() []string {
	return []string{fmt.Sprintf("%s/moviepages/%s/index.html", that.site, that.no)}
}

func (that *caribBeanComPr) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	for _, uri := range uris {
		var htmlBody string
		// get
		utils.GetIGout(uri, that.proxy, false).SetHeader(gout.H{
			"User-Agent": that.useragent,
			"Cookie":     that.cookies,
			"referer":    that.site,
		}).BindBody(&htmlBody).Do()

		// 编码转换
		reader := transform.NewReader(strings.NewReader(htmlBody), japanese.EUCJP.NewDecoder())

		if root, err2 := goquery.NewDocumentFromReader(reader); err2 != nil {
			log.Println("ERROR Reader:", err2)
			err = err2
			continue
		} else {
			// 查找是否获取到
			if -1 == root.Find(`h1`).Index() {
				err = errors.New("404 Not Found")
				continue
			}

			resp.No = that.no
			resp.WebSite = uri
			resp.Title = root.Find(`h1`).First().Text()

			// 获取简介
			resp.Intro = root.Find(`.movie-info p`).First().Text()

			root.Find(".movie-info li").Each(func(i int, s *goquery.Selection) {
				//log.Println(i, strings.TrimSpace(s.Find(".spec-content").Text()))

				// 演员
				if i == 0 {
					// 演员数组
					actors := make(map[string]string, 10)
					s.Find("a").Each(func(i int, s *goquery.Selection) {
						actors[strings.TrimSpace(s.Text())] = ""
					})
					resp.Actors = actors
				}

				//
				if i == 1 {
					resp.Runtime = strings.TrimSpace(s.Find(".spec-content").Text())
				}
				if i == 2 {
					resp.Studio = strings.TrimSpace(s.Find(".spec-content").Text())
				}
				if i == 3 {
					strings.TrimSpace(s.Find(".spec-content").Text())
				}
				if i == 4 {
					s.Find("a").Each(func(i int, s *goquery.Selection) {
						resp.Tags = append(resp.Tags, strings.TrimSpace(s.Text()))
					})
				}
			})

			resp.Director = ""
			resp.ReleaseDate = ""

			// 获取cover图片
			resp.Cover = that.site + "/moviepages/" + that.no + "/images/l_l.jpg"

			// 获取sample图片
			sample := []string{}
			root.Find(`.grid-item`).Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Find("a").First().Attr("href"); ok {
					if strings.Contains(v, "jpg") {
						sample = append(sample, v)
					}
				}
			})
			resp.SampleImg = sample

			if len(resp.Cover) == 0 {
				err = ErrorCoverNotFound
				continue
			}
			if len(resp.Title) < 10 {
				err = errors.New("title not right")
				continue
			}
		}
	}
	// log.Println(htmlBody)

	return resp, err
}
