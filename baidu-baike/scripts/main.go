package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	Version = "1.0.0"
)

// 响应体
type Response struct {
	RequestID string          `json:"request_id"`        // 请求ID
	Code      string          `json:"code"`              // 0 表示成功，其他都是异常
	Message   string          `json:"message,omitempty"` // 错误消息
	Result    json.RawMessage `json:"result,omitempty"`  // 义项内容
}

// https://cloud.baidu.com/doc/qianfan/s/bmh4stpbh
func baiduBaike(apiKey string, searchType string, searchKey string) (*Response, error) {
	client := &http.Client{}

	// 构建 URL 参数
	params := url.Values{}
	params.Add("search_type", searchType)
	params.Add("search_key", searchKey)

	url := fmt.Sprintf("https://appbuilder.baidu.com/v2/baike/lemma/get_content?%s", params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败，状态码: %d", resp.StatusCode)
	}

	var response *Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("无法解析响应: %w", err)
	}

	return response, nil
}

func main() {
	loadEnvFile()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go '<JSON>'")
		os.Exit(1)
	}

	query := os.Args[1]

	var parseData map[string]any
	if err := json.Unmarshal([]byte(query), &parseData); err != nil {
		fmt.Printf("JSON 解析错误: %v\n", err)
		os.Exit(1)
	}

	// 检索类型 lemmaTitle: 按百科词条名 lemmaId: 按百科词条ID
	if _, ok := parseData["search_type"]; !ok {
		fmt.Println("错误: search_type 必须在请求体中")
		os.Exit(1)
	}

	// 检索关键字 与searchType对应的检索条件
	if _, ok := parseData["search_key"]; !ok {
		fmt.Println("错误: search_key 必须在请求体中")
		os.Exit(1)
	}

	searchType := parseData["search_type"].(string)
	searchKey := parseData["search_key"].(string)

	apiKey := os.Getenv("BAIDU_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: BAIDU_API_KEY 必须在环境变量中设置")
		os.Exit(1)
	}

	response, err := baiduBaike(apiKey, searchType, searchKey)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("错误: 响应体序列化错误: %v\n", err)
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
