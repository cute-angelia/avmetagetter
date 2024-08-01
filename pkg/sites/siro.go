package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"log"
	"strings"
)

type siro struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewSiro(no string, useragent, cookies, proxy string) *siro {
	return &siro{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "",
	}
}

func (that *siro) GetPageUri() []string {
	return []string{
		fmt.Sprintf("https://www.mgstage.com/product/product_detail/%s/", that.no),
	}
}

func (that *siro) Fetch() (resp ScraperResp, err error) {
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

		if root, err2 := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err2 != nil {
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {
			resp.No = that.no
			resp.WebSite = uri

			// 标题
			resp.Title = utils.IntroFilter(root.Find(`h1.tag`).Text())

			resp.Intro = utils.IntroFilter(root.Find(`#introduction p.introduction`).Text())

			// 获取导演
			resp.Director = ""

			resp.ReleaseDate = root.Find(`th:contains("配信開始日")`).NextFiltered("td").Text()

			resp.Runtime = strings.TrimRight(root.Find(`th:contains("収録時間")`).NextFiltered("td").Text(), "min")

			val := root.Find(`th:contains("メーカー")`).NextFiltered("td").Text()
			if val == "" {
				val = root.Find(`th:contains("メーカー")`).NextFiltered("td").Find("a").Text()
			}

			resp.Studio = utils.IntroFilter(val)

			val = root.Find(`th:contains("シリーズ")`).NextFiltered("td").Text()
			if val == "" {
				val = root.Find(`th:contains("シリーズ")`).NextFiltered("td").Find("a").Text()
			}

			resp.Series = utils.IntroFilter(val)

			// 类别数组
			var tags []string
			resp.Tags = tags

			// 获取sample图片
			sample := []string{}
			root.Find(".sample-photo a").Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Attr("href"); ok {
					sample = append(sample, v)
				}
			})
			resp.SampleImg = sample

			// 获取cover图片
			fanart, _ := root.Find(`#EnlargeImage`).Attr("href")

			// 检查
			resp.Cover = fanart

			// 演员数组
			actors := make(map[string]string)

			// 循环获取
			root.Find(`th:contains("出演")`).NextFiltered("td").Find("a").Each(func(i int, item *goquery.Selection) {
				// 演员名字
				actors[strings.TrimSpace(item.Text())] = ""
			})

			// 是否获取到
			if len(actors) == 0 {
				// 重新获取
				name := root.Find(`th:contains("出演")`).NextFiltered("td").Text()
				// 获取演员名字
				actors[strings.TrimSpace(name)] = ""
			}
			resp.Actors = actors

			if len(resp.Cover) == 0 {
				err = ErrorCoverNotFound
				continue
			}
		}
	}
	// log.Println(htmlBody)

	return resp, err
}
