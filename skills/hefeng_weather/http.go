package hefeng_weather

import (
	"time"

	"github.com/dtapp/skills/config"
	"resty.dev/v3"
)

const (
	// APIBaseURL 和风天气 API 基础地址
	APIBaseURL = "https://api.qweather.com/v7"
	// APITimeout API 请求超时时间
	APITimeout = 30 * time.Second
	// RetryCount 超时重试次数
	RetryCount = 3
	// RetryWaitTime 首次重试等待时间
	RetryWaitTime = 1 * time.Second
	// RetryMaxWaitTime 最大重试等待时间
	RetryMaxWaitTime = 5 * time.Second
)

// BaseRsp 基础响应结构
type BaseRsp struct {
	Code       string `json:"code"`       // "200" 表示成功
	UpdateTime string `json:"updateTime"` // 更新时间
	FxLink     string `json:"fxLink"`     // 网页跳转链接
}

// getHTTPClient 获取 Resty HTTP 客户端
func getHTTPClient() *resty.Client {
	client := resty.New()
	client.SetBaseURL(APIBaseURL)
	client.SetTimeout(APITimeout)
	client.SetRetryCount(RetryCount)
	client.SetRetryWaitTime(RetryWaitTime)
	client.SetRetryMaxWaitTime(RetryMaxWaitTime)
	client.SetHeader("Content-Type", "application/json")

	// 随机获取一个和风天气 API Key
	apiKey := config.GetHefengAPIKeyRandom()
	if apiKey != "" {
		client.SetQueryParam("key", apiKey)
	}

	// 开启 Debug 模式（如果配置文件中启用）
	if config.IsRestyDebugEnabled() {
		client.SetDebug(true)
	}

	return client
}
