package baidu_baike

import (
	"encoding/json"
	"fmt"
)

// BaikeResponse 百科响应
type BaikeResponse struct {
	RequestID string          `json:"request_id"`        // 请求ID
	Code      string          `json:"code"`              // 0 表示成功，其他都是异常
	Message   string          `json:"message,omitempty"` // 错误消息
	Result    json.RawMessage `json:"result,omitempty"`  // 义项内容
}

// BaikeSimpleResult 百科简洁输出结构体
type BaikeSimpleResult struct {
	LemmaTitle    string   `json:"lemma_title,omitempty" description:"词条标题"`
	Summary       string   `json:"summary,omitempty" description:"词条摘要"`
	AbstractPlain string   `json:"abstract_plain,omitempty" description:"词条概述"`
	URL           string   `json:"url,omitempty" description:"词条页URL"`
	Classify      []string `json:"classify,omitempty" description:"词条分类"`
}

// BaikeOptions 百科查询选项
type BaikeOptions struct {
	SearchType string // 检索类型：lemmaTitle(按词条名)、lemmaId(按词条ID)
	SearchKey  string // 检索关键字
}

// GetBaiduBaike 查询百度百科
// https://cloud.baidu.com/doc/qianfan/s/bmh4stpbh
func GetBaiduBaike(apiKey string, opts BaikeOptions) (*BaikeResponse, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("search_type", opts.SearchType).
		SetQueryParam("search_key", opts.SearchKey)

	resp, err := req.Get("/v2/baike/lemma/get_content")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result BaikeResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查 API 返回码
	if result.Code != "" && result.Code != "0" {
		return nil, handleAPIError(result.Code, result.Message)
	}

	return &result, nil
}
