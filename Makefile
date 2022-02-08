.DEFAULT_GOAL = test
.PHONY: FORCE

SHELL := /bin/bash
BASEDIR = $(shell pwd)

version = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else echo master; fi)
commit = $(shell git rev-parse --short HEAD)
built = $(shell TZ=UTC date +%FT%T%z)
ldflags="-s -w -X main.version=${version} -X main.commit=${commit} -X main.built=${built}"

# enable module support across all go commands.
export GO111MODULE = on
# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://goproxy.io

# Build

.PHONY: build
build:
	go build
	AVMeta

build_race:
	go build -race -ldflags ${ldflags}
.PHONY: build_race

clean:
	rm -f AVMeta
.PHONY: clean

# Test
test: build
	go test -v ./...
.PHONY: test

AVMeta: FORCE
	go build -ldflags ${ldflags}

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

unexport GOFLAGS
vendor_free_build: FORCE
	go build -ldflags ${ldflags}

.PHONY: up tag
up:
	git add .
	git commit -am "update"
	git pull origin master
	git push origin master
	@echo "\n 代码提交发布..."


tag:
	git pull origin master
	git add .
	git commit -am "update"
	git push origin master
	git tag v1.0.8
	git push --tags
	@echo "\n tags 发布中..."

