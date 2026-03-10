package media

import "github.com/cute-angelia/avmetagetter/pkg/sites"

// media base
type (
	Media interface {
		ParseMedia(resp sites.ScraperResp)
		Marshal() ([]byte, error)
	}
	BuildMedia struct{}
)

func (_ *BuildMedia) Marshal() (_ []byte, _ error) {
	return
}

// xml base
type (
	// Inner 文字数据，为了避免某些内容被转义。
	Inner struct {
		Inner string `xml:",innerxml"`
	}

	// Actor 演员信息，保存演员姓名及头像地址。
	Actor struct {
		Name  string `xml:"name"`
		Thumb string `xml:"thumb"`
	}

	Art struct {
		Poster string `xml:"poster"`
		Fanart string `xml:"fanart"`
	}

	JellyfinMeta struct {
		Tmdbid int    `xml:"tmdbid"` // tmdbid
		Key    string `xml:"key"`    // cache key
		Update string `xml:"update"`
	}
)
