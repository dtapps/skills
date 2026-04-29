---
name: "weibo-search"
title: "微博搜索"
description: "微博搜索技能，支持微博热搜榜查询和微博智搜（搜索公开微博内容）。"
metadata:
  author: "dtapp"
  date: "2026-04-29 20:55:11"
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
| `--category` | `-t` | 热搜类别 | 是 | - |
| `--count` | `-c` | 返回结果数量 | 否 | 50 |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `category` | 分类标签 |
| `hot` | 热度值 |
| `rank` | 排名 |
| `title` | 热搜标题 |
| `url` | 链接 |

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
| `--count` | `-c` | 返回结果数量 | 否 | 20 |
| `--page` | `-p` | 页码 | 否 | 1 |
| `--query` | `-q` | 搜索关键词 | 是 | - |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `comments` | 评论数 |
| `created_at` | 发布时间 |
| `id` | 微博ID |
| `likes` | 点赞数 |
| `reposts` | 转发数 |
| `text` | 微博内容 |
| `user` | 发布者昵称 |

