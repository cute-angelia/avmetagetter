package sites

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/guonaihong/gout"
	"html"
	"log"
	"regexp"
	"strings"
)

type heyzo struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewHeyzo(no string, useragent, cookies, proxy string) *heyzo {
	return &heyzo{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "",
	}
}

func (that *heyzo) GetPageUri() []string {
	// 番号正则
	r := regexp.MustCompile(`[0-9]{4}`)
	number := r.FindString(that.no)
	return []string{
		fmt.Sprintf("https://www.heyzo.com/moviepages/%s/index.html", number),
	}
}

type jsonresp struct {
	Context  string `json:"@context"`
	Type     string `json:"@type"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Encoding struct {
		Type           string `json:"@type"`
		EncodingFormat string `json:"encodingFormat"`
	} `json:"encoding"`
	Actor struct {
		Type  string `json:"@type"`
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"actor"`
	Duration      string `json:"duration"`
	DateCreated   string `json:"dateCreated"`
	ReleasedEvent struct {
		Type      string `json:"@type"`
		StartDate string `json:"startDate"`
		Location  struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"location"`
	} `json:"releasedEvent"`
	Video struct {
		Type         string `json:"@type"`
		Description  string `json:"description"`
		Duration     string `json:"duration"`
		Name         string `json:"name"`
		Thumbnail    string `json:"thumbnail"`
		ThumbnailURL string `json:"thumbnailUrl"`
		UploadDate   string `json:"uploadDate"`
		Actor        string `json:"actor"`
		Provider     string `json:"provider"`
	} `json:"video"`
	AggregateRating struct {
		Type        string `json:"@type"`
		RatingValue string `json:"ratingValue"`
		BestRating  string `json:"bestRating"`
		ReviewCount string `json:"reviewCount"`
	} `json:"aggregateRating"`
}

func (that *heyzo) Fetch() (resp ScraperResp, err error) {
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

		// debug
		//log.Println(htmlBody)

		if root, err2 := goquery.NewDocumentFromReader(strings.NewReader(htmlBody)); err2 != nil {
			log.Println("ERROR:", err2)
			err = err2
			continue
		} else {
			// 获取json节点
			data, err3 := root.Find(`script[type="application/ld+json"]`).Html()
			if err3 != nil {
				err = err3
			}
			// 转码
			data = strings.ReplaceAll(html.UnescapeString(data), "\n", "")
			// 检查
			if err != nil {
				err = fmt.Errorf("%s [Find Json]: %s", uri, err)
				continue
			}

			// 转码
			// 转换为结构体
			js := jsonresp{}
			err = json.Unmarshal([]byte(data), &js)

			resp.No = that.no
			resp.WebSite = uri

			// 标题
			resp.Title = js.Name

			resp.Intro = js.Video.Description

			// 获取导演
			resp.Director = ""

			resp.ReleaseDate = js.DateCreated

			resp.Runtime = ""
			resp.Studio = js.Video.Provider
			resp.Series = js.Video.Provider

			// 类别数组
			var tags []string
			resp.Tags = tags

			// 获取sample图片
			sample := []string{}
			resp.SampleImg = sample

			// 获取cover图片
			if strings.Contains(js.Image, "http") {
				resp.Cover = js.Image
			} else {
				resp.Cover = "https:" + js.Image
			}

			// 演员数组
			actors := make(map[string]string)
			actors[js.Actor.Name] = "https:" + js.Actor.Image
			resp.Actors = actors

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
