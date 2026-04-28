package tencent_news

import (
	"encoding/json"
	"fmt"
	"time"
)

// WeatherSimpleResult 天气简洁结果
type WeatherSimpleResult struct {
	City          string `json:"city" description:"城市"`
	Province      string `json:"province" description:"省份"`
	Weather       string `json:"weather" description:"天气"`
	MinDegree     string `json:"min_degree,omitempty" description:"最低温度"`
	MaxDegree     string `json:"max_degree,omitempty" description:"最高温度"`
	WindDirection string `json:"wind_direction,omitempty" description:"风向"`
	WindPower     string `json:"wind_power,omitempty" description:"风力"`
	Humidity      string `json:"humidity,omitempty" description:"湿度"`
}

// HotNewsResponse 热点新闻响应
type HotNewsResponse struct {
	BaseRsp     BaseRsp         `json:"base_rsp"`     // 基础响应
	NewsList    json.RawMessage `json:"news_list"`    // 新闻列表
	HasMore     int             `json:"has_more"`     // 是否有更多 (0/1)
	CurrentPage int             `json:"current_page"` // 当前页码
	TraceID     string          `json:"trace_id"`     // 追踪 ID
	Timestamp   string          `json:"timestamp"`    // 时间戳
}

// GetHotNews 获取热点新闻
func GetHotNews(limit int, caller string) (*HotNewsResponse, error) {
	client := getHTTPClient(caller)

	resp, err := client.R().
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("caller", caller).
		Get("/api/v1/agent/hot_rank_list")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var hotResp HotNewsResponse
	if err := json.Unmarshal(resp.Bytes(), &hotResp); err != nil {
		return nil, fmt.Errorf("解析热点新闻数据失败: %w", err)
	}

	if hotResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(hotResp.BaseRsp.Code, hotResp.BaseRsp.Msg)
	}

	return &hotResp, nil
}

// SearchNewsResponse 搜索新闻响应
type SearchNewsResponse struct {
	BaseRsp     BaseRsp         `json:"base_rsp"`     // 基础响应
	NewsList    json.RawMessage `json:"news_list"`    // 新闻列表
	HasMore     int             `json:"has_more"`     // 是否有更多 (0/1)
	CurrentPage int             `json:"current_page"` // 当前页码
	QueryID     string          `json:"query_id"`     // 查询 ID
	TraceID     string          `json:"trace_id"`     // 追踪 ID
	Timestamp   string          `json:"timestamp"`    // 时间戳
}

// SearchNews 搜索新闻
func SearchNews(query string, limit int, page int, caller string) (*SearchNewsResponse, error) {
	client := getHTTPClient(caller)

	// 生成 query_id
	queryID := fmt.Sprintf("%d", time.Now().UnixMilli())

	// 构建请求 body，query 作为嵌套对象，包含 search 和 query_id
	reqBody := map[string]any{
		"query": map[string]string{
			"search":   query,
			"query_id": queryID,
		},
		"limit":  limit,
		"page":   page,
		"caller": caller,
	}

	resp, err := client.R().
		SetBody(reqBody).
		Post("/api/v1/agent/search")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var searchResp SearchNewsResponse
	if err := json.Unmarshal(resp.Bytes(), &searchResp); err != nil {
		return nil, fmt.Errorf("解析搜索数据失败: %w", err)
	}

	if searchResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(searchResp.BaseRsp.Code, searchResp.BaseRsp.Msg)
	}

	return &searchResp, nil
}

// NewsSimpleItem 新闻简洁条目
type NewsSimpleItem struct {
	Title     string `json:"title" description:"标题"`
	Abstract  string `json:"abstract,omitempty" description:"摘要"`
	URL       string `json:"url" description:"链接"`
	Timestamp string `json:"timestamp,omitempty" description:"时间戳"`
	Source    string `json:"source,omitempty" description:"来源"`
}

// BroadcastResponse 早晚报响应
type BroadcastResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 早晚报数据
}

// GetBroadcast 获取早晚报
// broadcastType: 1=早报, 2=晚报
func GetBroadcast(broadcastType int, caller string) (*BroadcastResponse, error) {
	client := getHTTPClient(caller)

	resp, err := client.R().
		SetQueryParam("broadcast_type", fmt.Sprintf("%d", broadcastType)).
		SetQueryParam("caller", caller).
		Get("/api/v1/agent/broadcast_info")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var broadcastResp BroadcastResponse
	if err := json.Unmarshal(resp.Bytes(), &broadcastResp); err != nil {
		return nil, fmt.Errorf("解析早晚报数据失败: %w", err)
	}

	if broadcastResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(broadcastResp.BaseRsp.Code, broadcastResp.BaseRsp.Msg)
	}

	return &broadcastResp, nil
}
