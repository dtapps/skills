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

// BaikeOptions 百科查询选项
type BaikeOptions struct {
	SearchType string // 检索类型：lemmaTitle(按词条名)、lemmaId(按词条ID)
	SearchKey  string // 检索关键字
}

// GetBaiduBaike 查询百度百科
// https://cloud.baidu.com/doc/qianfan/s/bmh4stpbh
func GetBaiduBaike(opts BaikeOptions) (json.RawMessage, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("search_type", opts.SearchType).
		SetQueryParam("search_key", opts.SearchKey)

	resp, err := req.Get("/v2/baike/lemma/get_content")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return data, nil
}
