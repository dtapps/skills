# 腾讯新闻 API 文档

> 通过抓包分析 tencent-news-cli 二进制文件获取

## 基础信息

- **正式环境**: `https://openapi.inews.qq.com`
- **测试环境**: `https://test-openapi.inews.qq.com`
- **API Key 申请**: https://news.qq.com/exchange?scene=appkey
- **环境变量**: `TENCENT_NEWS_APIKEY`

## API 端点

| 功能 | 方法 | 端点 | 状态 |
|------|------|------|------|
| 热点新闻榜 | GET | `/api/v1/agent/hot_rank_list` | ✅ 可用 |
| 搜索新闻 | GET | `/api/v1/agent/search` | ✅ 可用 |
| 早晚报 | GET | `/api/v1/agent/broadcast_info` | ✅ 可用 |
| IP 定位 | GET | `/api/v1/agent/ip2city` | ✅ 可用 |
| CLI 版本 | GET | `/api/v1/agent/cli/version` | ✅ 可用 |
| 反馈 | POST | `/api/v1/agent/skill/feedback` | ✅ 可用 |
| 天气 | GET | `/api/v1/agent/weather/comprehensive_weather_info` | ✅ 可用，需要 weather_type 参数 |
| AI 对话 | POST | `/v1/stream/agent/chat/completions` | ⚠️ 需要 agent_appid，可能有限流 |

## 请求头

```
Authorization: Bearer {api_key}
Caller-Skill: {skill_name}
User-Agent: tencent-news-cli/{version}
Content-Type: application/json (for POST)
```

## 通用响应结构

```json
{
  "code": 0,
  "msg": "success",
  "data": {},
  "trace_id": "string",
  "timestamp": "string"
}
```

## 热点新闻榜

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 返回新闻条数，默认 10 |
| caller | string | 否 | Skill 名称，用于上报 |

### 响应字段

```json
{
  "code": 0,
  "msg": "string",
  "data": {
    "news_list": [
      {
        "id": "string",
        "title": "string",
        "abstract": "string",
        "url": "string",
        "source": "string",
        "cover_image": "string",
        "article_type": "string"
      }
    ],
    "has_more": true,
    "current_page": 1
  },
  "trace_id": "string",
  "timestamp": "string"
}
```

## 搜索新闻

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| query | string | 是 | 搜索关键词 |
| limit | int | 否 | 返回条数，默认 10 |
| page | int | 否 | 页码，默认 1 |
| caller | string | 否 | Skill 名称 |

### 响应字段

```json
{
  "code": 0,
  "msg": "string",
  "data": {
    "news_list": [...],
    "has_more": true,
    "current_page": 1,
    "query_id": "string"
  },
  "trace_id": "string"
}
```

## 早晚报

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| broadcast_type | int | 是 | 类型：1(早报) / 2(晚报) |
| caller | string | 否 | Skill 名称 |

### 响应字段

```json
{
  "code": 0,
  "msg": "string",
  "data": {
    "broadcast_type": 1,
    "start_audio": "string",
    "end_audio": "string",
    "news_content": "string",
    "news_list": [...]
  }
}
```

## IP 定位

### 响应字段

```json
{
  "code": 0,
  "msg": "string",
  "data": {
    "ip": "string",
    "city": "string",
    "province": "string",
    "country": "string"
  }
}
```

## 天气查询

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| adcode | string | 是 | 地区编码，例如 110000(北京) |
| weather_type | string | 是 | 天气类型：observe(实况), forecast_1h(1小时预报), forecast_24h(24小时预报), air(空气质量), tips(生活指数), alarm(预警), rise(日出日落), index(指数), limit(限行), air_forecast(空气质量预报), limit_forecast(限行预报) |
| caller | string | 否 | Skill 名称 |

### 响应字段

根据 weather_type 不同，返回不同的数据：

- **observe**: 实况天气（温度、湿度、风向等）
- **forecast_1h/forecast_24h**: 天气预报
- **air**: 空气质量（AQI、等级等）
- **tips**: 生活指数（穿衣、雨伞、运动、感冒）
- **alarm**: 天气预警
- **rise**: 日出日落时间
- **index**: 各类指数
- **limit**: 限行信息
- **air_forecast**: 空气质量预报
- **limit_forecast**: 限行预报

### 请求示例

```bash
# 获取实况天气
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/weather/comprehensive_weather_info?adcode=110000&weather_type=observe"

# 获取24小时预报
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/weather/comprehensive_weather_info?adcode=110000&weather_type=forecast_24h"

# 获取空气质量
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/weather/comprehensive_weather_info?adcode=110000&weather_type=air"
```

## AI 对话 (Chat Completions)

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| query | string | 是 | 查询内容 |
| stream | bool | 否 | 是否流式返回 |
| agent_appid | string | 是 | 应用 ID |
| caller | string | 否 | Skill 名称 |

### 请求示例

```bash
curl -X POST "https://openapi.inews.qq.com/v1/stream/agent/chat/completions" \
  -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: jiaozhen-factcheck" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "喝柠檬水能抗癌是真的吗",
    "stream": false,
    "agent_appid": "tencent-news-cli"
  }'
```

### 注意事项

- 此接口可能有每日请求次数限制
- 需要正确的 agent_appid 才能访问
- 用于较真查证等 AI 功能

## 新闻条目字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 新闻唯一标识 |
| title | string | 新闻标题 |
| abstract | string | 新闻摘要 |
| url | string | 新闻详情页链接 |
| source | string | 新闻来源 |
| cover_image | string | 封面图片 URL |
| article_type | string | 文章类型 |
| audio_url | string | 音频链接（如果有） |
| audio_duration | string | 音频时长 |
| publish_time | string | 发布时间 |

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| -1 | 未设置 API Key 或请求格式错误 |
| 999 | authorization is empty |
| 10001 | 参数错误（如 weatherType is empty） |
| 4004 | 超出每日请求总量最大阈值 |
| 50001 | 参数值非法 |

## 使用示例

### 获取热点新闻

```bash
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/hot_rank_list?limit=10"
```

### 搜索新闻

```bash
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/search?query=AI&limit=10"
```

### 获取早报

```bash
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/broadcast_info?broadcast_type=1"
```

### IP 定位

```bash
curl -H "Authorization: Bearer {api_key}" \
  -H "Caller-Skill: my-skill" \
  "https://openapi.inews.qq.com/api/v1/agent/ip2city"
```

## 注意事项

1. 所有 API 调用都需要在请求头中设置 `Authorization` 和 `Caller-Skill`
2. 部分接口可能需要 API Key 认证
3. API 可能有频率限制
4. 响应数据可能需要处理分页（`has_more` 和 `current_page` 字段）
5. 天气 API 需要 `weatherType` 参数，但具体值需要进一步确认
6. AI 对话接口有每日请求次数限制
