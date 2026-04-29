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
	@echo "[配置] 获取配置文件路径..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . config list

# 获取 配置文件路径
# make config-path
.PHONY: config-path
config-path:
	@echo "[配置] 获取配置文件路径..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . config path

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

# 获取 统计数据路径
# make stats-list
.PHONY: stats-list
stats-list:
	@echo "[统计] 获取统计数据路径..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . stats list

#####################

# 获取 百度搜索
# make baidu-search QUERY="人工智能" COUNT=10
# make baidu-search QUERY="人工智能" COUNT=10 VERBOSE="-v"
.PHONY: baidu-search
baidu-search:
	@echo "[获取] 百度搜索..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-search search --query "$(QUERY)" --count $(COUNT) $(VERBOSE)

#####################

# 获取 百度百科
# make baidu-baike KEYWORD="北京"
# make baidu-baike KEYWORD="北京" VERBOSE="-v"
.PHONY: baidu-baike
baidu-baike:
	@echo "[获取] 百度百科..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . baidu-baike search --key "$(KEYWORD)" $(VERBOSE)

#####################

# 获取 腾讯新闻（搜索）
# make tencent-news-search QUERY="人工智能"
# make tencent-news-search QUERY="人工智能" VERBOSE="-v"
.PHONY: tencent-news-search
tencent-news-search:
	@echo "[获取] 腾讯新闻（搜索）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news search --query "$(QUERY)" $(VERBOSE)

# 获取 腾讯新闻（热点）
# make tencent-news-hot
# make tencent-news-hot VERBOSE="-v"
.PHONY: tencent-news-hot
tencent-news-hot:
	@echo "[获取] 腾讯新闻（热点）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news hot $(VERBOSE)

# 获取 腾讯新闻（早报）
# make tencent-news-morning
# make tencent-news-morning VERBOSE="-v"
.PHONY: tencent-news-morning
tencent-news-morning:
	@echo "[获取] 腾讯新闻（早报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news morning $(VERBOSE)

# 获取 腾讯新闻（晚报）
# make tencent-news-evening
# make tencent-news-evening VERBOSE="-v"
.PHONY: tencent-news-evening
tencent-news-evening:
	@echo "[获取] 腾讯新闻（晚报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-news evening $(VERBOSE)

#####################

# 获取 腾讯天气（当前）
# make tencent-weather-current ADCODE=110000
# make tencent-weather-current ADCODE=110000 VERBOSE="-v"
.PHONY: tencent-weather-current
tencent-weather-current:
	@echo "[获取] 腾讯天气（当前）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather current --adcode "$(ADCODE)" $(VERBOSE)

# 获取 腾讯天气（预报）
# make tencent-weather-forecast ADCODE=110000 TYPE=forecast_24h
# make tencent-weather-forecast ADCODE=110000 TYPE=forecast_24h VERBOSE="-v"
.PHONY: tencent-weather-forecast
tencent-weather-forecast:
	@echo "[获取] 腾讯天气（预报）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather forecast --adcode "$(ADCODE)" --type "$(TYPE)" $(VERBOSE)

# 获取 腾讯天气（空气质量）
# make tencent-weather-airquality ADCODE=110000 
# make tencent-weather-airquality ADCODE=110000 VERBOSE="-v"
.PHONY: tencent-weather-airquality
tencent-weather-airquality:
	@echo "[获取] 腾讯天气（空气质量）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather airquality --adcode "$(ADCODE)" $(VERBOSE)

# 获取 腾讯天气（IP定位获取城市）
# make tencent-weather-ip2city
.PHONY: tencent-weather-ip2city
tencent-weather-ip2city:
	@echo "[获取] 腾讯天气（IP定位获取城市）..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . tencent-weather ip2city

#####################

# 腾讯较真事实查证
# make jiaozhen-factcheck-check QUERY="喝柠檬水能抗癌是真的吗"
# make jiaozhen-factcheck-check QUERY="喝柠檬水能抗癌是真的吗" VERBOSE="-v"
.PHONY: jiaozhen-factcheck-check
jiaozhen-factcheck-check:
	@echo "[查证] 腾讯较真事实查证..."
	@mkdir -p $(DTSKILLS_WORK_DIR)
	DTSKILLS_WORK_DIR=$(DTSKILLS_WORK_DIR) go run . jiaozhen-factcheck check --query "$(QUERY)" $(VERBOSE)

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

# 检测是否存在 upx
UPX = $(shell command -v upx 2> /dev/null)

# 定义压缩逻辑：参数 $(1) 为文件路径
define compress_file
	@if [ "$(UPX)" ]; then \
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
.PHONY: build-all
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
	@echo "[构建] Linux AMD64 (通用服务器/PC x86_64)"
	GOOS=linux GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64)

build-linux-arm64:
	@echo "[构建] Linux ARM64 (树莓派4/5, AWS Graviton, 苹果M系Linux虚拟机)"
	GOOS=linux GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64)

build-linux-arm-v7:
	@echo "[构建] Linux ARM v7 (32位, 旧款树莓派, 大部分NAS, 嵌入式盒子)"
	GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v7 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v7)

build-linux-mipsle:
	@echo "[构建] Linux MIPS LE (小端序, 常见于 OpenWrt 路由器/光猫)"
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-mipsle .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-mipsle)

build-linux-riscv64:
	@echo "[构建] Linux RISC-V 64位 (新兴开源指令集架构)"
	GOOS=linux GOARCH=riscv64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64)

build-linux-loong64:
	@echo "[构建] Linux Loong64 (龙芯新版架构)"
	GOOS=linux GOARCH=loong64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-loong64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-loong64)

build-linux-386:
	@echo "[构建] Linux 386 (32位 x86)"
	GOOS=linux GOARCH=386 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-386 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-386)

build-linux-arm-v6:
	@echo "[构建] Linux ARM v6 (树莓派 Zero/1)"
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v6 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-linux-arm-v6)

build-windows-amd64:
	@echo "[构建] Windows AMD64 (主流 Windows PC/服务器)"
	GOOS=windows GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe)

build-windows-arm64:
	@echo "[构建] Windows ARM64 (Surface Pro X, 骁龙笔记本等)"
	GOOS=windows GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe)

build-windows-386:
	@echo "[构建] Windows 386 (32位 Windows PC)"
	GOOS=windows GOARCH=386 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe)

build-darwin-amd64:
	@echo "[构建] macOS Intel (旧款 Macbook/iMac)"
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64)

build-darwin-arm64:
	@echo "[构建] macOS Apple Silicon (M1/M2/M3 芯片)"
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD) -o ./$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	$(call compress_file,./$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64)

#####################
## 打包 (Linux deb/rpm)
#####################

# 生成 nfpm 临时配置文件（替换变量）
define generate_nfpm_config
	@sed -e 's|$${ARCH}|$(1)|g' -e 's|$${VERSION}|$(VERSION)|g' .nfpm.yaml > .nfpm-$(1).yaml
endef

# 打包 Linux deb
.PHONY: package-deb
package-deb: build-linux-amd64
	@echo "[打包] 创建 Linux deb 包..."
	@mkdir -p dist
	$(call generate_nfpm_config,amd64)
	nfpm package -p deb -f .nfpm-amd64.yaml -t dist/dtskills_amd64.deb
	@rm -f .nfpm-amd64.yaml
	@echo "✓ deb 包已创建: dist/dtskills_amd64.deb"

# 打包 Linux rpm
.PHONY: package-rpm
package-rpm: build-linux-amd64
	@echo "[打包] 创建 Linux rpm 包..."
	@mkdir -p dist
	@cp dist/dtskills-linux-amd64 dist/dtskills-linux-x86_64 2>/dev/null || true
	$(call generate_nfpm_config,x86_64)
	nfpm package -p rpm -f .nfpm-x86_64.yaml -t dist/dtskills.x86_64.rpm
	@rm -f .nfpm-x86_64.yaml dist/dtskills-linux-x86_64
	@echo "✓ rpm 包已创建: dist/dtskills.x86_64.rpm"

# 打包所有 Linu deb/rpm 包
.PHONY: package-linux-all
package-linux-all: package-deb package-rpm
	@echo "✓ 所有 Linux deb/rpm  包已创建"

#####################
## 打包 (Windows MSI)
#####################

# 确保版本号是纯数字，例如从 v1.0.0 变成 1.0.0
MSI_VERSION=$(shell echo $(VERSION) | sed 's/v//')

# 打包 Windows MSI (AMD64)
.PHONY: package-win-amd64
package-win-amd64: build-windows-amd64
	@echo "[打包] 创建 Windows MSI 包 (AMD64)..."
	@mkdir -p dist
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="dist/dtskills-windows-amd64.exe"|g' \
	     installer.wxs > dist/installer-amd64.wxs
	wixl -v -a x64 dist/installer-amd64.wxs -o dist/dtskills_amd64.msi
	@rm -f dist/installer-amd64.wxs
	@echo "✓ MSI 包已创建: dist/dtskills_amd64.msi"

# 打包 Windows MSI (ARM64)
.PHONY: package-win-arm64
package-win-arm64: build-windows-arm64
	@echo "[打包] 创建 Windows MSI 包 (ARM64)..."
	@mkdir -p dist
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="dist/dtskills-windows-arm64.exe"|g' \
	     installer.wxs > dist/installer-arm64.wxs
	wixl -v -a arm64 dist/installer-arm64.wxs -o dist/dtskills_arm64.msi
	@rm -f dist/installer-arm64.wxs
	@echo "✓ MSI 包已创建: dist/dtskills_arm64.msi"

# 打包 Windows MSI (386)
.PHONY: package-win-386
package-win-386: build-windows-386
	@echo "[打包] 创建 Windows MSI 包 (386)..."
	@mkdir -p dist
	@sed -e 's|$${VERSION}|$(MSI_VERSION)|g' \
	     -e 's|Source="dtskills.exe"|Source="dist/dtskills-windows-386.exe"|g' \
	     installer.wxs > dist/installer-386.wxs
	wixl -v -a x86 dist/installer-386.wxs -o dist/dtskills_386.msi
	@rm -f dist/installer-386.wxs
	@echo "✓ MSI 包已创建: dist/dtskills_386.msi"

# 打包所有 Windows MSI 包
.PHONY: package-win-all
package-win-all: package-win-amd64 package-win-386
	@echo "✓ 所有 Windows MSI 包已创建"