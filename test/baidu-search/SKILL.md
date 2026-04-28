---
name: "baidu-search"
title: "百度搜索"
description: "百度搜索服务基于百度搜索能力与AI技术，专为生成式AI提供全网实时信息检索服务，具备海量内容站点、高权威性，强时效性，提供全新的搜索体验。"
metadata:
  author: "dtapp"
  date: "2026-04-28 22:09:54"
  version: "1.1.4"
---

# 百度搜索（baidu-search）

## 描述

百度搜索服务基于百度搜索能力与AI技术，专为生成式AI提供全网实时信息检索服务，具备海量内容站点、高权威性，强时效性，提供全新的搜索体验。

## 使用方法

```bash
dtskills baidu-search [command] [flags]
```

## 命令
### search

执行百度搜索

**用法：**
```bash
dtskills baidu-search search -q <query>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--count` | `-c` | 返回结果数量 | 否 | 10 |
| `--freshness` | `-f` | 时间筛选（如：day, week, month） | 否 | - |
| `--query` | `-q` | 搜索关键词 | 是 | - |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `content` | 网页内容，显示2000字以内的相关信息原文片段 |
| `date` | 网页日期 |
| `site` | 站点 |
| `title` | 网页标题 |
| `type` | 检索资源类型 |
| `url` | 网页地址 |
| `website` | 站点名称 |
