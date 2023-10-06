package sites

import (
	"github.com/spf13/viper"
	"log"
	"metagetter/pkg/utils"
	"path"
)

type javbus struct {
	no        string
	useragent string
	cookies   string
	proxy     string
}

func NewJavBus(no, useragent, cookies, proxy string) *javbus {
	return &javbus{
		no:        no,
		useragent: useragent,
		cookies:   cookies,
		proxy:     proxy,
	}
}

func (that *javbus) Fetch() (SiteResp, error) {
	var resp SiteResp
	uri := path.Base(viper.GetString("javbus.site")) + that.no
	log.Println(uri)
	utils.GetIGout().GET(uri)
	return resp, nil
}
