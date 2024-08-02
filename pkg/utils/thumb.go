package utils

import (
	"github.com/cute-angelia/go-utils/utils/iimage"
)

// MakeThumbCover 图片进行裁剪
func MakeThumbCover(imgIn string, imgOut string) error {
	return iimage.CropJavCover(imgIn, imgOut, 21)
}
