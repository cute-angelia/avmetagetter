.PHONY: up

up:
	git add .
	git commit -am "update"
	git pull origin v2
	git push origin v2
	@echo "\n 代码提交发布..."

tag:
	git pull origin v2
	git add .
	git commit -am "update"
	git push origin v2
	git tag v1.1.1
	git push --tags
	@echo "\n tags 发布中..."

.PHONY: run
no := $(firstword $(MAKECMDGOALS))
# make run no=ADN-499
run:
	cd cmd/cli_jav && go run main.go $(no)
