package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"log"
	"strings"
)

type fc2 struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewFc2(no string, useragent, cookies, proxy string) *fc2 {
	return &fc2{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "",
	}
}

func (that *fc2) GetPageUri() []string {
	return []string{
		fmt.Sprintf("https://adult.contents.fc2.com/article/%s/", strings.Replace(strings.Replace(that.no, "FC2-", "", 1), "fc2", "", 1)),
		fmt.Sprintf("https://fc2club.top/html/%s.html", that.no),
	}
}

func (that *fc2) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	for _, uri := range uris {
		if !strings.Contains(uri, "http") {
			err = errors.New("error url address:" + uri)
			log.Println(err)
			continue
		}

		log.Println(uri)

		var htmlBody string
		// get
		utils.GetIGout(uri, that.proxy, false).SetHeader(gout.H{
			"User-Agent": that.useragent,
			"Cookie":     that.cookies,
			"referer":    that.site,
		}).BindBody(&htmlBody).Do()

		if root, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err != nil {
			log.Println("ERROR:", err)
			continue
		} else {
			// 查找是否获取到
			if -1 == root.Find(`h3`).Index() {
				err = errors.New("404 Not Found")
				continue
			}

			resp.No = that.no
			resp.WebSite = uri

			// 标题
			title := root.Find(`h3`).First().Text()
			resp.Title = title

			resp.Intro = ""

			// 获取导演
			director := root.Find(`.items_article_headerInfo li:nth-child(3) a`).Text()
			// 检查
			if director == "" {
				director = root.Find(`.main h5:nth-child(5) a:nth-child(2)`).Text()
			}
			resp.Director = director

			resp.ReleaseDate = strings.ReplaceAll(strings.ReplaceAll(root.Find(`.items_article_Releasedate p`).Text(), "上架时间 :", ""), "販売日 :", "")

			resp.Runtime = ""
			resp.Studio = "FC2"
			resp.Series = "FC2"

			// 类别数组
			var tags []string
			resp.Tags = tags

			// 获取sample图片
			sample := []string{}
			root.Find(`.items_article_SampleImagesArea a`).Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Attr("href"); ok {
					sample = append(sample, v)
				}
			})
			resp.SampleImg = sample

			// 获取cover图片
			var fanart string
			if len(resp.SampleImg) > 0 {
				fanart = resp.SampleImg[0]
			}
			// 检查
			resp.Cover = fanart

			// 演员数组
			actors := make(map[string]string)
			resp.Actors = actors

			if len(resp.Title) > 10 {
				break
			}
		}
	}
	// log.Println(htmlBody)

	return resp, nil
}
