package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	Version = "1.0.1"
)

type Message struct {
	Role    string `json:"role"`    // 角色设定，可选值： user：用户 assistant：模型
	Content string `json:"content"` // content为文本时, 对应对话内容，即用户的query问题。说明：不能为空，长度限制在72个字符以内(一个汉字占两个字符)，输入过长的内容时，只取前72个字符检索
}

type ResourceFilter struct {
	Type string `json:"type"`  // 搜索资源类型。可选值： web：网页 video：视频 image：图片 aladdin：阿拉丁
	TopK int    `json:"top_k"` // 指定模态最大返回个数。
}

type TimeRange struct {
	GTE string `json:"gte"`
	LT  string `json:"lt"`
}

type PageTime struct {
	PageTime TimeRange `json:"page_time"`
}

type SearchFilter struct {
	Range PageTime `json:"range"` // 条件查询
}

type RequestBody struct {
	Messages           []Message        `json:"messages"`             // 搜索输入
	Edition            string           `json:"edition"`              // 搜索版本。默认为standard。可选值： standard：完整版本。 lite：标准版本，对召回规模和精排条数简化后的版本，时延表现更好，效果略弱于完整版。
	SearchSource       string           `json:"search_source"`        // 使用的搜索引擎版本。固定值：baidu_search_v2
	ResourceTypeFilter []ResourceFilter `json:"resource_type_filter"` // 支持设置网页、视频、图片、阿拉丁搜索模态，网页top_k最大取值为50，视频top_k最大为10，图片top_k最大为30，阿拉丁top_k最大为5
	SearchFilter       SearchFilter     `json:"search_filter"`        // 根据SearchFilter下的子条件做检索过滤，使用方式参考SearchFilter表详情。
}

// 响应体
type Response struct {
	RequestID  string          `json:"request_id"`           // 请求ID。
	Code       string          `json:"code,omitempty"`       // 错误码，当发生异常时返回。
	Message    string          `json:"message,omitempty"`    // 错误消息，当发生异常时返回。
	References json.RawMessage `json:"references,omitempty"` // 模型回答详情列表，参考Reference对象表详情。
}

// https://cloud.baidu.com/doc/qianfan-api/s/Wmbq4z7e5
func baiduSearch(apiKey string, requestBody RequestBody) (*Response, error) {
	client := &http.Client{}

	// 构建 Body 参数
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("无法序列化请求体: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://qianfan.baidubce.com/v2/ai_search/web_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	req.Header.Set("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败，状态码: %d", resp.StatusCode)
	}

	var results *Response
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("无法解析响应体: %w", err)
	}

	return results, nil
}

func parseRequestBody(query string) (RequestBody, error) {
	var parseData map[string]any
	if err := json.Unmarshal([]byte(query), &parseData); err != nil {
		return RequestBody{}, fmt.Errorf("JSON 解析错误: %w", err)
	}

	// 搜索关键词
	if _, ok := parseData["query"]; !ok {
		return RequestBody{}, fmt.Errorf("查询关键词 query 必须在请求体中")
	}

	// 返回结果数量，范围 1-50
	count := 10
	if val, ok := parseData["count"].(float64); ok {
		count = int(val)
		if count <= 0 {
			count = 10
		} else if count > 50 {
			count = 50
		}
	}

	currentTime := time.Now()
	endDate := currentTime.AddDate(0, 0, 1).Format("2006-01-02")
	searchFilter := SearchFilter{}

	// 时间范围过滤
	if freshness, ok := parseData["freshness"].(string); ok {
		pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}to\d{4}-\d{2}-\d{2}`)
		switch freshness {
		case "pd":
			startDate := currentTime.AddDate(0, 0, -1).Format("2006-01-02")
			searchFilter = SearchFilter{Range: PageTime{PageTime: TimeRange{GTE: startDate, LT: endDate}}}
		case "pw":
			startDate := currentTime.AddDate(0, 0, -6).Format("2006-01-02")
			searchFilter = SearchFilter{Range: PageTime{PageTime: TimeRange{GTE: startDate, LT: endDate}}}
		case "pm":
			startDate := currentTime.AddDate(0, 0, -30).Format("2006-01-02")
			searchFilter = SearchFilter{Range: PageTime{PageTime: TimeRange{GTE: startDate, LT: endDate}}}
		case "py":
			startDate := currentTime.AddDate(-1, 0, 0).Format("2006-01-02")
			searchFilter = SearchFilter{Range: PageTime{PageTime: TimeRange{GTE: startDate, LT: endDate}}}
		default:
			if pattern.MatchString(freshness) {
				parts := regexp.MustCompile(`to`).Split(freshness, 2)
				startDate := parts[0]
				endDate = parts[1]
				searchFilter = SearchFilter{Range: PageTime{PageTime: TimeRange{GTE: startDate, LT: endDate}}}
			} else {
				return RequestBody{}, fmt.Errorf("freshness must be pd, pw, pm, py, or match YYYY-MM-DDtoYYYY-MM-DD format")
			}
		}
	}

	requestBody := RequestBody{
		Messages: []Message{
			{
				Content: parseData["query"].(string),
				Role:    "user",
			},
		},
		SearchSource: "baidu_search_v2",
		ResourceTypeFilter: []ResourceFilter{
			{
				Type: "web",
				TopK: count,
			},
		},
		SearchFilter: searchFilter,
	}

	return requestBody, nil
}

func main() {
	loadEnvFile()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go '<JSON>'")
		os.Exit(1)
	}

	query := os.Args[1]

	requestBody, err := parseRequestBody(query)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("BAIDU_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: BAIDU_API_KEY 必须在环境变量中设置")
		os.Exit(1)
	}

	response, err := baiduSearch(apiKey, requestBody)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func findEnvFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			return envFile
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

func loadEnvFile() {
	envFile := findEnvFile()
	if envFile == "" {
		return
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}
