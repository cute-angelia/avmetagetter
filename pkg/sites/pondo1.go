package sites

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cute-angelia/avmetagetter/pkg/utils"
	"github.com/guonaihong/gout"
)

type pondo1 struct {
	BuildInScraper
	no        string
	useragent string
	cookies   string
	proxy     string
	site      string
}

func NewPondo1(no string, useragent, cookies, proxy string) *pondo1 {
	return &pondo1{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		site:      "https://www.1pondo.tv",
	}
}

func (that *pondo1) GetPageUri() []string {
	return []string{
		fmt.Sprintf("%s/dyn/phpauto/movie_details/movie_id/%s.json", that.site, that.no),
		fmt.Sprintf("%s/dyn/dla/json/movie_gallery/%s.json", that.site, that.no),
	}
}

type MovieDetails struct {
	Actor         string   `json:"Actor"`
	ActorID       []int    `json:"ActorID"`
	ActressesJa   []string `json:"ActressesJa"`
	ActressesEn   []string `json:"ActressesEn"`
	ActressesList struct {
		Num4572 struct {
			NameJa string `json:"NameJa"`
			NameEn string `json:"NameEn"`
			Sizes  string `json:"Sizes"`
			Age    int    `json:"Age"`
		} `json:"4572"`
	} `json:"ActressesList"`
	AvgRating           float64     `json:"AvgRating"`
	CanStream           bool        `json:"CanStream"`
	Conditions          interface{} `json:"Conditions"`
	Desc                string      `json:"Desc"`
	DescEn              string      `json:"DescEn"`
	Duration            int         `json:"Duration"`
	Expire              interface{} `json:"Expire"`
	HasFlash            bool        `json:"HasFlash"`
	NoListDisplay       bool        `json:"NoListDisplay"`
	SampleExcludeFlag   bool        `json:"SampleExcludeFlag"`
	Gallery             bool        `json:"Gallery"`
	AffZip              bool        `json:"AffZip"`
	HasGallery          bool        `json:"HasGallery"`
	HasMemberGalleryZip bool        `json:"HasMemberGalleryZip"`
	HasSampleGalleryZip bool        `json:"HasSampleGalleryZip"`
	MetaMovieID         int         `json:"MetaMovieID"`
	MovieID             string      `json:"MovieID"`
	MovieSeq            interface{} `json:"MovieSeq"`
	MovieThumb          string      `json:"MovieThumb"`
	RealMetaMovieID     interface{} `json:"RealMetaMovieID"`
	Release             string      `json:"Release"`
	Series              interface{} `json:"Series"`
	SeriesEn            interface{} `json:"SeriesEn"`
	SeriesID            interface{} `json:"SeriesID"`
	SiteID              int         `json:"SiteID"`
	Status              bool        `json:"Status"`
	ThumbHigh           string      `json:"ThumbHigh"`
	ThumbLow            string      `json:"ThumbLow"`
	ThumbMed            string      `json:"ThumbMed"`
	ThumbUltra          string      `json:"ThumbUltra"`
	Title               string      `json:"Title"`
	TitleEn             string      `json:"TitleEn"`
	Type                int         `json:"Type"`
	Year                string      `json:"Year"`
	UC                  []int       `json:"UC"`
	UCNAME              []string    `json:"UCNAME"`
	UCNAMEEn            []string    `json:"UCNAMEEn"`
	UcNameList          struct {
		Num1 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"1"`
		Num13 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"13"`
		Num17 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"17"`
		Num25 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"25"`
		Num32 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"32"`
		Num38 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"38"`
		Num46 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"46"`
		Num48 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"48"`
		Num52 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"52"`
		Num54 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"54"`
		Num58 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"58"`
		Num61 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"61"`
		Num67 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"67"`
		Num69 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"69"`
		Num70 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"70"`
		Num60000 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"60000"`
		Num60001 struct {
			NameEn string `json:"NameEn"`
			NameJa string `json:"NameJa"`
		} `json:"60001"`
	} `json:"UcNameList"`
	IsTicketOnly bool `json:"IsTicketOnly"`
	MemberFiles  []struct {
		FileName    string `json:"FileName"`
		FileSize    int    `json:"FileSize"`
		MetaMovieID int    `json:"MetaMovieID"`
		SiteID      string `json:"SiteID"`
		URL         string `json:"URL"`
	} `json:"MemberFiles"`
	SampleFiles []struct {
		FileName    string `json:"FileName"`
		FileSize    int    `json:"FileSize"`
		MetaMovieID int    `json:"MetaMovieID"`
		SiteID      string `json:"SiteID"`
		URL         string `json:"URL"`
	} `json:"SampleFiles"`
	PPVPrice struct {
		Regular  int `json:"Regular"`
		Discount int `json:"Discount"`
		Campaign int `json:"Campaign"`
	} `json:"PPVPrice"`
}
type MovieGallery struct {
	MovieID string `json:"MovieID"`
	Rows    []struct {
		Img       string `json:"Img"`
		Protected bool   `json:"Protected"`
	} `json:"Rows"`
}

func (that *pondo1) Fetch() (resp ScraperResp, err error) {
	uris := that.GetPageUri()

	resp.No = that.no
	resp.WebSite = fmt.Sprintf("https://www.1pondo.tv/movies/%s", that.no)

	for _, uri := range uris {
		if strings.Contains(uri, "movie_details") {
			var respDetail MovieDetails
			utils.GetIGout(uri, that.proxy, false).SetHeader(gout.H{
				"User-Agent": that.useragent,
				"Cookie":     that.cookies,
				"referer":    that.site,
			}).BindJSON(&respDetail).Do()

			resp.Title = respDetail.Title

			if len(resp.Title) < 10 {
				err = errors.New("title not right")
				break
			}

			// 获取简介
			resp.Intro = respDetail.Desc

			resp.Director = ""

			resp.ReleaseDate = respDetail.Release
			resp.Runtime = fmt.Sprintf("%d", respDetail.Duration)

			resp.Studio = "1pondo"
			resp.Series = ""

			resp.Tags = respDetail.UCNAME

			// 获取cover图片
			resp.Cover = respDetail.ThumbUltra

			if len(resp.Cover) == 0 {
				err = ErrorCoverNotFound
				break
			}

			// 演员数组
			actors := make(map[string]string, len(respDetail.ActressesJa))
			for _, s := range respDetail.ActressesJa {
				actors[s] = ""
			}
			resp.Actors = actors

		}

		if strings.Contains(uri, "movie_gallery") {
			var respDetail MovieGallery
			utils.GetIGout(uri, that.proxy, false).SetHeader(gout.H{
				"User-Agent": that.useragent,
				"Cookie":     that.cookies,
				"referer":    that.site,
			}).BindJSON(&respDetail).Do()

			for _, row := range respDetail.Rows {
				// https://www.1pondo.tv/dyn/dla/images/movie_gallery/sample/022426_001/035e3ec5a0c6d5be81fed55848df14a6.jpg
				// https://www.1pondo.tv/dyn/dla/images/movie_gallery/sample/022426_001/0d58fac4d10a327661e19b2b955b1468.jpg
				// movie_gallery/member/022426_001/aefba6b9eb189be45fc21acc6ba030be.jpg

				img := row.Img
				newimg := ""
				ioms := strings.Split(img, "/")
				if len(ioms) > 1 {
					newimg = ioms[len(ioms)-1]
				}

				resp.SampleImg = append(resp.SampleImg, that.site+"//dyn/dla/images/movie_gallery/sample/"+that.no+"/"+newimg)
			}

			// 只获取5个,5个后面需要会员。。
			maxlen := len(resp.SampleImg)
			if len(resp.SampleImg) >= 5 {
				maxlen = 5
			}
			resp.SampleImg = resp.SampleImg[0:maxlen]
		}
	}
	return resp, err
}
