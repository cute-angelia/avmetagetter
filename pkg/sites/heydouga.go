package sites

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/metagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"log"
	"regexp"
	"strings"
)

type heydouga struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string

	code1 string
	code2 string
}

func NewHeydouga(no string, useragent, cookies, proxy string) *heydouga {
	return &heydouga{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "",
	}
}

func (that *heydouga) GetPageUri() []string {

	// 转换大写
	code := strings.ToUpper(that.no)
	// 番号正则
	r := regexp.MustCompile(`([0-9]{4}).+?([0-9]{3,4})`)
	// 临时番号
	code = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(r.FindString(code), "PPV", ""), "HEYDOUGA", ""))
	// 检查是否为空
	if code == "" {
		return []string{}
	}

	// 番号分割
	cs := strings.Split(code, "-")
	// 检查是否有两个
	if len(cs) < 2 {
		return []string{}
	}

	// 设置番号前后缀
	code1 := cs[0]
	code2 := cs[1]

	that.code1 = code1
	that.code2 = code2

	// 组合地址
	uri := fmt.Sprintf("https://www.heydouga.com/moviepages/%s/%s/index.html", code1, code2)

	// 设置番号前后缀
	// 重新组合地址
	uri2 := fmt.Sprintf("https://www.heydouga.com/moviepages/%s/%s/index.html", cs[0], "ppv-"+cs[1])

	return []string{
		uri,
		uri2,
	}
}

func (that *heydouga) Fetch() (resp ScraperResp, err error) {
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
			// 查找是否获取到
			resp.No = that.no
			resp.WebSite = uri

			// 标题
			resp.Title = root.Find(`div#title-bg h1`).Text()
			if len(resp.Title) < 4 {
				err = errors.New("not found")
				continue
			}

			resp.Intro = utils.IntroFilter(root.Find(`div[class="movie-description"] p`).Text())

			// 获取导演
			resp.Director = root.Find(`span:contains("提供元")`).Next().Find(`a[href*="/listpages/provider"]`).Text()

			resp.ReleaseDate = root.Find(`span:contains("配信日")`).Next().Text()

			resp.Runtime = strings.ReplaceAll(root.Find(`span:contains("動画再生時間")`).Next().Text(), "分", "")
			resp.Studio = "Hey動画"
			resp.Series = "Hey動画 PPV"

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
			fanart := fmt.Sprintf("https://image01-www.heydouga.com/contents/%s/%s/player_thumb.jpg", that.code1, that.code2)
			resp.Cover = fanart

			// 演员数组
			actors := make(map[string]string)
			// 定义一个临时演员数组
			var tmpActors []string

			// 循环获取
			root.Find(`span:contains("主演")`).Next().Find(`a`).Each(func(i int, item *goquery.Selection) {
				// 获取演员信息
				act := strings.TrimSpace(item.Text())
				// 检查
				if act == "" {
					return
				}
				// 分割数据
				acts1 := strings.Split(act, "、")
				acts2 := strings.Split(act, " ")
				// 循环加入数组
				for _, a := range acts1 {
					tmpActors = append(tmpActors, strings.TrimSpace(strings.ReplaceAll(a, "素人", "")))
				}
				for _, a := range acts2 {
					tmpActors = append(tmpActors, strings.TrimSpace(strings.ReplaceAll(a, "素人", "")))
				}
			})

			// 循环加入map
			for _, actor := range tmpActors {
				actors[actor] = ""
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
