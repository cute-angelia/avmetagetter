package utils

import "strings"

func CleanNo(nostr string) string {
	nostr = strings.ToUpper(nostr)
	nostr = strings.ReplaceAll(nostr, "-UC", "")
	nostr = strings.ReplaceAll(nostr, "-C", "")
	return nostr
}
