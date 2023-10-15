package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type javbus struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewJavBus(no string, useragent, cookies, proxy string) *javbus {
	return &javbus{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      viper.GetString("javbus.site"),
	}
}

func (that *javbus) GetPageUri() []string {
	if len(that.site) == 0 {
		that.site = "https://www.javbus.com/"
	}
	return []string{that.site + that.no}
}

func (that *javbus) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	for _, uri := range uris {

		if !strings.Contains(uri, "http") {
			err = errors.New("error url address:" + uri)
			continue
		}
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
			resp.Title = root.Find("h3").Text()
			resp.Intro = ""
			resp.Director = root.Find(`a[href*="/director/"]`).Text()
			resp.ReleaseDate = strings.ReplaceAll(root.Find(`p:contains("發行日期:")`).Text(), "發行日期: ", "")
			resp.Runtime = strings.ReplaceAll(strings.TrimRight(root.Find(`p:contains("長度:")`).Text(), "分鐘"), "長度: ", "")
			resp.Studio = root.Find(`a[href*="/studio/"]`).Text()
			resp.Series = root.Find(`a[href*="/series/"]`).Text()

			// 类别数组
			var tags []string
			// 循环获取
			root.Find(`span.genre a[href*="/genre/"]`).Each(func(i int, item *goquery.Selection) {
				tags = append(tags, utils.T2S(strings.TrimSpace(item.Text())))
			})
			resp.Tags = tags

			// 获取cover图片
			fanart, _ := root.Find(`a.bigImage img`).Attr("src")
			resp.Cover = fmt.Sprintf("%s%s", that.site, strings.TrimLeft(fanart, "/"))

			// 获取sample图片
			sample := []string{}
			root.Find(`a.sample-box`).Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Attr("href"); ok {
					sample = append(sample, v)
				}
			})
			resp.SampleImg = sample

			// 演员数组
			actors := make(map[string]string)
			// 循环获取
			root.Find(`div.star-box li > a`).Each(func(i int, item *goquery.Selection) {
				// 获取演员图片
				img, _ := item.Find(`img`).Attr("src")
				img = fmt.Sprintf("%s%s", that.site, strings.TrimLeft(img, "/"))

				// 获取演员名字
				name, _ := item.Find("img").Attr("title")
				// 加入列表
				actors[strings.TrimSpace(name)] = strings.TrimSpace(img)
			})
			resp.Actors = actors
		}
	}

	// log.Println(htmlBody)

	return resp, nil
}
