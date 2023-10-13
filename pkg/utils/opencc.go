package utils

import (
	"github.com/liuzl/gocc"
)

func T2S(in string) string {
	s2t, _ := gocc.New("t2s")
	out, _ := s2t.Convert(in)
	return out
}
