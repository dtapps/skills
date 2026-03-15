---
name: "baidu-baike"
description: "通过百度百科 API 查询百科内容。当需要获取百度百科词条信息时调用。"
---

# 百度百科 Go 版本

这是一个使用 Go 语言实现的百度百科查询技能，通过百度千帆 API 查询百科内容。

## 使用方法

```bash
export BAIDU_API_KEY="your_api_key"
go run scripts/main.go '<JSON>'
```

## 请求参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| search_type | string | 是 | - | 搜索类型，目前支持 `lemmaTitle` |
| search_key | string | 是 | - | 搜索关键词（词条标题） |

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

详细 API 文档请参考：[百度百科 API 文档](https://cloud.baidu.com/doc/qianfan/s/bmh4stpbh)

## 示例

### 查询词条

```bash
go run scripts/main.go '{"search_type":"lemmaTitle","search_key":"刘德华"}'
go run scripts/main.go '{"search_type":"lemmaTitle","search_key":"人工智能"}'
go run scripts/main.go '{"search_type":"lemmaTitle","search_key":"Go语言"}'
```

## 输出格式

返回 JSON 格式的结果，包含：
- request_id: 请求 ID
- code: 状态码（0 表示成功，其他表示异常）
- message: 错误消息
- result: 词条结果，包含：
  - lemma_id: int - 词条 ID
  - lemma_title: string - 词条名
  - lemma_desc: string - 义项描述
  - url: string - 词条页 URL（移动设备支持自动适配）
  - summary: string - 词条摘要（优先提取概述，无概述时提取正文前 n 个字，不超过 400 字节）
  - abstract_plain: string - 词条概述（纯文本格式，\n 换行）
  - abstract_html: string - 词条概述（HTML 格式）
  - abstract_structured: any - 词条概述（结构化数据）
  - pic_url: string - 概述图片 URL
  - square_pic_url: string - 概述图片 URL（自动裁切方图，可用于分享、icon 等）
  - square_pic_url_wap: string - 移动端概述图片 URL（自动裁切方图，可用于分享、icon 等）
  - card: list[dict] - 基本信息栏项目
  - card_type: int - 基本信息栏类型
  - albums: list[Album] - 图册列表
  - classify: list[string] - 词条分类，如：["演员","体育人物"]
  - star_map: list[StarMap] - 星图列表
  - videos: list[Video] - 视频列表，包含：
    - second_id: string - 秒懂 ID
    - cover_pic_url: string - 封面原图
    - second_title: string - 秒懂词条名
    - page_url: string - 视频播放页 URL
  - relations: list[Relation] - 关系列表，包含：
    - lemma_id: int - 关联词条 ID
    - lemma_title: string - 关联词条标题
    - relation_name: string - 关系名
    - square_pic_url: string - 关联词条方图

## 注意事项

1. API Key 必须通过环境变量 `BAIDU_API_KEY` 设置
2. `search_type` 目前只支持 `lemmaTitle`
3. `search_key` 是词条的完整标题