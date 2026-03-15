package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
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

type Reference struct {
	Icon      string `json:"icon,omitempty"`    // 网站图标地址。
	ID        int64  `json:"id"`                // 引用编号1、2、3
	Title     string `json:"title"`             // 网页标题
	URL       string `json:"url"`               // 网页地址
	WebAnchor string `json:"web_anchor"`        // 网站锚文本或网站标题
	Website   string `json:"website,omitempty"` // 站点名称。
	Snippet   string `json:"snippet,omitempty"` // 搜索接口原文片段（非网页原文）
	Content   string `json:"content,omitempty"` // 网站全文或片段，仅当用户开白且enable_full_content设为true时，显示全文
	Date      string `json:"date,omitempty"`    // 网页日期
	Type      string `json:"type,omitempty"`    // 检索资源类型：web:网页 image:图像内容 video:视频内容
	Image     struct {
		URL    string `json:"url,omitempty"`    // 图片链接
		Height string `json:"height,omitempty"` // 图片高度
		Width  string `json:"width,omitempty"`  // 图片宽度
	} `json:"image,omitempty"` // 图片详情
	Video struct {
		URL      string `json:"url,omitempty"`       // 视频链接
		Height   string `json:"height,omitempty"`    // 视频高度
		Width    string `json:"width,omitempty"`     // 视频宽度
		Size     string `json:"size,omitempty"`      // 视频大小，单位Bytes
		Duration string `json:"duration,omitempty"`  // 视频长度，单位秒
		HoverPic string `json:"hover_pic,omitempty"` // 视频封面图
	} `json:"video,omitempty"` // 视频详情
	IsAladdin     bool   `json:"is_aladdin,omitempty"` // 是否为阿拉丁内容。
	Aladdin       string `json:"aladdin,omitempty"`    // 阿拉丁详细内容，参考文档。
	WebExtensions struct {
		Images []struct {
			URL    string `json:"url,omitempty"`    // 图片链接
			Height string `json:"height,omitempty"` // 图片高度
			Width  string `json:"width,omitempty"`  // 图片宽度
		} `json:"images,omitempty"` // 网页相关图片
	} `json:"web_extensions,omitempty"` // 网页相关图片
}

type SearchResult struct {
	RequestID  string      `json:"request_id"`           // 请求ID。
	Code       string      `json:"code,omitempty"`       // 错误码，当发生异常时返回。
	Message    string      `json:"message,omitempty"`    // 错误消息，当发生异常时返回。
	References []Reference `json:"references,omitempty"` // 模型回答详情列表，参考Reference对象表详情。
}

func baiduSearch(apiKey string, requestBody RequestBody) ([]Reference, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://qianfan.baidubce.com/v2/ai_search/web_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var results SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if results.Code != "" {
		return nil, fmt.Errorf("API error: %s", results.Message)
	}

	return results.References, nil
}

func parseRequestBody(query string) (RequestBody, error) {
	var parseData map[string]interface{}
	if err := json.Unmarshal([]byte(query), &parseData); err != nil {
		return RequestBody{}, fmt.Errorf("JSON parse error: %w", err)
	}

	if _, ok := parseData["query"]; !ok {
		return RequestBody{}, fmt.Errorf("query must be present in request body")
	}

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
	var query string

	if len(os.Args) >= 2 {
		query = os.Args[1]
	} else {
		fmt.Print("Enter search query (JSON format): ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		query = input
	}

	requestBody, err := parseRequestBody(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("BAIDU_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: BAIDU_API_KEY must be set in environment")
		os.Exit(1)
	}

	results, err := baiduSearch(apiKey, requestBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
