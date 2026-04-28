package baidu_search

import (
	"fmt"
	"time"

	"github.com/dtapp/skills/config"
	"resty.dev/v3"
)

const (
	// APIBaseURL API 基础地址
	APIBaseURL = "https://qianfan.baidubce.com"
	// APITimeout API 请求超时时间
	APITimeout = 30 * time.Second
)

// getHTTPClient 获取 Resty HTTP 客户端
func getHTTPClient() *resty.Client {
	client := resty.New()
	client.SetBaseURL(APIBaseURL)
	client.SetTimeout(APITimeout)
	client.SetHeader("Content-Type", "application/json")

	// 添加 API Key（如果有）
	if apiKey := config.GetBaiduAPIKey(); apiKey != "" {
		client.SetHeader("X-Appbuilder-Authorization", "Bearer "+apiKey)
	}

	// 开启 Debug 模式（如果配置文件中启用）
	if config.IsRestyDebugEnabled() {
		client.SetDebug(true)
	}

	return client
}

// handleAPIError 处理 API 错误
func handleAPIError(code string, msg string) error {
	return fmt.Errorf("API 错误: %s (code: %s)", msg, code)
}
