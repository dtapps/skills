package baidu_map

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetWeather 获取天气信息
// https://lbs.baidu.com/docs/webapi?title=weatherinquiry/weather
func GetWeather(districtID, location string, isChina bool, caller string) (json.RawMessage, error) {
	client, err := NewMCPClient(caller)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// 构建参数（is_china 需要作为字符串传递）
	params := map[string]any{
		"is_china": "true",
	}
	if !isChina {
		params["is_china"] = "false"
	}
	if districtID != "" {
		params["district_id"] = districtID
	}
	if location != "" {
		params["location"] = location
	}

	result, err := client.CallTool(ctx, "map_weather", params)
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

	// 检查返回内容是否是错误信息（非 JSON 格式）
	text := textContent.Text
	if len(text) > 0 && (text[0] != '{' && text[0] != '[') {
		return nil, fmt.Errorf("MCP 服务返回错误: %s", text)
	}

	var data json.RawMessage
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, fmt.Errorf("解析天气 JSON 失败: %w", err)
	}

	return data, nil
}

// GetIpLocation 获取IP定位信息
// https://lbs.baidu.com/docs/webapi?title=locationip/ip-api
func GetIpLocation(ip string, caller string) (json.RawMessage, error) {
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

	result, err := client.CallTool(ctx, "map_ip_location", params)
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

	// 检查返回内容是否是错误信息（非 JSON 格式）
	text := textContent.Text
	if len(text) > 0 && (text[0] != '{' && text[0] != '[') {
		return nil, fmt.Errorf("MCP 服务返回错误: %s", text)
	}

	var data json.RawMessage
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, fmt.Errorf("解析 IP定位 JSON 失败: %w", err)
	}

	return data, nil
}
