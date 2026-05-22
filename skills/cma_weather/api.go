package cma_weather

import (
	"encoding/json"
	"fmt"
)

// GetCurrentWeather 获取实时天气
// https://weather.cma.cn/api/now/{cityID}
func GetCurrentWeather(cityID string) (json.RawMessage, error) {
	client := getHTTPClient()

	url := fmt.Sprintf("https://weather.cma.cn/api/now/%s", cityID)

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析天气数据失败: %w", err)
	}

	return data, nil
}
