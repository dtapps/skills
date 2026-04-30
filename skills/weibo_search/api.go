package weibo_search

import (
	"encoding/json"
	"fmt"
)

// HotSearchResponse 微博热搜响应
type HotSearchResponse struct {
	Code int             `json:"code"` // 0 表示成功
	Msg  string          `json:"msg"`  // 消息
	Data json.RawMessage `json:"data"` // 热搜数据
}

// GetHotSearchOptions 获取微博热搜榜选项
type GetHotSearchOptions struct {
	Category string // 榜单类型
	Count    int    // 返回结果数量
}

// GetHotSearch 获取微博热搜榜
func GetHotSearch(opts GetHotSearchOptions) (*HotSearchResponse, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("category", opts.Category)

	if opts.Count > 0 {
		req.SetQueryParam("count", fmt.Sprintf("%d", opts.Count))
	}

	resp, err := req.Get("/open/weibo/hot_search")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var hotResp HotSearchResponse
	if err := json.Unmarshal(resp.Bytes(), &hotResp); err != nil {
		return nil, fmt.Errorf("解析热搜数据失败: %w", err)
	}

	if hotResp.Code != 0 {
		return nil, handleAPIError(hotResp.Code, hotResp.Msg)
	}

	return &hotResp, nil
}

// SmartSearchResponse 微博智搜响应
type SmartSearchResponse struct {
	Code int             `json:"code"` // 0 表示成功
	Msg  string          `json:"msg"`  // 消息
	Data json.RawMessage `json:"data"` // 搜索结果数据
}

// SmartSearchOptions 智搜选项
type SmartSearchOptions struct {
	Query string // 搜索关键词
}

// SmartSearch 执行微博智搜
func SmartSearch(opts SmartSearchOptions) (*SmartSearchResponse, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("query", opts.Query)

	resp, err := req.Get("/open/wis/search_query")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var searchResp SmartSearchResponse
	if err := json.Unmarshal(resp.Bytes(), &searchResp); err != nil {
		return nil, fmt.Errorf("解析搜索数据失败: %w", err)
	}

	if searchResp.Code != 0 {
		return nil, handleAPIError(searchResp.Code, searchResp.Msg)
	}

	return &searchResp, nil
}
