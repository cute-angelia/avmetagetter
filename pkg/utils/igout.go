package utils

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"strings"
	"time"
)

var igout *gout.Client

func init() {
	igout = gout.NewWithOpt(gout.WithTimeout(time.Second * 15))
}

func GetIGout(uri string, proxy string, debug bool) *dataflow.DataFlow {
	if len(proxy) > 0 {
		if strings.Contains(proxy, "socks5") {
			proxy = strings.Replace(proxy, "socks5://", "", -1)
			return igout.GET(uri).Debug(debug).SetSOCKS5(proxy)
		} else {
			return igout.GET(uri).Debug(debug).SetProxy(proxy)
		}
	}
	return igout.GET(uri).Debug(debug)
}
