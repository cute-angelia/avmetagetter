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

type caribBeanCom struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewCaribBeanCom(no string, useragent, cookies, proxy string) *caribBeanCom {
	return &caribBeanCom{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://www.caribbeancom.com",
	}
}

func (that *caribBeanCom) GetPageUri() []string {
	return []string{fmt.Sprintf("%s/moviepages/%s/index.html", that.site, that.no)}
}

func (that *caribBeanCom) Fetch() (resp ScraperResp, err error) {
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
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {
			// 查找是否获取到
			if -1 == root.Find(`h3`).Index() {
				err = errors.New("404 Not Found")
				continue
			}

			resp.No = that.no
			resp.WebSite = uri
			resp.Title = root.Find(`h1[itemprop="name"]`).Text()

			// 获取简介
			intro, err3 := root.Find(`p[itemprop="description"]`).Html()
			// 检查
			if err3 != nil {
				intro = ""
			} else {
				intro = utils.IntroFilter(intro)
			}

			resp.Intro = intro
			resp.Director = ""
			resp.ReleaseDate = root.Find(`span[itemprop="uploadDate"]`).Text()
			resp.Runtime = strings.TrimSpace(root.Find(`span[itemprop="duration"]`).Text())

			resp.Studio = "カリビアンコム"
			resp.Series = root.Find(`a[href*="/series/"]`).Text()

			// 类别数组
			var tags []string
			// 循环获取
			// 日语 tag 翻译不精确，算了
			//root.Find(`a[itemprop="genre"]`).Each(func(i int, item *goquery.Selection) {
			//	tags = append(tags, strings.TrimSpace(item.Text()))
			//})
			resp.Tags = tags

			// 获取cover图片
			resp.Cover = fmt.Sprintf("https://www.caribbeancom.com/moviepages/%s/images/l_l.jpg", that.no)

			// 获取sample图片
			// https://www.caribbeancom.com/moviepages/010113-225/images/g_big012.jpg
			// https://c0.jdbstatic.com/samples/ev/EvEd3_l_4.jpg
			// https://www.caribbeancom.com/moviepages/102715-008/images/s/015.jpg
			// https://www.caribbeancom.com/moviepages/102715-008/images/b/015.jpg
			// https://www.caribbeancom.com/moviepages/102715-008/images/g_big015.jpg
			// https://www.caribbeancom.com/moviepages/102715-008/images/g_s/015.jpg
			sample := []string{}
			root.Find(`div.gallery .grid-item`).Each(func(i int, selection *goquery.Selection) {
				if v, ok := selection.Find("a").First().Attr("href"); ok {
					v = that.site + v
					v = strings.Replace(v, "/member", "", 1)
					sample = append(sample, v)
				}
			})
			resp.SampleImg = sample

			// 演员数组
			actors := make(map[string]string)
			// 循环获取
			root.Find(`a[class="spec__tag"] span[itemprop="name"]`).Each(func(i int, item *goquery.Selection) {
				// 演员名称
				actors[strings.TrimSpace(item.Text())] = ""
			})
			resp.Actors = actors

			if len(resp.Title) < 10 {
				err = errors.New("title not right")
				continue
			}
		}
	}
	// log.Println(htmlBody)

	return resp, err
}
