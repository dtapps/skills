---
name: "baidu-baike"
title: "百度百科"
description: "百度百科提供百科官方内容检索服务，帮助用户快速、全面地获取权威的百科信息。支持百科词条义项查询、词条详情查询以及秒懂百科视频获取。适用于AI问答、内容创作、教育教学、知识检索等场景，为模型提供可信的知识支撑。"
metadata:
  author: "dtapp"
  date: "2026-04-28 22:09:54"
  version: "1.1.0"
---

# 百度百科（baidu-baike）

## 描述

百度百科提供百科官方内容检索服务，帮助用户快速、全面地获取权威的百科信息。支持百科词条义项查询、词条详情查询以及秒懂百科视频获取。适用于AI问答、内容创作、教育教学、知识检索等场景，为模型提供可信的知识支撑。

## 使用方法

```bash
dtskills baidu-baike [command] [flags]
```

## 命令
### search

查询百度百科词条

**用法：**
```bash
dtskills baidu-baike search -k <key>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--key` | `-k` | 检索关键字 | 是 | - |
| `--type` | `-t` | 检索类型：lemmaTitle(按词条名)、lemmaId(按词条ID) | 否 | lemmaTitle |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `abstract_plain` | 词条概述 |
| `classify` | 词条分类 |
| `lemma_title` | 词条标题 |
| `summary` | 词条摘要 |
| `url` | 词条页URL |
