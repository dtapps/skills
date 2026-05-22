# 默认目标
.DEFAULT_GOAL := generate

#####################
## 变量定义
#####################

# 二进制文件名
BINARY_NAME := dtskills

# 构建目录
BUILD_DIR := dist

# 工作目录（配置文件和统计数据都存储在这里）
DTSKILLS_WORK_DIR = .

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

# 和风天气 获取天气空气质量
# make hefeng-weather-air-quality LOCATION=101280601
# make hefeng-weather-air-quality LOCATION=101280601 MINIFY="--minify"
.PHONY: hefeng-weather-air-quality
hefeng-weather-air-quality:
	@echo "[和风天气] 获取天气空气质量..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . hefeng-weather air-quality --location "$(LOCATION)" $(MINIFY)

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
## 检查相关
#####################

# 安装依赖
.PHONY: deps
deps:
	@echo "[依赖] 安装/更新依赖..."
	go version
	go get -u ./...
	go mod tidy
	go mod verify
	@echo "✓ 依赖安装完成"

# 定义 Lint 过程中需要激活的 Go Build Tags
LINT_TAGS = goolm,stdjson

check:
	@echo "[Check] 自动语法转换与格式化..."
	go fix ./...
	go fmt ./...
	@echo "[Check] 运行 golangci-lint 深度扫描 (Tags: $(LINT_TAGS))..."
	@echo "  -> Linux 平台 (amd64)"
	GOOS=linux GOARCH=amd64 golangci-lint fmt
	GOOS=linux GOARCH=amd64 golangci-lint run --build-tags $(LINT_TAGS)
	GOOS=linux GOARCH=amd64 golangci-lint run --fix --build-tags $(LINT_TAGS)
	
	@echo "  -> macOS 平台 (arm64)"
	GOOS=darwin GOARCH=arm64 golangci-lint fmt
	GOOS=darwin GOARCH=arm64 golangci-lint run --build-tags $(LINT_TAGS)
	GOOS=darwin GOARCH=arm64 golangci-lint run --fix --build-tags $(LINT_TAGS)
	
	@echo "  -> Windows 平台 (amd64)"
	GOOS=windows GOARCH=amd64 golangci-lint fmt
	GOOS=windows GOARCH=amd64 golangci-lint run --build-tags $(LINT_TAGS)
	GOOS=windows GOARCH=amd64 golangci-lint run --fix --build-tags $(LINT_TAGS)

#####################
## 构建相关
#####################

# 功能开关（设置为 1 开启，0 关闭）
ENABLE_UPX ?= 1
ENABLE_ARCHIVE ?= 1

# 检测是否存在 upx
UPX = $(shell command -v upx 2> /dev/null)
# 检测是否存在 tar 和 zip
TAR = $(shell command -v tar 2> /dev/null)
ZIP = $(shell command -v zip 2> /dev/null)

# 定义压缩逻辑：参数 $(1) 为文件路径
define compress_file
	@if [ "$(ENABLE_UPX)" = "1" ] && [ "$(UPX)" ]; then \
		FILE_SIZE=$$(stat -c%s "$(1)" 2>/dev/null || stat -f%z "$(1)"); \
		ORIG_SIZE_H=$$(du -h "$(1)" | cut -f1); \
		\
		if [ $$FILE_SIZE -lt 10485760 ]; then \
			LEVEL="--best"; \
			DESC="小文件-极限压缩"; \
		elif [ $$FILE_SIZE -lt 52428800 ]; then \
			LEVEL="-6"; \
			DESC="中等文件-平衡模式"; \
		else \
			LEVEL="-1"; \
			DESC="大文件-极速模式"; \
		fi; \
		\
		echo "  📦 正在以 [$$DESC ($$LEVEL)] 压缩..."; \
		ERR_MSG=$$($(UPX) $$LEVEL "$(1)" 2>&1); \
		\
		if [ $$? -eq 0 ]; then \
			NEW_SIZE_H=$$(du -h "$(1)" | cut -f1); \
			echo "  ⚡ UPX 压缩完成: $(1) [ $${ORIG_SIZE_H} -> $${NEW_SIZE_H} ]"; \
		else \
			echo "  ⚠️  UPX 压缩跳过: $(1) ($$ERR_MSG)"; \
		fi \
	elif [ "$(ENABLE_UPX)" = "1" ] && [ ! "$(UPX)" ]; then \
		echo "  ⏭️  UPX 未安装，跳过压缩"; \
	elif [ "$(ENABLE_UPX)" = "0" ]; then \
		echo "  ⏭️  UPX 已禁用"; \
	fi
endef

# 定义打包压缩逻辑：参数 $(1) 为源文件路径, $(2) 为目标平台 (如 linux-amd64)
# Linux/macOS 使用 tar.gz, Windows 使用 zip
define archive_binary
	@if [ "$(ENABLE_ARCHIVE)" = "1" ]; then \
		SOURCE_FILE=$(1); \
		PLATFORM=$(2); \
		ARCHIVE_NAME="$(BUILD_DIR)/$(BINARY_NAME)-$${PLATFORM}"; \
		\
		if echo "$${PLATFORM}" | grep -q "windows"; then \
			if [ "$(ZIP)" ]; then \
				ARCHIVE_NAME="$${ARCHIVE_NAME}.zip"; \
				echo "  📦 正在打包 Windows 压缩包: $${ARCHIVE_NAME}..."; \
				cd $(BUILD_DIR) && zip -q "$$(basename $${ARCHIVE_NAME})" "$$(basename $${SOURCE_FILE})"; \
				cd ..; \
				ARCHIVE_SIZE=$$(du -h "$${ARCHIVE_NAME}" | cut -f1); \
				BIN_SIZE=$$(du -h "$${SOURCE_FILE}" | cut -f1); \
				echo "  ✓ 打包完成: $${ARCHIVE_NAME} [ 二进制: $${BIN_SIZE} -> 压缩包: $${ARCHIVE_SIZE} ]"; \
			else \
				echo "  ⚠️  zip 未安装，无法打包 Windows 压缩包"; \
			fi \
		else \
			if [ "$(TAR)" ]; then \
				ARCHIVE_NAME="$${ARCHIVE_NAME}.tar.gz"; \
				echo "  📦 正在打包 Unix 压缩包: $${ARCHIVE_NAME}..."; \
				tar -czf "$${ARCHIVE_NAME}" -C $(BUILD_DIR) "$$(basename $${SOURCE_FILE})"; \
				ARCHIVE_SIZE=$$(du -h "$${ARCHIVE_NAME}" | cut -f1); \
				BIN_SIZE=$$(du -h "$${SOURCE_FILE}" | cut -f1); \
				echo "  ✓ 打包完成: $${ARCHIVE_NAME} [ 二进制: $${BIN_SIZE} -> 压缩包: $${ARCHIVE_SIZE} ]"; \
			else \
				echo "  ⚠️  tar 未安装，无法打包 Unix 压缩包"; \
			fi \
		fi; \
	elif [ "$(ENABLE_ARCHIVE)" = "0" ]; then \
		echo "  ⏭️  压缩包已禁用"; \
	fi
endef

# 定义清理二进制文件逻辑：如果存在对应压缩包则删除二进制
# 参数 $(1) 为二进制文件路径, $(2) 为目标平台 (如 linux-amd64)
define clean_binary
	@SOURCE_FILE=$(1); \
	PLATFORM=$(2); \
	if echo "$${PLATFORM}" | grep -q "windows"; then \
		ARCHIVE_NAME="$(BUILD_DIR)/$(BINARY_NAME)-$${PLATFORM}.zip"; \
	else \
		ARCHIVE_NAME="$(BUILD_DIR)/$(BINARY_NAME)-$${PLATFORM}.tar.gz"; \
	fi; \
	if [ -f "$${ARCHIVE_NAME}" ] && [ -f "$${SOURCE_FILE}" ]; then \
		rm -f "$${SOURCE_FILE}"; \
		echo "  🗑️  已清理二进制: $${SOURCE_FILE} (存在 $${ARCHIVE_NAME})"; \
	else \
		echo "  ⏭️  未找到压缩包或二进制，跳过: $${SOURCE_FILE}"; \
	fi
endef

# 版本信息（支持外部传入）
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDTIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# ldflags 模板
# -s -w: 移除符号表和调试信息，减小体积
# -buildid=: 移除构建指纹
LDFLAGS = -s -w -buildid= \
    -X 'main.Version=$(VERSION)' \
    -X 'main.GitCommit=$(COMMIT)' \
    -X 'main.BuildTime=$(BUILDTIME)' \
    -X 'github.com/dtapp/skills/types.ProjectName=$(BINARY_NAME)'

# 统一定义构建命令
# CGO_ENABLED=0: 禁用 CGO，实现完全静态链接
# -trimpath: 移除本地路径，保护隐私
BUILD_CMD = CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)"

#####################
## 默认目标
#####################

.PHONY: all
all: build

#####################
## 本地构建
#####################

# 构建当前平台
.PHONY: build
build:
	@echo "[构建] 本地平台..."
	$(BUILD_CMD) -o $(BINARY_NAME) .
	@echo "✓ 构建完成: ./$(BINARY_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  提交: $(COMMIT)"
	@echo "  时间: $(BUILDTIME)"

# 清理构建产物
.PHONY: clean
clean:
	@echo "[Clean] 清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@mkdir -p $(BUILD_DIR)
	@echo "✓ 清理完成"

#####################
## 多平台构建
#####################

# 构建所有平台
.PHONY: build-all build-linux-amd64 build-linux-arm64 build-linux-arm-v7 build-linux-mipsle build-linux-riscv64 build-linux-loong64 build-linux-386 build-linux-arm-v6 build-windows-amd64 build-windows-arm64 build-windows-386 build-darwin-amd64 build-darwin-arm64
build-all: clean \
	build-linux-amd64 \
	build-linux-arm64 \
	build-linux-arm-v7 \
	build-linux-mipsle \
	build-linux-riscv64 \
	build-linux-loong64 \
	build-linux-386 \
	build-linux-arm-v6 \
	build-windows-amd64 \
	build-windows-arm64 \
	build-windows-386 \
	build-darwin-amd64 \
	build-darwin-arm64

build-linux-amd64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux AMD64 (通用服务器/PC x86_64)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64,linux-amd64)

build-linux-arm64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux ARM64 (树莓派4/5, AWS Graviton, 苹果M系Linux虚拟机)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64,linux-arm64)

build-linux-arm-v7:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux ARM v7 (32位, 旧款树莓派, 大部分NAS, 嵌入式盒子)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v7 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v7,linux-arm-v7)

build-linux-mipsle:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux MIPS LE (小端序, 常见于 OpenWrt 路由器/光猫)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-mipsle .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-mipsle,linux-mipsle)

build-linux-riscv64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux RISC-V 64位 (新兴开源指令集架构)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=riscv64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64,linux-riscv64)

build-linux-loong64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux Loong64 (龙芯新版架构)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=loong64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-loong64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-loong64,linux-loong64)

build-linux-386:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux 386 (32位 x86)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=386 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-386 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-386,linux-386)

build-linux-arm-v6:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Linux ARM v6 (树莓派 Zero/1)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v6 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v6,linux-arm-v6)

build-windows-amd64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Windows AMD64 (主流 Windows PC/服务器)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=windows GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe,windows-amd64)

build-windows-arm64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Windows ARM64 (Surface Pro X, 骁龙笔记本等)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=windows GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe,windows-arm64)

build-windows-386:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] Windows 386 (32位 Windows PC)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=windows GOARCH=386 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe,windows-386)

build-darwin-amd64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] macOS Intel (旧款 Macbook/iMac)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64,darwin-amd64)

build-darwin-arm64:
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [构建] macOS Apple Silicon (M1/M2/M3 芯片)"
	@echo "└────────────────────────────────────────────────────────────"
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	$(call archive_binary,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64,darwin-arm64)

# 清理所有已压缩的二进制文件（上传前执行）
.PHONY: clean-binaries
clean-binaries:
	@echo "[清理] 删除已压缩的二进制文件..."
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64,linux-amd64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64,linux-arm64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v7,linux-arm-v7)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-mipsle,linux-mipsle)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64,linux-riscv64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-loong64,linux-loong64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-386,linux-386)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v6,linux-arm-v6)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe,windows-amd64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe,windows-arm64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe,windows-386)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64,darwin-amd64)
	$(call clean_binary,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64,darwin-arm64)
	@echo "✓ 清理完成"

#####################
## 打包 (Linux deb/rpm)
#####################

# 生成 nfpm 临时配置文件（替换变量）
define generate_nfpm_config
	@sed -e 's|$${ARCH}|$(1)|g' -e 's|$${VERSION}|$(VERSION)|g' .nfpm.yaml > .nfpm-$(1).yaml
endef

# 打包 Linux deb
.PHONY: package-linux-deb
package-linux-deb: build-linux-amd64
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [打包] 创建 Linux deb 包..."
	@echo "└────────────────────────────────────────────────────────────"
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64)
	$(call generate_nfpm_config,amd64)
	nfpm package -p deb -f .nfpm-amd64.yaml -t $(BUILD_DIR)/dtskills_amd64.deb
	@rm -f .nfpm-amd64.yaml
	@echo "✓ deb 包已创建: $(BUILD_DIR)/dtskills_amd64.deb"

# 打包 Linux rpm
.PHONY: package-linux-rpm
package-linux-rpm: build-linux-amd64
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [打包] 创建 Linux rpm 包..."
	@echo "└────────────────────────────────────────────────────────────"
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64)
	@cp $(BUILD_DIR)/dtskills-linux-amd64 $(BUILD_DIR)/dtskills-linux-x86_64 2>/dev/null || true
	$(call generate_nfpm_config,x86_64)
	nfpm package -p rpm -f .nfpm-x86_64.yaml -t $(BUILD_DIR)/dtskills.x86_64.rpm
	@rm -f .nfpm-x86_64.yaml $(BUILD_DIR)/dtskills-linux-x86_64
	@echo "✓ rpm 包已创建: $(BUILD_DIR)/dtskills.x86_64.rpm"

# 打包所有 Linu deb/rpm 包
.PHONY: package-linux-all
package-linux-all: package-linux-deb package-linux-rpm
	@echo ""
	@echo "╔════════════════════════════════════════════════════════════"
	@echo "║ ✓ 所有 Linux deb/rpm 包已创建"
	@echo "╚════════════════════════════════════════════════════════════"

#####################
## 打包 (Windows MSI)
#####################

# 确保版本号是纯数字，例如从 v1.0.0 变成 1.0.0
MSI_VERSION=$(shell echo $(VERSION) | sed 's/v//')

# 打包 Windows MSI (AMD64)
.PHONY: package-win-amd64
package-win-amd64: build-windows-amd64
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [打包] 创建 Windows MSI 包 (AMD64)..."
	@echo "└────────────────────────────────────────────────────────────"
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe)
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="$(BUILD_DIR)/dtskills-windows-amd64.exe"|g' \
	     installer.wxs > $(BUILD_DIR)/installer-amd64.wxs
	wixl -v -a x64 $(BUILD_DIR)/installer-amd64.wxs -o $(BUILD_DIR)/dtskills_amd64.msi
	@rm -f $(BUILD_DIR)/installer-amd64.wxs
	@echo "✓ MSI 包已创建: $(BUILD_DIR)/dtskills_amd64.msi"

# 打包 Windows MSI (ARM64)
.PHONY: package-win-arm64
package-win-arm64: build-windows-arm64
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [打包] 创建 Windows MSI 包 (ARM64)..."
	@echo "└────────────────────────────────────────────────────────────"
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe)
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="$(BUILD_DIR)/dtskills-windows-arm64.exe"|g' \
	     installer.wxs > $(BUILD_DIR)/installer-arm64.wxs
	wixl -v -a arm64 $(BUILD_DIR)/installer-arm64.wxs -o $(BUILD_DIR)/dtskills_arm64.msi
	@rm -f $(BUILD_DIR)/installer-arm64.wxs
	@echo "✓ MSI 包已创建: $(BUILD_DIR)/dtskills_arm64.msi"

# 打包 Windows MSI (386)
.PHONY: package-win-386
package-win-386: build-windows-386
	@echo ""
	@echo "┌────────────────────────────────────────────────────────────"
	@echo "│ [打包] 创建 Windows MSI 包 (386)..."
	@echo "└────────────────────────────────────────────────────────────"
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe)
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="$(BUILD_DIR)/dtskills-windows-386.exe"|g' \
	     installer.wxs > $(BUILD_DIR)/installer-386.wxs
	wixl -v -a x86 $(BUILD_DIR)/installer-386.wxs -o $(BUILD_DIR)/dtskills_386.msi
	@rm -f $(BUILD_DIR)/installer-386.wxs
	@echo "✓ MSI 包已创建: $(BUILD_DIR)/dtskills_386.msi"

# 打包所有 Windows MSI 包
.PHONY: package-win-all
package-win-all: package-win-amd64 package-win-386
	@echo ""
	@echo "╔════════════════════════════════════════════════════════════"
	@echo "║ ✓ 所有 Windows MSI 包已创建"
	@echo "╚════════════════════════════════════════════════════════════"