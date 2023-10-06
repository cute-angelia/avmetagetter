package sites

import "metagetter/pkg/utils"

type javbus struct {
	no        string
	useragent string
	cookies   string
	proxy     string
	siteuri   string
}

func NewJavBus(no, useragent, cookies, proxy string) *javbus {
	return &javbus{
		no:        no,
		useragent: useragent,
		cookies:   useragent,
		proxy:     proxy,
		siteuri:   "https://www.javbus.com/",
	}
}

func (that *javbus) Fetch() (SiteResp, error) {
	var resp SiteResp
	utils.GetIGout().GET("")
	return resp, nil
}
