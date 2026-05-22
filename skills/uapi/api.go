package uapi

import (
	"encoding/json"
	"fmt"
)

// GetWeather 获取天气信息
// https://uapis.cn/docs/api-reference/get-misc-weather
func GetWeather(city string, adcode string) (json.RawMessage, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("city", city).
		SetQueryParam("adcode", adcode)

	resp, err := req.Get("/misc/weather")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析天气数据失败: %w", err)
	}

	return data, nil
}
