package utils

import "strings"

// IntroFilter 简介信息过滤
func IntroFilter(intro string) string {
	// 替换<br>
	intro = strings.ReplaceAll(intro, "<br>", "\n")
	intro = strings.ReplaceAll(intro, "<br/>", "\n")
	intro = strings.ReplaceAll(intro, "<br />", "\n")
	// 替换\r\n
	intro = strings.ReplaceAll(intro, "\r\n", "\n")
	// 替换\r
	intro = strings.ReplaceAll(intro, "\r", "\n")
	// 替换\n\n
	intro = strings.ReplaceAll(intro, "\n\n", "\n")

	// 清除多余空白
	return strings.TrimSpace(intro)
}
