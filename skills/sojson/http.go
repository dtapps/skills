package sojson

import (
	"time"

	"github.com/dtapp/skills/config"
	"resty.dev/v3"
)

const (
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
	client.SetTimeout(APITimeout)
	client.SetRetryCount(RetryCount)
	client.SetRetryWaitTime(RetryWaitTime)
	client.SetRetryMaxWaitTime(RetryMaxWaitTime)
	client.SetHeader("Content-Type", "application/json")

	// 开启 Debug 模式（如果配置文件中启用）
	if config.IsRestyDebugEnabled() {
		client.SetDebug(true)
	}

	return client
}
