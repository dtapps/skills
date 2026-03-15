---
name: "baidu-search"
description: "通过百度搜索 API 进行网络搜索。当需要搜索最新信息、新闻或特定主题的网络内容时调用。支持时间范围过滤和结果数量设置。"
---

# 百度搜索 Go 版本

这是一个使用 Go 语言实现的百度搜索技能，通过百度千帆 API 进行网络搜索。

## 使用方法

```bash
export BAIDU_API_KEY="your_api_key"
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

或者创建 `.env` 文件（推荐）：

```bash
BAIDU_API_KEY=your_api_key_here
```

## API 文档

详细 API 文档请参考：[百度搜索 API 文档](https://cloud.baidu.com/doc/qianfan-api/s/Wmbq4z7e5)

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

返回 JSON 格式的结果，包含：
- request_id: string - 请求 ID
- code: string - 错误码，当发生异常时返回
- message: string - 错误消息，当发生异常时返回
- references: object - 模型回答详情列表，每个引用包含：
  - icon: string - 网站图标地址
  - id: integer - 引用编号 1、2、3
  - title: string - 网页标题
  - url: string - 网页地址
  - web_anchor: string - 网站锚文本或网站标题
  - website: string - 站点名称
  - snippet: string - 搜索接口原文片段（非网页原文）
  - content: string - 网站全文或片段，仅当用户开白且 enable_full_content 设为 true 时，显示全文
  - date: string - 网页日期
  - type: string - 检索资源类型：web（网页）、image（图像内容）、video（视频内容）
  - image: object - 图片详情
    - url: string - 图片链接
    - height: string - 图片高度
    - width: string - 图片宽度
  - video: object - 视频详情
    - url: string - 视频链接
    - height: string - 视频高度
    - width: string - 视频宽度
    - size: string - 视频大小，单位 Bytes
    - duration: string - 视频长度，单位秒
    - hover_pic: string - 视频封面图
  - is_aladdin: boolean - 是否为阿拉丁内容
  - aladdin: string - 阿拉丁详细内容
  - web_extensions: object - 网页相关图片
    - images: array - 网页相关图片
      - url: string - 图片链接
      - height: string - 图片高度
      - width: string - 图片宽度

## 注意事项

1. API Key 必须通过环境变量 `BAIDU_API_KEY` 设置
2. `count` 参数最大值为 50
3. 时间范围如果未指定，默认不限制时间
4. 日期格式必须严格匹配 `YYYY-MM-DDtoYYYY-MM-DD` 格式
5. `content` 字段仅在白名单用户且 `enable_full_content` 设为 true 时返回