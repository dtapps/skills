package jiaozhen_factcheck

import (
	"fmt"
	"time"

	"github.com/dtapp/skills/config"
	"resty.dev/v3"
)

const (
	// APIBaseURL 腾讯较真 API 基础地址
	APIBaseURL = "https://openapi.inews.qq.com"
	// APITimeout API 请求超时时间
	APITimeout = 30 * time.Second
)

// BaseRsp 基础响应结构
type BaseRsp struct {
	Code int    `json:"code"` // 0 表示成功
	Msg  string `json:"msg"`  // 消息
}

// getHTTPClient 获取 Resty HTTP 客户端
func getHTTPClient(caller string) *resty.Client {
	client := resty.New()
	client.SetBaseURL(APIBaseURL)
	client.SetTimeout(APITimeout)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("User-Agent", "tencent-news-cli/1.0.0")
	client.SetHeader("Caller-Skill", caller)

	// 添加 API Key（如果有）
	if apiKey := config.GetTencentNewsAPIKey(); apiKey != "" {
		client.SetHeader("Authorization", "Bearer "+apiKey)
	}

	// 开启 Debug 模式（如果配置文件中启用）
	if config.IsRestyDebugEnabled() {
		client.SetDebug(true)
	}

	return client
}

// handleAPIError 处理 API 错误
func handleAPIError(code int, msg string) error {
	// if code == 999 || code == 4006 {
	// 	return fmt.Errorf("API 错误: %s (code: %d)\n\n请访问 https://news.qq.com/exchange?scene=appkey 申请 API Key，\n然后在配置文件中设置 api_keys.tencent_news", msg, code)
	// }
	return fmt.Errorf("API 错误: %s (code: %d)", msg, code)
}
