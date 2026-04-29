---
name: "weibo-user"
title: "微博用户"
description: "微博用户技能，支持查询当前授权用户发布的微博列表和用户信息。需要 OAuth2 授权。"
metadata:
  author: "dtapp"
  date: "2026-04-29 20:56:19"
  version: "1.0.1"
---

# 微博用户（weibo-user）

## 描述

微博用户技能，支持查询当前授权用户发布的微博列表和用户信息。需要 OAuth2 授权。

## 使用方法

```bash
dtskills weibo-user [command] [flags]
```

## 命令
### my-info

获取当前授权用户信息

**用法：**
```bash
dtskills weibo-user my-info
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `avatar` | 头像地址 |
| `followers` | 粉丝数 |
| `friends` | 关注数 |
| `id` | 用户ID |
| `name` | 昵称 |

### my-posts

获取当前用户发布的微博列表

**用法：**
```bash
dtskills weibo-user my-posts
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--count` | `-c` | 返回结果数量 | 否 | 20 |
| `--page` | `-p` | 页码 | 否 | 1 |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `comments` | 评论数 |
| `created_at` | 发布时间 |
| `id` | 微博ID |
| `likes` | 点赞数 |
| `reposts` | 转发数 |
| `source` | 来源 |
| `text` | 微博内容 |

