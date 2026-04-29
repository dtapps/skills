package weibo_search

import (
	"fmt"
	"time"

	"github.com/dtapp/skills/config"
	"resty.dev/v3"
)

const (
	// APIBaseURL API 基础地址
	APIBaseURL = "https://open-im.api.weibo.com"
	// APITimeout API 请求超时时间
	APITimeout = 30 * time.Second
	// RetryCount 超时重试次数
	RetryCount = 3
	// RetryWaitTime 首次重试等待时间
	RetryWaitTime = 1 * time.Second
	// RetryMaxWaitTime 最大重试等待时间
	RetryMaxWaitTime = 5 * time.Second
)

// getHTTPClient 获取 Resty HTTP 客户端
func getHTTPClient() *resty.Client {
	client := resty.New()
	client.SetBaseURL(APIBaseURL)
	client.SetTimeout(APITimeout)
	client.SetRetryCount(RetryCount)
	client.SetRetryWaitTime(RetryWaitTime)
	client.SetRetryMaxWaitTime(RetryMaxWaitTime)
	client.SetHeader("Content-Type", "application/json")

	// 添加 API Token（随机选择一个）
	if apiToken := config.GetWeiboAPITokenRandom(); apiToken != "" {
		client.SetQueryParam("token", apiToken)
	}

	// 开启 Debug 模式（如果配置文件中启用）
	if config.IsRestyDebugEnabled() {
		client.SetDebug(true)
	}

	return client
}

// handleAPIError 处理 API 错误
func handleAPIError(code int, msg string) error {
	return fmt.Errorf("API 错误: %s (code: %d)", msg, code)
}
