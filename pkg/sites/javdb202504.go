package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/cute-angelia/go-utils/syntax/iurl"
	"github.com/guonaihong/gout"
	"log"
	"strings"
)

type javDb202504 struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewJavDb202504(no string, useragent, cookies, proxy string) *javDb202504 {
	return &javDb202504{
		no:        strings.ToUpper(no),
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://javdb.com/",
	}
}

// GetPageUri 获取页面地址
func (that *javDb202504) GetPageUri() []string {
	return []string{fmt.Sprintf("%s/search?q=%s&f=all", iurl.GetDomainWithOutSlant(that.site), that.no)}
}

func (that *javDb202504) Fetch() (resp ScraperResp, err error) {
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
			if -1 < root.Find(`.empty-message:contains("暫無內容")`).Index() {
				err = errors.New("404 Not Found")
				continue
			}

			// 查找是否获取到
			resp.No = that.no
			resp.WebSite = uri

			t := ""
			root.Find(`.movie-list .video-title`).Each(func(i int, selection *goquery.Selection) {
				titlefind := strings.ToUpper(selection.Text())
				//log.Println(titlefind, "xxtitlefind")
				if strings.Contains(titlefind, strings.ToUpper(that.no)) {
					t += selection.Text() + " "
					contentHref, _ := selection.Parent().Attr("href")
					resp.WebSite = fmt.Sprintf("%s%s?locale=zh", iurl.GetDomainWithOutSlant(that.site), contentHref)
				}
			})
			resp.Title = strings.TrimSpace(t)

			if len(resp.Title) < 6 {
				err = errors.New("title not right")
				break
			}

			// 获取内页信息
			log.Println(resp.WebSite)
			var htmlBodyContent string
			utils.GetIGout(resp.WebSite, that.proxy, false).SetHeader(gout.H{
				"User-Agent": that.useragent,
				"Cookie":     that.cookies,
				"referer":    that.site,
			}).BindBody(&htmlBodyContent).Do()
			root2, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlBodyContent))

			//log.Println(htmlBodyContent)
			root2.Find(".movie-panel-info div").Each(func(i int, selection *goquery.Selection) {

				// 导演信息
				if i == 3 {
					resp.Director = selection.Find("span").Text()
				}
				// 发布日期
				if i == 1 {
					resp.ReleaseDate = selection.Find("span").Text()
				}
				// 时长
				if i == 2 {
					resp.Runtime = strings.TrimRight(selection.Find("span").Text(), " minute(s)")
				}

				// 片商
				if i == 4 {
					resp.Studio = selection.Find("span").Text()
				}

				// 系列
				if i == 5 {
					resp.Series = selection.Find("span").Text()
				}

				// 标签
				if i == 7 {
					var tags []string
					selection.Find("span a").Each(func(i int, selection *goquery.Selection) {
						tags = append(tags, utils.T2S(strings.TrimSpace(selection.Text())))
					})
					resp.Tags = tags
				}

				// actor
				if i == 8 {
					actors := make(map[string]string)
					selection.Find("span a").Each(func(i int, selection *goquery.Selection) {
						// 演员列表
						actors[strings.TrimSpace(selection.Text())] = ""
					})
					resp.Actors = actors
				}
			})

			resp.Intro = ""

			// 获取cover图片
			// 获取图片
			fanart, _ := root2.Find(`div.column-video-cover a img`).Attr("src")
			resp.Cover = fanart

			// 获取sample图片
			samples := []string{}
			root2.Find(".preview-images .tile-item").Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("href")
				samples = append(samples, href)
			})
			resp.SampleImg = samples

			if len(resp.Cover) == 0 {
				err = ErrorCoverNotFound
				continue
			}

			log.Println("✅javDb202504:", that.no, resp.WebSite)
		}

		// log.Println(htmlBody)

	}
	return resp, err
}
