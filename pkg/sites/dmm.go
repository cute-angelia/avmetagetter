package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"log"
	"strings"
)

type dmm struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewDmm(no string, useragent, cookies, proxy string) *dmm {
	return &dmm{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://www.dmm.co.jp",
	}
}

func (that *dmm) GetPageUri() []string {
	return []string{
		"https://www.dmm.co.jp/digital/videoa/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/mono/dvd/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/digital/anime/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/mono/anime/-/detail/=/cid=%s",
	}
}

func (that *dmm) Fetch() (resp ScraperResp, err error) {

	// 组合地址列表
	uris := that.GetPageUri()

	for _, uri := range uris {

		log.Println(fmt.Sprintf(uri, that.no))

		var htmlBody string
		// get
		if err := utils.GetIGout(fmt.Sprintf(uri, that.no), that.proxy, false).SetHeader(gout.H{
			"User-Agent": that.useragent,
			"Cookie":     "age_check_done:1;",
			"referer":    that.site,
		}).BindBody(&htmlBody).Do(); err != nil {
			log.Println(err)
			continue
		}

		// 编码转换
		reader := transform.NewReader(strings.NewReader(htmlBody), japanese.EUCJP.NewDecoder())

		if root, err2 := goquery.NewDocumentFromReader(reader); err2 != nil {
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {
			// 判断是否返回了地域限制
			foreignError := root.Find(`.foreignError__desc`).Text()
			if foreignError != "" {
				err = fmt.Errorf(foreignError)
				continue
			}

			// 查找是否获取到
			if -1 == root.Find(`h3`).Index() {
				err = errors.New("404 Not Found")
				continue
			}

			resp.No = that.no
			resp.WebSite = uri
			resp.Title = root.Find(`h1#title`).Text()

			// 获取简介
			resp.Intro = utils.IntroFilter(root.Find(`tr td div.mg-b20.lh4 p.mg-b20`).Text())

			// 获取导演
			director := root.Find(`td:contains("監督：")`).Next().Find("a").Text()
			// 如果没有
			if director == "" {
				director = root.Find(`td:contains("監督：")`).Next().Text()
			}

			resp.Director = director

			// 获取发行时间
			release := root.Find(`td:contains("発売日：")`).Next().Find("a").Text()
			// 没获取到
			if release == "" {
				release = root.Find(`td:contains("発売日：")`).Next().Text()
			}
			// 替换
			resp.ReleaseDate = strings.ReplaceAll(release, "/", "-")

			resp.Runtime = strings.TrimRight(root.Find(`td:contains("収録時間：")`).Next().Text(), "分")

			// 获取厂商
			studio := root.Find(`td:contains("メーカー：")`).Next().Find("a").Text()
			// 是否获取到
			if studio == "" {
				studio = root.Find(`td:contains("メーカー：")`).Next().Text()
			}

			resp.Studio = studio

			// 获取系列
			set := root.Find(`td:contains("シリーズ：")`).Next().Find("a").Text()
			// 是否获取到
			if set == "" {
				set = root.Find(`td:contains("シリーズ：")`).Next().Text()
			}
			resp.Series = set

			// 类别数组
			var tags []string
			// 循环获取
			root.Find(`td:contains("ジャンル：")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
				// 加入数组
				tags = append(tags, strings.TrimSpace(item.Text()))
			})
			resp.Tags = tags

			// 获取cover图片
			// 获取图片
			fanart, _ := root.Find(`#` + that.no).Attr("href")

			if fanart == "" {
				root.Find(`td:contains("品番：")`).Next().Each(func(i int, item *goquery.Selection) {
					// 获取演员名字
					number := strings.TrimSpace(item.Text())
					if number != "" {
						cover, _ := root.Find(`#` + number).Attr("href")
						fanart = cover
					}
				})
			}
			resp.Cover = fanart

			// 获取sample图片
			sample := []string{}
			//root.Find(`div.gallery .grid-item`).Each(func(i int, selection *goquery.Selection) {
			//	if v, ok := selection.Find("a").First().Attr("href"); ok {
			//		v = that.site + v
			//		v = strings.Replace(v, "/member", "", 1)
			//		sample = append(sample, v)
			//	}
			//})
			resp.SampleImg = sample

			// 演员数组
			actors := make(map[string]string)
			// 循环获取
			root.Find(`td:contains("出演者：")`).Next().Find("span a").Each(func(i int, item *goquery.Selection) {
				// 获取演员名字
				actors[strings.TrimSpace(item.Text())] = ""
			})
			resp.Actors = actors

			if len(resp.Title) < 10 {
				err = errors.New("title not right")
				continue
			}
		}

		// log.Println(htmlBody)
	}
	return resp, err
}
