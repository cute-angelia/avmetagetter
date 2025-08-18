package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/cute-angelia/go-utils/syntax/itime"
	"github.com/cute-angelia/go-utils/syntax/iurl"
	"github.com/guonaihong/gout"
	"log"
	"strings"
)

// test : go run main.go -no=START-353 -scraper=memojav202508

type memojav202508 struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewMemojav202508(no string, useragent, cookies, proxy string) *memojav202508 {
	return &memojav202508{
		no:        strings.ToUpper(no),
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://memojav.com/",
	}
}

// GetPageUri 获取页面地址
func (that *memojav202508) GetPageUri() []string {
	return []string{fmt.Sprintf("%s/video/%s", iurl.GetDomainWithOutSlant(that.site), that.no)}
}

func (that *memojav202508) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	for _, uri := range uris {

		if !strings.Contains(uri, "http") {
			err = errors.New("error url address:" + uri)
			continue
		}
		var htmlBody string
		// get
		//code := 0
		//utils.GetIGout(uri, that.proxy, true).SetHeader(gout.H{
		//	"User-Agent": that.useragent,
		//	"Referer":    that.site,
		//}).BindBody(&htmlBody).Code(&code).Do()

		// get
		err = utils.GetIGout(uri, that.proxy, false).SetHeader(gout.H{
			"User-Agent": that.useragent,
			"Cookie":     that.cookies,
			"referer":    that.site,
		}).BindBody(&htmlBody).Do()

		if root, err2 := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err2 != nil {
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {
			//log.Println(htmlBody)
			// 查找是否获取到
			if strings.Contains(root.Find("body").Text(), "It may have moved or does not exist") {
				err = errors.New("404 Not Found")
				continue
			}

			// 查找是否获取到
			resp.No = that.no
			resp.WebSite = uri

			t := root.Find(`#title`).Text()
			resp.Title = strings.TrimSpace(t)

			if len(resp.Title) < 6 {
				err = errors.New("title not right")
				break
			}

			// 简介
			resp.Intro = root.Find(`#title-description`).Text()

			//log.Println(htmlBodyContent)
			root.Find(".details").Each(func(i int, selectiondetails *goquery.Selection) {
				selectiondetails.Find("tr").Each(func(i int, selection *goquery.Selection) {

					// 导演信息
					if strings.Contains(selection.Find("th").First().Text(), "Director") {
						resp.Director = selection.Find("td").First().Text()
					}

					//// 时长
					//if i == 2 {
					//	resp.Runtime = strings.TrimRight(selection.Find("span").Text(), " minute(s)")
					//}
					//

					// 片商
					if strings.Contains(selection.Find("th").First().Text(), "Studio") {
						resp.Studio = selection.Find("td").First().Text()
					}

					// Label
					if strings.Contains(selection.Find("th").First().Text(), "Label") {
						resp.Label = selection.Find("td").First().Text()
					}

					// 系列
					if strings.Contains(selection.Find("th").First().Text(), "Series") {
						resp.Series = selection.Find("td a").First().Text()
					}

					// 标签
					if strings.Contains(selection.Find("th").First().Text(), "Categories") {
						var tags []string
						selection.Find("td a").Each(func(i int, selection *goquery.Selection) {
							tags = append(tags, strings.TrimSpace(selection.Text()))
						})
						resp.Tags = tags
					}

					// actor
					if strings.Contains(selection.Find("th").First().Text(), "Actress") {
						actors := make(map[string]string)
						selection.Find("td a").Each(func(i int, selection *goquery.Selection) {
							// 演员列表
							actors[strings.TrimSpace(selection.Text())] = ""
						})
						resp.Actors = actors
					}
				})
			})

			// 发布日期
			root.Find(".details-block").Each(func(i int, selection *goquery.Selection) {
				if strings.Contains(selection.Find("h2").Text(), "Release Date") {
					tz := selection.Find("p").Text()
					tzz := itime.NewFormatLayout(tz, "2006/01/02")
					resp.ReleaseDate = tzz.FormatDate()
				}
			})

			// 获取cover图片
			// 获取图片
			fanart, _ := root.Find(`#poster`).Attr("src")
			resp.Cover = fanart

			//// 获取sample图片
			samples := []string{}
			root.Find(".thumb-list img").Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("src")
				samples = append(samples, href)
			})
			root.Find(".thumb-preview img").Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("src")

				href = strings.ReplaceAll(href, "_s.jpg", ".jpg")

				samples = append(samples, href)
			})
			resp.SampleImg = samples

			if len(resp.Cover) == 0 {
				err = ErrorCoverNotFound
				continue
			}

			log.Println("✅memojav202508:", that.no, resp.WebSite)
		}

		// log.Println(htmlBody)

	}
	return resp, err
}
