package crawler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetGoldPriceData 获取黄金价格全部数据
// tool: get_gold_price_data
func GetGoldPriceData(date, source string) (json.RawMessage, error) {
	client, err := NewMCPClient("crawler")
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	params := map[string]any{
		"date":   date,
		"source": source,
	}

	result, err := client.CallTool(ctx, "get_gold_price_data", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	return parseMCPResult(result)
}

// GetLatestGoldPrice 获取最新黄金价格
// tool: get_latest_gold_price
func GetLatestGoldPrice(date, source string) (json.RawMessage, error) {
	client, err := NewMCPClient("crawler")
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	params := map[string]any{
		"date":   date,
		"source": source,
	}

	result, err := client.CallTool(ctx, "get_latest_gold_price", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	return parseMCPResult(result)
}

// GetJDFinanceGoldStatus 获取京东金融黄金市场状态
// tool: get_jdfinance_gold_status
func GetJDFinanceGoldStatus(datetime string) (json.RawMessage, error) {
	client, err := NewMCPClient("crawler")
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	params := map[string]any{
		"datetime": datetime,
	}

	result, err := client.CallTool(ctx, "get_jdfinance_gold_status", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	return parseMCPResult(result)
}

// GetHotData 获取热搜数据
// tool: get_hot_data
func GetHotData(date, source string) (json.RawMessage, error) {
	client, err := NewMCPClient("crawler")
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	params := map[string]any{
		"date": date,
	}
	if source != "" {
		params["source"] = source
	}

	result, err := client.CallTool(ctx, "get_hot_data", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	return parseMCPResult(result)
}

// GetWeatherData 获取天气数据
// tool: get_weather_data
func GetWeatherData(cityName, date, source string) (json.RawMessage, error) {
	client, err := NewMCPClient("crawler")
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	params := map[string]any{
		"city_name": cityName,
	}
	if date != "" {
		params["date"] = date
	}
	if source != "" {
		params["source"] = source
	}

	result, err := client.CallTool(ctx, "get_weather_data", params)
	if err != nil {
		return nil, fmt.Errorf("调用 MCP Tool 失败: %w", err)
	}

	return parseMCPResult(result)
}

// parseMCPResult 解析 MCP 返回结果
func parseMCPResult(result *mcp.CallToolResult) (json.RawMessage, error) {
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("MCP 返回结果为空")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok || textContent == nil {
		return nil, fmt.Errorf("MCP 返回结果格式错误")
	}

	text := textContent.Text
	if len(text) > 0 && (text[0] != '{' && text[0] != '[') {
		return nil, fmt.Errorf("MCP 服务返回错误: %s", text)
	}

	var data json.RawMessage
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return data, nil
}
