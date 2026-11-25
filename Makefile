# 设置默认目标
.DEFAULT_GOAL := setup

# 环境变量优先级高于 .env：include 的赋值会覆盖环境值，先暂存一份
ifeq ($(origin GIT_REMOTE),environment)
    GIT_REMOTE_ENV := $(GIT_REMOTE)
endif

# 加载 .env 文件（如果存在）
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifneq ($(GIT_REMOTE_ENV),)
    GIT_REMOTE := $(GIT_REMOTE_ENV)
endif

# 同步/推送目标远程，逗号分隔：GIT_REMOTE=origin,github,gitlab
# 可通过 .env / 环境变量 / 命令行覆盖，未配置时默认 origin。
GIT_REMOTE ?= origin
export GIT_REMOTE

#####################
## 变量定义
#####################

# 二进制文件名
BINARY_NAME := dtskills

# 构建目录
BUILD_DIR := dist

# 工作目录（配置文件和统计数据都存储在这里）
DTSKILLS_WORK_DIR = .

# 一键设置项目（默认执行）
setup: deps
	@echo "🎉 Project setup completed!"

# 安装全部依赖
.PHONY: deps
deps: deps-go
	@echo "[Deps] 完成。"

# 安装全部 Go 依赖
.PHONY: deps-go
deps-go:
	@echo "[Deps] 安装 Go 模块依赖..."
	go mod tidy
	go mod download

# 更新全部 Go 依赖
.PHONY: update-deps-go
update-deps-go:
	@echo "[Deps] 更新 Go 模块依赖..."
	go get -u ./...
	go mod tidy

#####################
## 开发工具
#####################

# 生成技能包
# make generate
.PHONY: generate
generate:
	@echo "[生成] 生成技能包..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . generate

#####################

# 获取 配置文件路径
# make config-list
.PHONY: config-list
config-list:
	@echo "[配置] 获取配置文件列表..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . config list

# 获取 配置文件路径
# make config-path
.PHONY: config-path
config-path:
	@echo "[配置] 获取配置文件路径..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . config path

# 交互删除 配置文件项
# make config-remove
.PHONY: config-remove
config-remove:
	@echo "[配置] 交互删除配置文件项..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . config remove

#####################

# 获取 统计数据
# make stats-list
.PHONY: stats-list
stats-list:
	@echo "[统计] 获取统计数据..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . stats list

# 获取 统计数据路径
# make stats-path
.PHONY: stats-path
stats-path:
	@echo "[统计] 获取统计数据路径..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . stats path

#####################

# 获取 百度搜索
# make baidu-search QUERY="人工智能" COUNT=10
# make baidu-search QUERY="人工智能" COUNT=10 MINIFY="--minify"
.PHONY: baidu-search
baidu-search:
	@echo "[获取] 百度搜索..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-search search --query "$(QUERY)" --count $(COUNT) $(MINIFY)

#####################

# 获取 百度百科
# make baidu-baike KEYWORD="北京"
# make baidu-baike KEYWORD="北京" MINIFY="--minify"
.PHONY: baidu-baike
baidu-baike:
	@echo "[获取] 百度百科..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-baike search --key "$(KEYWORD)" $(MINIFY)

#####################

# 获取 腾讯新闻（搜索）
# make tencent-news-search QUERY="人工智能"
# make tencent-news-search QUERY="人工智能" MINIFY="--minify"
.PHONY: tencent-news-search
tencent-news-search:
	@echo "[获取] 腾讯新闻（搜索）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news search --query "$(QUERY)" $(MINIFY)

# 获取 腾讯新闻（热点）
# make tencent-news-hot
# make tencent-news-hot MINIFY="--minify"
.PHONY: tencent-news-hot
tencent-news-hot:
	@echo "[获取] 腾讯新闻（热点）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news hot $(MINIFY)

# 获取 腾讯新闻（早报）
# make tencent-news-morning
# make tencent-news-morning MINIFY="--minify"
.PHONY: tencent-news-morning
tencent-news-morning:
	@echo "[获取] 腾讯新闻（早报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news morning $(MINIFY)

# 获取 腾讯新闻（晚报）
# make tencent-news-evening
# make tencent-news-evening MINIFY="--minify"
.PHONY: tencent-news-evening
tencent-news-evening:
	@echo "[获取] 腾讯新闻（晚报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news evening $(MINIFY)

#####################

# 获取 腾讯天气（当前）
# make tencent-weather-current ADCODE=110000
# make tencent-weather-current ADCODE=110000 MINIFY="--minify"
.PHONY: tencent-weather-current
tencent-weather-current:
	@echo "[获取] 腾讯天气（当前）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather current --adcode "$(ADCODE)" $(MINIFY)

# 获取 腾讯天气（预报）
# make tencent-weather-forecast ADCODE=110000 TYPE=forecast_24h
# make tencent-weather-forecast ADCODE=110000 TYPE=forecast_24h MINIFY="--minify"
.PHONY: tencent-weather-forecast
tencent-weather-forecast:
	@echo "[获取] 腾讯天气（预报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather forecast --adcode "$(ADCODE)" --type "$(TYPE)" $(MINIFY)

# 获取 腾讯天气（空气质量）
# make tencent-weather-air-quality ADCODE=110000 
# make tencent-weather-air-quality ADCODE=110000 MINIFY="--minify"
.PHONY: tencent-weather-airquality
tencent-weather-air-quality:
	@echo "[获取] 腾讯天气（空气质量）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather air-quality --adcode "$(ADCODE)" $(MINIFY)

# 获取 腾讯天气（IP定位获取城市）
# make tencent-weather-ip2city
# make tencent-weather-ip2city MINIFY="--minify"
.PHONY: tencent-weather-ip2city
tencent-weather-ip2city:
	@echo "[获取] 腾讯天气（IP定位获取城市）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather ip2city $(MINIFY)

#####################

# 腾讯较真事实查证
# make jiaozhen-factcheck-check QUERY="喝柠檬水能抗癌是真的吗"
# make jiaozhen-factcheck-check QUERY="喝柠檬水能抗癌是真的吗" MINIFY="--minify"
.PHONY: jiaozhen-factcheck-check
jiaozhen-factcheck-check:
	@echo "[查证] 腾讯较真事实查证..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . jiaozhen-factcheck check --query "$(QUERY)" $(MINIFY)

#####################

# 微博搜索 获取微博热搜榜
# make weibo-search-hot-search CATEGORY="主榜"
# make weibo-search-hot-search CATEGORY="主榜" MINIFY="--minify"
.PHONY: weibo-search-hot-search
weibo-search-hot-search:
	@echo "[微博搜索] 获取微博热搜榜..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . weibo-search hot-search --category "$(CATEGORY)" $(MINIFY)

# 微博搜索 执行微博智搜
# make weibo-search-smart-search QUERY="小米"
# make weibo-search-smart-search QUERY="小米" MINIFY="--minify"
.PHONY: weibo-search-smart-search
weibo-search-smart-search:
	@echo "[微博搜索] 执行微博智搜..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . weibo-search smart-search --query "$(QUERY)" $(MINIFY)

#####################

# 和风天气 获取当前天气
# make hefeng-weather-current LOCATION=101280601
# make hefeng-weather-current LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-current
hefeng-weather-current:
	@echo "[和风天气] 获取当前天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather current --location "$(LOCATION)" $(MINIFY)

# 和风天气 获取每日天气预报
# make hefeng-weather-forecast LOCATION=101280601
# make hefeng-weather-forecast LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-forecast
hefeng-weather-forecast:
	@echo "[和风天气] 获取每日天气预报..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather forecast --location "$(LOCATION)" $(MINIFY)

# 和风天气 获取天气指数
# make hefeng-weather-indices LOCATION=101280601
# make hefeng-weather-indices LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-indices
hefeng-weather-indices:
	@echo "[和风天气] 获取天气指数..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather indices --location "$(LOCATION)" $(MINIFY)

# 和风天气 获取日出日落
# make hefeng-weather-sunrise-sunset LOCATION=101280601
# make hefeng-weather-sunrise-sunset LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-sunrise-sunset
hefeng-weather-sunrise-sunset:
	@echo "[和风天气] 获取日出日落..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather sunrise-sunset --location "$(LOCATION)" $(MINIFY)

# 和风天气 获取月升月落和月相
# make hefeng-weather-moon-phase LOCATION=101280601
# make hefeng-weather-moon-phase LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-moon-phase
hefeng-weather-moon-phase:
	@echo "[和风天气] 获取月升月落和月相..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather moon-phase --location "$(LOCATION)" $(MINIFY)

#####################

# 百度地图 获取天气
# make baidu-map-weather DISTRICT_ID=440300
# make baidu-map-weather DISTRICT_ID=440300 MINIFY="--minify"
.PHONY: baidu-map-weather
baidu-map-weather:
	@echo "[百度地图] 获取天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-map weather --district_id "$(DISTRICT_ID)" $(MINIFY)

# 百度地图 获取IP定位
# make baidu-map-ip-location
# make baidu-map-ip-location MINIFY="--minify"
.PHONY: baidu-map-ip-location
baidu-map-ip-location:
	@echo "[百度地图] 获取IP定位..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-map ip-location $(MINIFY)

#####################

# 高德地图 获取天气
# make amap-map-weather CITY=440300
# make amap-map-weather CITY=440300 MINIFY="--minify"
.PHONY: amap-map-weather
amap-map-weather:
	@echo "[高德地图] 获取天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . amap-map weather --city "$(CITY)" $(MINIFY)

# 高德地图 获取IP定位
# make amap-map-ip-location IP=114.247.50.2
# make amap-map-ip-location IP=114.247.50.2 MINIFY="--minify"
.PHONY: amap-map-ip-location
amap-map-ip-location:
	@echo "[高德地图] 获取IP定位..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . amap-map ip-location --ip "$(IP)" $(MINIFY)

#####################

# 腾讯地图 获取当前天气
# make tencent-map-current-weather ADCODE=440300
.PHONY: tencent-map-current-weather
tencent-map-current-weather:
	@echo "[腾讯地图] 获取当前天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-map current-weather --adcode "$(ADCODE)"

# 腾讯地图 获取天气预报
# make tencent-map-forecast-weather ADCODE=440300
.PHONY: tencent-map-forecast-weather
tencent-map-forecast-weather:
	@echo "[腾讯地图] 获取天气预报..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-map forecast-weather --adcode "$(ADCODE)"

# 腾讯地图 获取IP定位
# make tencent-map-ip-location
.PHONY: tencent-map-ip-location
tencent-map-ip-location:
	@echo "[腾讯地图] 获取IP定位..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-map ip-location

#####################

# 中国气象局天气 获取实时天气
# make cma-weather-current-now CITY=59493
# make cma-weather-current-now CITY=59493 MINIFY="--minify"
.PHONY: cma-weather-current-now
cma-weather-current-now:
	@echo "[中国气象局天气] 获取实时天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . cma-weather now --city "$(CITY)" $(MINIFY)

#####################

# SoJSON 获取实时天气
# make sojson-weather CITY=101030100
# make sojson-weather CITY=101030100 MINIFY="--minify"
.PHONY: sojson-weather
sojson-weather:
	@echo "[SoJSON] 获取实时天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . sojson weather --city "$(CITY)" $(MINIFY)

#####################

# UAPI 获取实时天气
# make uapi-weather ADCODE=440300
# make uapi-weather ADCODE=440300 MINIFY="--minify"
.PHONY: uapi-weather
uapi-weather:
	@echo "[UAPI] 获取实时天气..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . uapi weather --adcode "$(ADCODE)" $(MINIFY)

#####################

# 爬虫 获取黄金价格全部数据
# make crawler-gold-price-data DATE=2026-05-22
# make crawler-gold-price-data DATE=2026-05-22 MINIFY="--minify"
.PHONY: crawler-gold-price-data
crawler-gold-price-data:
	@echo "[爬虫] 获取黄金价格全部数据..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . crawler gold-price-data --date "$(DATE)" $(MINIFY)

# 爬虫 获取最新黄金价格
# make crawler-latest-gold-price DATE=2026-05-22
# make crawler-latest-gold-price DATE=2026-05-22 MINIFY="--minify"
.PHONY: crawler-latest-gold-price
crawler-latest-gold-price:
	@echo "[爬虫] 获取最新黄金价格..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . crawler latest-gold-price --date "$(DATE)" $(MINIFY)

# 爬虫 获取京东金融黄金市场状态
# make crawler-jd-finance-gold-status
# make crawler-jd-finance-gold-status FORMAT="--format is_open"
# make crawler-jd-finance-gold-status MINIFY="--minify"
.PHONY: crawler-jd-finance-gold-status
crawler-jd-finance-gold-status:
	@echo "[爬虫] 获取京东金融黄金市场状态..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . crawler jd-finance-gold-status $(FORMAT) $(MINIFY)

# 爬虫 获取热搜数据
# make crawler-hot-data DATE=2026-05-22
# make crawler-hot-data DATE=2026-05-22 MINIFY="--minify"
.PHONY: crawler-hot-data
crawler-hot-data:
	@echo "[爬虫] 获取热搜数据..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . crawler hot-data --date "$(DATE)" $(MINIFY)

# 爬虫 获取天气数据
# make crawler-weather-data CITY=广州 DATE=2026-05-22
# make crawler-weather-data CITY=广州 DATE=2026-05-22 MINIFY="--minify"
.PHONY: crawler-weather-data
crawler-weather-data:
	@echo "[爬虫] 获取天气数据..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . crawler weather-data --city "$(CITY)" --date "$(DATE)" $(MINIFY)

#####################
## 构建 / 打包 相关
#####################

# GoReleaser 版本（流水线与本地共用，升级只改这里）
GORELEASER_VERSION ?= v2.20.2

# goreleaser 可执行文件（不在 PATH 时回退到 $HOME/go/bin）
GORELEASER ?= $(shell command -v goreleaser 2>/dev/null || echo $(HOME)/go/bin/goreleaser)

# goreleaser 构建时注入的版本名：流水线在 tag 创建前构建，用 GORELEASER_CURRENT_TAG 覆盖当前 tag；
# 本仓库的 tag 可能被 git describe 识别成 refs/tags/...，这里统一剥掉前缀。
GORELEASER_ENV = GORELEASER_CURRENT_TAG=$(patsubst refs/tags/%,%,$(VERSION))

# 同时编译/链接的目标数。默认按 CPU 数并发，CI 容器内存有限时容易在链接阶段被
# OOM Killer 杀掉（`link: signal: killed`）。本地默认 1，CI 环境（CI=true）默认 2，
# 两者都可显式覆盖：make goreleaser-snapshot GORELEASER_PARALLELISM=4
GORELEASER_PARALLELISM ?= $(if $(filter true,$(CI)),2,1)

# 非 CI（本地）快照只编当前平台：`--single-target` 是 `goreleaser build` 的参数，
# `goreleaser release` 没有（传了会直接报 unknown flag），所以本地走 build，只产出
# 当前平台的三个二进制，不生成压缩包 / Linux 系统包；CI 走 release 出全平台产物。
# 本地要全平台产物（含系统包）：CI=true make goreleaser-snapshot
GORELEASER_SNAPSHOT_CMD ?= $(if $(filter true,$(CI)),release --snapshot,build --snapshot --single-target)

# 安装固定版本的 goreleaser（已存在则跳过）
.PHONY: goreleaser-install
goreleaser-install:
	@if ! $(GORELEASER) --version >/dev/null 2>&1; then \
		echo "安装 goreleaser $(GORELEASER_VERSION) → $(dir $(GORELEASER))"; \
		GOBIN=$(dir $(GORELEASER)) go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION); \
	fi

# 校验 .goreleaser.yaml
.PHONY: goreleaser-check
goreleaser-check: goreleaser-install
	$(GORELEASER) check

# 快照构建（无需 tag）：产物在 dist/。本地非 CI 只编当前平台。
.PHONY: goreleaser-snapshot
goreleaser-snapshot: goreleaser-install
	$(GORELEASER) $(GORELEASER_SNAPSHOT_CMD) --clean --parallelism $(GORELEASER_PARALLELISM)

# 发布构建：产物在 dist/。
# GORELEASER_CURRENT_TAG 只提供版本名，validate 仍要求 tag 指向当前提交，
# 因此必须 --skip=validate，否则会报 "git tag ... was not made against commit"。
.PHONY: goreleaser-release
goreleaser-release: goreleaser-install
	$(GORELEASER_ENV) $(GORELEASER) release --clean --skip=validate --parallelism $(GORELEASER_PARALLELISM)

#####################
## 格式化 相关
#####################

# 格式化后端文件
.PHONY: format-go
format-go:
	@echo "[Format] 格式化后端文件 (GO)..."
	gofmt -w -s .
	go fmt ./...
	go fix ./...

# 格式化 YAML 文件
.PHONY: format-yaml
format-yaml:
	@bash scripts/format.sh yaml

# 格式化严格 JSON 文件（gojq）
.PHONY: format-json
format-json:
	@bash scripts/format.sh json

# 格式化 JSON5 文件（jsonfmt，保留注释 / 无引号 key / 尾逗号）
.PHONY: format-json5
format-json5:
	@bash scripts/format.sh json5

# 格式化 SQL 文件
.PHONY: format-sql
format-sql:
	@bash scripts/format.sh sql

# 格式化 Markdown 文件
.PHONY: format-markdown
format-markdown:
	@bash scripts/format.sh markdown

# 格式化 Dockerfile 文件
.PHONY: format-docker
format-docker:
	@bash scripts/format.sh docker

# 格式化 Shell 脚本
.PHONY: format-shell
format-shell:
	@bash scripts/format.sh shell

# 格式化/修复项目
.PHONY: format
format: fetch-scripts format-go format-yaml format-json format-json5 format-sql format-markdown format-docker format-shell

#####################
## 检查 / 测试 相关
#####################

# 统一定义构建命令
# CGO_ENABLED=0: 禁用 CGO，实现完全静态链接
# CGO_ENABLED=1: 启用 CGO，实现动态链接
LINT_ENV = CGO_ENABLED=0

# 定义 Lint 过程中需要激活的 Go Build Tags
LINT_TAGS = goolm,stdjson

# Lint 扫描的目标平台，默认当前平台；CI 通过 LINT_PLATFORMS=... 覆盖为全平台
LINT_PLATFORMS ?= $(shell go env GOOS)/$(shell go env GOARCH)

# 运行 Go 静态检查
# -printf=false 关闭 printf 分析，避免 i18n.T() 作日志 format 时触发错误
.PHONY: check-go-vet
check-go-vet:
	@echo "[Check] -> 运行 Go vet 静态检查 (Tags: $(LINT_TAGS))..."
	$(LINT_ENV) go vet -tags=$(LINT_TAGS) -printf=false ./...

# 运行 Go lint 扫描
.PHONY: check-go-lint
check-go-lint:
	@echo "[Check] -> 运行 Go lint 扫描 (平台: $(LINT_PLATFORMS), Tags: $(LINT_TAGS))..."
	@for os_arch in $(LINT_PLATFORMS); do \
		os=$$(echo $$os_arch | cut -d'/' -f1); \
		arch=$$(echo $$os_arch | cut -d'/' -f2); \
		echo "  -> $$os 平台 ($$arch)"; \
		$(LINT_ENV) GOOS=$$os GOARCH=$$arch golangci-lint fmt; \
		$(LINT_ENV) GOOS=$$os GOARCH=$$arch golangci-lint run --build-tags $(LINT_TAGS); \
	done

# 运行 Go fix 问题
.PHONY: check-go-fix
check-go-fix:
	@echo "[Check] -> 运行 Go fix 问题 (Tags: $(LINT_TAGS))..."
	$(LINT_ENV) golangci-lint run --fix --build-tags $(LINT_TAGS)

# 运行 Go 单元测试
.PHONY: check-go-test
check-go-test:
	@echo "[Check] -> 运行 Go 单元测试 (Tags: $(LINT_TAGS))..."
	$(LINT_ENV) go test -tags=$(LINT_TAGS) $(GO_COVER_TEST_FLAGS) $(GO_COVER_SCAN)

# 运行 Go 依赖漏洞检查
.PHONY: check-go-vuln
check-go-vuln:
	@echo "[Check] -> 运行 Go 依赖漏洞检查 (Tags: $(LINT_TAGS))..."
	$(LINT_ENV) govulncheck -tags=$(LINT_TAGS) -scan=package -show verbose ./...

# 检查后端项目
.PHONY: check-go
check-go: check-go-vet check-go-lint check-go-fix check-go-test check-go-vuln 

# 检查/测试项目
.PHONY: check
check: check-go 

#####################
## 更新 / 拉取 / 推送 相关
#####################

# 按 GIT_REMOTE 顺序先全部拉取合并，再统一推送（保留本地未提交改动）
.PHONY: sync
sync:
	@IFS=,; for r in $$GIT_REMOTE; do \
		git fetch "$$r" && git merge --ff-only "$$r/master" || exit 1; \
	done
	@IFS=,; for r in $$GIT_REMOTE; do git push "$$r" HEAD || exit 1; done
	@echo "同步完成！"

# 推送到所有远程仓库
.PHONY: push
push:
	@IFS=,; for r in $$GIT_REMOTE; do git push "$$r" HEAD || exit 1; done
	@echo "推送完成！"

# 强制推送到所有远程仓库（忽略冲突）
.PHONY: push-force
push-force:
	@IFS=,; for r in $$GIT_REMOTE; do git push --force "$$r" HEAD || exit 1; done
	@echo "强制推送完成！"

#####################
## 测试覆盖率 相关
#####################

# ── 清单配置：新增目录 / 生成代码只改这里，下面的命令不散落路径 ──────────────
# Go 扫描范围：本仓库没有 internal/，入口 main.go 在根目录、业务在各子包，
# 全量扫 ./... 才不会漏掉主包（覆盖率分母要包含命令行入口）。
GO_COVER_SCAN = ./...
# Go 排除目录（包路径前缀，空格分隔）：当前没有 ent 之类的生成代码目录，留空即可。
# 以后引入生成代码时在这里按 /path 列出，下面的命令不用改。
GO_COVER_EXCLUDE =
# 关闭 vet 检查，加快测试与覆盖率统计速度
GO_COVER_TEST_FLAGS = -vet=off

# -coverpkg 只收「扫描范围内且不在排除列表」的包。用 $$(...) 让 go list 在配方执行时才跑，
# 避免任何 make 目标都触发一次包扫描；grep 参数由排除列表生成，加目录不用改命令。
GO_COVER_EXCLUDE_ARGS = $(foreach d,$(GO_COVER_EXCLUDE),-e '$(d)')
ifeq ($(strip $(GO_COVER_EXCLUDE)),)
GO_COVER_PKGS = $$(go list $(GO_COVER_SCAN) | paste -sd, -)
else
GO_COVER_PKGS = $$(go list $(GO_COVER_SCAN) | grep -v $(GO_COVER_EXCLUDE_ARGS) | paste -sd, -)
endif
GO_COVER_ARGS = $(GO_COVER_TEST_FLAGS) -coverpkg=$(GO_COVER_PKGS) -coverprofile=coverage.out $(GO_COVER_SCAN)

# Go 覆盖率（golang 格式）→ coverage.out；CNB 的 testing:coverage 与 Codecov 都读这份
.PHONY: cover-go
cover-go:
	@echo "[Coverage] Go 测试覆盖率 → coverage.out"
	go test $(GO_COVER_ARGS)
	go tool cover -func=coverage.out

# Go 覆盖率 + JUnit（Codecov Test Analytics 用），依赖 CI 镜像预装的 gotestsum
.PHONY: junit-go
junit-go:
	@echo "[Coverage] Go 测试 → coverage.out + go-test-report.xml"
	gotestsum --junitfile go-test-report.xml -- $(GO_COVER_ARGS)
	go tool cover -func=coverage.out

# Coveralls 的 golang parser 会把包前缀压成 basename，改用 Cobertura 保留完整 filename 路径。
# gocover-cobertura 由统一 CI 镜像预装；-strict 让「包/模块信息缺失被静默跳过」直接失败。
.PHONY: cover-go-cobertura
cover-go-cobertura: cover-go
	@echo "[Coverage] coverage.out → coverage.xml（Cobertura）"
	gocover-cobertura -by-files -strict < coverage.out > coverage.xml

#####################
## 同步 相关
#####################

# 同步 CI 模板脚本（缺失下载，sha256 不同则更新）
.PHONY: fetch-scripts
fetch-scripts:
	@echo "[Scripts] 同步 CI 模板脚本..."
	sh scripts/ci/fetch-scripts.sh
