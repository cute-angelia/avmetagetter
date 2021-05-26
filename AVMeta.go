package main

import (
	"github.com/ylqjgm/AVMeta/pkg/cmd"
	"log"
)

var (
	version = "master"
	commit  = "?"
	built   = ""
)

func main() {
	// 设置 flag
	flag := log.LstdFlags | log.Lshortfile | log.Lmsgprefix
	log.SetFlags(flag)

	e := cmd.NewExecutor(version, commit, built)

	if err := e.Execute(); err != nil {
		panic(err)
	}
}
