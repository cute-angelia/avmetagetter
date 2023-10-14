package utils

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/guonaihong/gout"
	"github.com/spf13/viper"
	"log"
	"time"
)

func GetWs() string {
	addr := viper.GetString("chromedp.addr")
	var jsonconfig JsonConfig
	gout.GET(addr).BindJSON(&jsonconfig).Do()
	return jsonconfig.WebSocketDebuggerURL
}

type JsonConfig struct {
	Browser              string `json:"Browser"`
	ProtocolVersion      string `json:"Protocol-Version"`
	UserAgent            string `json:"User-Agent"`
	V8Version            string `json:"V8-Version"`
	WebKitVersion        string `json:"WebKit-Version"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func GetBody(uri string, waitVisible string) (string, error) {

	log.Println(uri, waitVisible)
	// getws
	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), GetWs())
	defer cancelActxt()

	ctxt, cancelCtxt := chromedp.NewContext(actxt, chromedp.WithBrowserOption(
		chromedp.WithDialTimeout(time.Second*10),
	)) // create new tab
	defer cancelCtxt() // close tab afterwards
	var title string
	var htmlContent string
	var urlLocation string
	if err := chromedp.Run(ctxt,
		chromedp.Navigate(uri),
		chromedp.WaitVisible(waitVisible),
		chromedp.Title(&title),
		chromedp.Location(&urlLocation),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	); err != nil {
		log.Printf("Failed: %v", err)
		return "", err
	}

	return htmlContent, nil
}
