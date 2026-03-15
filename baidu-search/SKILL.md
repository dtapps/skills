---
name: "baidu-search"
description: "通过百度搜索 API 进行网络搜索。当需要搜索最新信息、新闻或特定主题的网络内容时调用。支持时间范围过滤和结果数量设置。"
---

# 百度搜索 Go 版本

这是一个使用 Go 语言实现的百度搜索技能，通过百度千帆 API 进行网络搜索。

## 使用方法

```bash
go run scripts/main.go '<JSON>'
```

## 请求参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| query | string | 是 | - | 搜索关键词 |
| count | int | 否 | 10 | 返回结果数量，范围 1-50 |
| freshness | string | 否 | null | 时间范围过滤，支持两种格式：<br>1. 简短格式：pd(最近24小时)、pw(最近7天)、pm(最近31天)、py(最近365天)<br>2. 日期格式：YYYY-MM-DDtoYYYY-MM-DD |

## 环境变量

需要设置以下环境变量：

```bash
export BAIDU_API_KEY="your_api_key_here"
```

## 示例

### 基础搜索

```bash
go run scripts/main.go '{"query":"人工智能"}'
```

### 时间范围 - 简短格式

```bash
go run scripts/main.go '{"query":"最新新闻","freshness":"pd"}'
go run scripts/main.go '{"query":"科技动态","freshness":"pw"}'
go run scripts/main.go '{"query":"行业报告","freshness":"pm"}'
go run scripts/main.go '{"query":"历史资料","freshness":"py"}'
```

### 时间范围 - 日期格式

```bash
go run scripts/main.go '{"query":"2025年新闻","freshness":"2025-01-01to2025-12-31"}'
```

### 设置返回数量

```bash
go run scripts/main.go '{"query":"旅游景点","count": 20}'
```

### 组合使用

```bash
go run scripts/main.go '{"query":"Go语言教程","count": 15, "freshness":"pw"}'
```

## 输出格式

返回 JSON 格式的结果，包含搜索引用（references），每个引用包含：
- title: 标题
- url: 链接
- source: 来源
- page_time: 页面时间
- snippet_html: HTML 格式的摘要

## 注意事项

1. API Key 必须通过环境变量 `BAIDU_API_KEY` 设置
2. `count` 参数最大值为 50
3. 时间范围如果未指定，默认不限制时间
4. 日期格式必须严格匹配 `YYYY-MM-DDtoYYYY-MM-DD` 格式

## 打包说明

由于 Go 编译后的二进制文件与操作系统和 CPU 架构相关，**请在你的本地环境自行打包**：

```bash
# 进入 scripts 目录
cd scripts

# 编译（根据你的系统选择目标平台）
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o baidu-search .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o baidu-search .

# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o baidu-search .

# Windows (x86_64)
GOOS=windows GOARCH=amd64 go build -o baidu-search.exe .

# 运行
./baidu-search '{"query":"旅游景点","count":20}'
```
