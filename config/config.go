package config

import (
	_ "embed"
	"github.com/cute-angelia/go-utils/utils/conf"
)

//go:embed config.local.toml
var configLocal []byte

//go:embed config.product.toml
var configProduct []byte

func InitConfig(envStr string) {
	switch envStr {
	case "local":
		conf.MustLoadConfigByte(configLocal, "toml")
	default:
		conf.MustLoadConfigByte(configProduct, "toml")
	}
}
