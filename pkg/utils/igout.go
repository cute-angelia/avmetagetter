package utils

import (
	"github.com/guonaihong/gout"
	"time"
)

var igout *gout.Client

func init() {
	igout = gout.NewWithOpt(gout.WithTimeout(time.Second * 15))
}

func GetIGout() *gout.Client {
	return igout
}
