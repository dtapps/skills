---
name: "tencent-weather"
title: "腾讯天气"
description: "天气信息查询工具，覆盖中国市级和区县级行政区。当用户查询实况天气、天气预报信息时使用。"
metadata:
  author: "TencentNews"
  date: "2026-04-29 20:56:19"
  version: "1.0.0"
---

# 腾讯天气（tencent-weather）

## 描述

天气信息查询工具，覆盖中国市级和区县级行政区。当用户查询实况天气、天气预报信息时使用。

## 使用方法

```bash
dtskills tencent-weather [command] [flags]
```

## 命令
### air-quality

获取天气空气质量

**用法：**
```bash
dtskills tencent-weather air-quality -a <adcode>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--adcode` | `-a` | 地区编码（如北京：110000） | 是 | - |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `aqi` | AQI 指数 |
| `aqi_name` | 空气质量名称 |
| `co` | CO |
| `no2` | NO2 |
| `o3` | O3 |
| `pm10` | PM10 |
| `pm25` | PM2.5 |
| `so2` | SO2 |
| `update_time` |  |

### current

获取实时天气

**用法：**
```bash
dtskills tencent-weather current -a <adcode>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--adcode` | `-a` | 地区编码（如北京：110000） | 是 | - |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `degree` | 温度 |
| `humidity` | 湿度 |
| `precipitation` | 降水量 |
| `pressure` | 气压 |
| `weather` | 天气 |
| `weather_short` | 天气简写 |
| `wind_direction_name` | 风向 |
| `wind_power` | 风力 |

### forecast

获取天气预报
###### 获取天气预报，支持 forecast_1h(1小时预报) 或 forecast_24h(24小时预报)

**用法：**
```bash
dtskills tencent-weather forecast -a <adcode>
```
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
| `--adcode` | `-a` | 地区编码（如北京：110000） | 是 | - |
| `--type` | `-t` | 预报类型: forecast_1h 或 forecast_24h | 否 | forecast_24h |
| `--verbose` | `-v` | 输出完整结果 | 否 | false |

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `aqi` | 空气质量 |
| `aqi_level` | 空气质量等级 |
| `aqi_name` | 空气质量名称 |
| `day_weather` | 白天天气 |
| `day_weather_short` | 白天天气简写 |
| `day_wind_direction` | 白天风向 |
| `day_wind_power` | 白天风力 |
| `max_degree` | 最高温度 |
| `min_degree` | 最低温度 |
| `night_weather` | 夜晚天气 |
| `night_weather_short` | 夜晚天气简写 |
| `night_wind_direction` | 夜晚风向 |
| `night_wind_power` | 夜晚风力 |
| `time` | 时间 |

### ip2city

IP定位获取城市

**用法：**
```bash
dtskills tencent-weather ip2city
```

**输出字段：**

| 字段名 | 描述 |
|--------|------|
| `city` | 城市 |
| `country` | 国家 |
| `district` | 区县 |
| `district_code` | 区县编码 |
| `ip` | IP 地址 |
| `province` | 省份 |

