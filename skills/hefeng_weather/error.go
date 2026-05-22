package hefeng_weather

import "fmt"

// handleAPIError 处理 API 错误
func handleAPIError(code string, msg string) error {
	return fmt.Errorf("和风天气 API 错误: %s (code: %s)", msg, code)
}
