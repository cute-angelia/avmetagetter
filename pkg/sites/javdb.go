package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/go-utils/syntax/iurl"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"log"
	"strings"
)

type javDb struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewJavDb(no string, useragent, cookies, proxy string) *javDb {
	return &javDb{
		no:        strings.ToUpper(no),
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://javdb.com/",
	}
}

// GetPageUri 获取页面地址
func (that *javDb) GetPageUri() string {
	return fmt.Sprintf("%s/search?q=%s&f=all", iurl.GetDomainWithOutSlant(that.site), that.no)
}

func (that *javDb) Fetch() (resp ScraperResp, err error) {
	uri := that.GetPageUri()
	if !strings.Contains(uri, "http") {
		return resp, errors.New("error url address:" + uri)
	}
	var htmlBody string
	// get
	//code := 0
	//utils.GetIGout(uri, that.proxy, true).SetHeader(gout.H{
	//	"User-Agent": that.useragent,
	//	"Referer":    that.site,
	//}).BindBody(&htmlBody).Code(&code).Do()

	//if code == 0 || code == 403 {
	if htmlBody, err = utils.GetBody(uri, "body > section > div > div.movie-list.h.cols-4.vcols-8"); err != nil {
		return resp, err
	}
	//}

	if rootPre, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err != nil {
		log.Println("ERROR:", err)
		return resp, err
	} else {
		// 特殊处理
		newUrl := ""
		rootPre.Find(".movie-list").Each(func(i int, selection *goquery.Selection) {
			selection.Find(".item").Each(func(i int, selection2 *goquery.Selection) {
				findCode := selection2.Find("strong").First().Text()
				if that.no == findCode {
					newUrl, _ = selection2.Find("a").First().Attr("href")
				}
			})
		})
		if htmlBody, err = utils.GetBody(iurl.GetDomainWithOutSlant(that.site)+newUrl, "body > section > div"); err != nil {
			return resp, err
		}

		if root, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err != nil {
			log.Println("ERROR:", err)
			return resp, err
		} else {
			//log.Println(htmlBody)
			// 查找是否获取到
			if -1 < root.Find(`.empty-message:contains("暫無內容")`).Index() {
				return resp, errors.New("404 Not Found")
			}

			// 查找是否获取到
			resp.No = that.no
			resp.WebSite = uri

			t := ""
			root.Find(`h2.title strong`).Each(func(i int, selection *goquery.Selection) {
				t += selection.Text() + " "
			})
			resp.Title = strings.TrimSpace(t)
			resp.Intro = ""

			// 获取数据
			val := root.Find(`strong:contains("導演")`).NextFiltered(`span.value`).Text()
			// 检查
			if val == "" {
				val = root.Find(`strong:contains("導演")`).NextFiltered(`span.value`).Find("a").Text()
			}
			resp.Director = val

			// 获取数据
			val = root.Find(`strong:contains("日期")`).NextFiltered(`span.value`).Text()
			// 检查
			if val == "" {
				val = root.Find(`strong:contains("日期")`).NextFiltered(`span.value`).Find("a").Text()
			}
			resp.ReleaseDate = val

			// 获取数据
			val = root.Find(`strong:contains("時長")`).NextFiltered(`span.value`).Text()
			// 检查
			if val == "" {
				val = root.Find(`strong:contains("時長")`).NextFiltered(`span.value`).Find("a").Text()
			}
			// 去除多余
			val = strings.TrimRight(val, "分鍾")
			resp.Runtime = val

			// 获取数据
			val = root.Find(`strong:contains("片商")`).NextFiltered(`span.value`).Text()
			// 检查
			if val == "" {
				val = root.Find(`strong:contains("片商")`).NextFiltered(`span.value`).Find("a").Text()
			}
			resp.Studio = val

			// 获取数据
			val = root.Find(`strong:contains("系列")`).NextFiltered(`span.value`).Text()
			// 检查
			if val == "" {
				val = root.Find(`strong:contains("系列")`).NextFiltered(`span.value`).Find("a").Text()
			}
			resp.Series = val

			// 类别数组
			var tags []string
			// 循环获取
			root.Find(`strong:contains("類別")`).NextFiltered(`span.value`).Find("a").Each(func(i int, item *goquery.Selection) {
				tags = append(tags, utils.T2S(strings.TrimSpace(item.Text())))
			})
			resp.Tags = tags

			// 获取cover图片
			// 获取图片
			fanart, _ := root.Find(`div.column-video-cover a img`).Attr("src")
			resp.Cover = fanart

			// 获取sample图片
			samples := []string{}
			root.Find(".preview-images .tile-item").Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("href")
				samples = append(samples, href)
			})
			resp.SampleImg = samples

			// 演员数组
			// 演员列表
			actors := make(map[string]string)

			// 循环获取
			root.Find(`strong:contains("演員")`).NextFiltered(`span.value`).Find("a").Each(func(i int, item *goquery.Selection) {
				// 演员名称
				actors[strings.TrimSpace(item.Text())] = ""
			})
			resp.Actors = actors
		}
	}

	// log.Println(htmlBody)

	return resp, nil
}
