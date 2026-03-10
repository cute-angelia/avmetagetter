package utils

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

// GetBackgroundImage 提取指定选择器的 background-image URL
func GetBackgroundImage(s *goquery.Selection) string {
	// 获取 style 属性内容
	style, exists := s.Attr("style")
	if !exists {
		return ""
	}

	// 正则匹配 url(...) 内部内容，忽略单双引号及转义符
	// 匹配 url( 后面非 ) 的内容
	re := regexp.MustCompile(`url\(['"&quot;]*([^'"&quot;)]+)['"&quot;]*\)`)
	match := re.FindStringSubmatch(style)

	if len(match) > 1 {
		// 返回捕获组中的路径部分
		return match[1]
	}

	return ""
}
