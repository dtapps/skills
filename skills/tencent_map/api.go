package tencent_map

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCurrentWeather 获取实时天气
// https://lbs.qq.com/service/webService/webServiceGuide/weatherinfo
func GetCurrentWeather(adcode, location string, caller string) (*string, error) {
	client, err := NewMCPClient(caller)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// 构建参数（type 需要作为字符串传递）
	params := map[string]any{
		"type": "now",
	}
	if adcode != "" {
		params["adcode"] = adcode
	}
	if location != "" {
		params["location"] = location
	}

	result, err := client.CallTool(ctx, "weather", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	// 解析 MCP 返回的结果
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("MCP 返回结果为空")
	}

	// 获取文本内容
	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok || textContent == nil {
		return nil, fmt.Errorf("MCP 返回结果格式错误")
	}

	return &textContent.Text, nil
}

// GetForecastWeather 获取天气预报
// https://lbs.qq.com/service/webService/webServiceGuide/weatherinfo
func GetForecastWeather(adcode, location string, caller string) (*string, error) {
	client, err := NewMCPClient(caller)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// 构建参数（type 需要作为字符串传递）
	params := map[string]any{
		"type": "future",
	}
	if adcode != "" {
		params["adcode"] = adcode
	}
	if location != "" {
		params["location"] = location
	}

	result, err := client.CallTool(ctx, "weather", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	// 解析 MCP 返回的结果
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("MCP 返回结果为空")
	}

	// 获取文本内容
	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok || textContent == nil {
		return nil, fmt.Errorf("MCP 返回结果格式错误")
	}

	return &textContent.Text, nil
}

// GetIpLocation 获取IP定位
// https://lbs.qq.com/service/webService/webServiceGuide/position/webServiceIp
func GetIpLocation(ip string, caller string) (*string, error) {
	client, err := NewMCPClient(caller)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// 构建参数
	params := map[string]any{}
	if ip != "" {
		params["ip"] = ip
	}

	result, err := client.CallTool(ctx, "ipLocation", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	// 解析 MCP 返回的结果
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("MCP 返回结果为空")
	}

	// 获取文本内容
	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok || textContent == nil {
		return nil, fmt.Errorf("MCP 返回结果格式错误")
	}

	return &textContent.Text, nil
}
