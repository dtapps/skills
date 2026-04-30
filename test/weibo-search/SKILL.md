---
name: "weibo-search"
title: "微博搜索"
description: "微博搜索技能，支持微博热搜榜查询和微博智搜（搜索公开微博内容）。"
metadata:
  author: "dtapp"
  date: "2026-04-30 11:55:15"
  version: "1.0.0"
---

# 微博搜索（weibo-search）

## 描述

微博搜索技能，支持微博热搜榜查询和微博智搜（搜索公开微博内容）。

## 使用方法

```bash
dtskills weibo-search [command] [flags]
```

## 命令
### hot-search

获取微博热搜榜

###### 微博热搜榜工具。当用户需要查看微博热搜、获取热门话题、了解当前热点新闻时激活。


**用法：**
```bash
dtskills weibo-search hot-search -t <category>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--category` | `-t` | 榜单类型（中文）：主榜、文娱榜、社会榜、生活榜、acg榜、科技榜、体育榜 | 是 | 主榜 |
| `--limit` | `-l` | 返回结果数量 | 否 | 50 |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `h5_query_link` | H5 页面链接 |
| `id` | 热搜排名 |
| `num` | 热度值 |
| `word` | 热搜词 |

### smart-search

执行微博智搜

###### 微博智搜工具，通过关键词获取微博智搜内容。当用户需要搜索微博内容、查找微博话题、搜索微博用户发布的内容时激活。


**用法：**
```bash
dtskills weibo-search smart-search -q <query>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--query` | `-q` | 搜索关键词 | 是 | - |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

