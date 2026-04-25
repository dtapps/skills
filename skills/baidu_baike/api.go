package baidu_baike

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
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

// Search 查询百度百科
// https://cloud.baidu.com/doc/qianfan/s/bmh4stpbh
func Search(apiKey string, opts BaikeOptions) (*BaikeResponse, error) {
	params := url.Values{}
	params.Add("search_type", opts.SearchType)
	params.Add("search_key", opts.SearchKey)

	apiURL := fmt.Sprintf("https://appbuilder.baidu.com/v2/baike/lemma/get_content?%s", params.Encode())

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败，状态码: %d", resp.StatusCode)
	}

	var result *BaikeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("无法解析响应: %w", err)
	}

	return result, nil
}
