package baidu_search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// Message 搜索消息
type Message struct {
	Role    string `json:"role"`    // 角色设定，可选值： user：用户 assistant：模型
	Content string `json:"content"` // content为文本时, 对应对话内容
}

// ResourceFilter 资源过滤器
type ResourceFilter struct {
	Type string `json:"type"`  // 搜索资源类型：web、video、image、aladdin
	TopK int    `json:"top_k"` // 指定模态最大返回个数
}

// TimeRange 时间范围
type TimeRange struct {
	GTE string `json:"gte"`
	LT  string `json:"lt"`
}

// PageTime 页面时间
type PageTime struct {
	PageTime TimeRange `json:"page_time"`
}

// SearchFilter 搜索过滤器
type SearchFilter struct {
	Range PageTime `json:"range"`
}

// SearchRequest 搜索请求体
type SearchRequest struct {
	Messages           []Message        `json:"messages"`             // 搜索输入
	Edition            string           `json:"edition"`              // 搜索版本：standard、lite
	SearchSource       string           `json:"search_source"`        // 搜索引擎版本：baidu_search_v2
	ResourceTypeFilter []ResourceFilter `json:"resource_type_filter"` // 资源类型过滤
	SearchFilter       SearchFilter     `json:"search_filter"`        // 搜索过滤条件
}

// SearchResponse 搜索响应
type SearchResponse struct {
	RequestID  string          `json:"request_id"`           // 请求ID
	Code       string          `json:"code,omitempty"`       // 错误码
	Message    string          `json:"message,omitempty"`    // 错误消息
	References json.RawMessage `json:"references,omitempty"` // 模型回答详情
}

// SearchOptions 搜索选项
type SearchOptions struct {
	Query     string // 搜索关键词
	Count     int    // 返回结果数量 1-50
	Freshness string // 时间范围：pd(一天内)、pw(一周内)、pm(一月内)、py(一年内) 或 YYYY-MM-DDtoYYYY-MM-DD
}

// DefaultCount 默认返回数量
const DefaultCount = 10

// Search 执行百度搜索
// https://cloud.baidu.com/doc/qianfan-api/s/Wmbq4z7e5
func Search(apiKey string, opts SearchOptions) (*SearchResponse, error) {
	requestBody := buildRequestBody(opts)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("无法序列化请求体: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		"https://qianfan.baidubce.com/v2/ai_search/web_search",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	req.Header.Set("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

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

	var result *SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("无法解析响应体: %w", err)
	}

	return result, nil
}

// buildRequestBody 构建请求体
func buildRequestBody(opts SearchOptions) SearchRequest {
	count := opts.Count
	if count <= 0 {
		count = DefaultCount
	} else if count > 50 {
		count = 50
	}

	currentTime := time.Now()
	endDate := currentTime.AddDate(0, 0, 1).Format("2006-01-02")
	searchFilter := SearchFilter{}

	// 处理时间范围
	if opts.Freshness != "" {
		searchFilter = buildTimeFilter(opts.Freshness, currentTime, endDate)
	}

	return SearchRequest{
		Messages: []Message{
			{
				Content: opts.Query,
				Role:    "user",
			},
		},
		Edition:      "standard",
		SearchSource: "baidu_search_v2",
		ResourceTypeFilter: []ResourceFilter{
			{
				Type: "web",
				TopK: count,
			},
		},
		SearchFilter: searchFilter,
	}
}

// buildTimeFilter 构建时间过滤器
func buildTimeFilter(freshness string, currentTime time.Time, endDate string) SearchFilter {
	pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}to\d{4}-\d{2}-\d{2}`)

	var startDate string
	switch freshness {
	case "pd": // 一天内
		startDate = currentTime.AddDate(0, 0, -1).Format("2006-01-02")
	case "pw": // 一周内
		startDate = currentTime.AddDate(0, 0, -6).Format("2006-01-02")
	case "pm": // 一月内
		startDate = currentTime.AddDate(0, 0, -30).Format("2006-01-02")
	case "py": // 一年内
		startDate = currentTime.AddDate(-1, 0, 0).Format("2006-01-02")
	default:
		if pattern.MatchString(freshness) {
			parts := regexp.MustCompile(`to`).Split(freshness, 2)
			return SearchFilter{
				Range: PageTime{
					PageTime: TimeRange{
						GTE: parts[0],
						LT:  parts[1],
					},
				},
			}
		}
		return SearchFilter{}
	}

	return SearchFilter{
		Range: PageTime{
			PageTime: TimeRange{
				GTE: startDate,
				LT:  endDate,
			},
		},
	}
}
