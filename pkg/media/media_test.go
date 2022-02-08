package media

import (
	"github.com/cute-angelia/AVMeta/pkg/util"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"log"
	"testing"
)

func TestSearch(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	cfg := util.ConfigStruct{
		Base: util.BaseStruct{
			//Proxy:  proxyUri,
			//Socket: strings.Replace(proxyUri, "socks5://", "", -1),
		},
		Site: util.SiteStruct{
			JavBus: "https://www.javbus.com/",
			JavDB:  "https://javdb.com/",
		},
		Path: util.PathStruct{
			Filter: "-hd||hd-||_hd||hd_||[||]||【||】||asfur||~||-full||_full||3xplanet||monv||云中飘荡||@||tyhg999.com||xxxxxxxx||-fhd||_fhd||thz.la",
		},
	}

	nos := []string{
		//"PRED-372",
		//"101116-279",
		//"fc2-2653105",
		//"ABW-198",
		//"534IND-028",
		"467SHINKI-075",
	}

	for _, no := range nos {
		if m, err := Search(no, &cfg); err != nil {
			log.Println("抓取失败", err.Error())
			return
		} else {
			log.Println(ijson.Pretty(m))
		}
	}
}
