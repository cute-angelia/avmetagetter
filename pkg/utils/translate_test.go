package utils

import (
	"github.com/xuthus5/BaiduTranslate"
	"log"
	"testing"
)

func TestTranlate(t *testing.T) {
	bi := BaiduTranslate.BaiduInfo{AppID: "x", Salt: BaiduTranslate.Salt(5), SecretKey: "x", From: "auto", To: "zh"}
	ss := []string{
		"初裏",
		"スレンダー",
		"ぶっかけ",
	}
	bi.From = "jp"
	bi.To = "zh"

	for _, s := range ss {
		bi.Text = s
		log.Println(bi.Translate())
	}
}
