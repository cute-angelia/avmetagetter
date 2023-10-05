package util

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/guonaihong/gout"
	"log"
	"strings"
	"time"
)

var gclent *gout.Client

func init() {
	gclent = gout.NewWithOpt(gout.WithTimeout(time.Second * 16))
}

type igout struct {
	useragent string
	cookies   string
	proxy     string
	client    *gout.Client
}

func NewIGout(useragent, cookies, proxy string) *igout {
	return &igout{
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
		client:    gclent,
	}
}

func (that *igout) Get(uri string) ([]byte, error) {
	var body []byte
	status := 0
	zgo := that.client.GET(uri)
	if len(that.proxy) > 0 {
		proxySocks5 := strings.Replace(that.proxy, "socks5://", "", -1)
		log.Println("proxySocks5=>", uri, proxySocks5)
		zgo.SetSOCKS5(proxySocks5)
	}
	err := zgo.SetHeader(gout.H{
		"cookie":     that.cookies,
		"user-agent": USER_AGENT,
	}).BindBody(&body).Code(&status).Do()
	// 检查错误
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (that *igout) Dom(body []byte) (*goquery.Document, error) {
	// 转换为节点数据
	root, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	// 检查错误
	if err != nil {
		return nil, err
	}
	return root, nil
}
