package sojson

import (
	"encoding/json"
	"fmt"
)

// GetWeather 获取天气信息
// https://www.sojson.com/api/weather.html
func GetWeather(city string) (json.RawMessage, error) {
	client := getHTTPClient()

	url := fmt.Sprintf("http://t.weather.sojson.com/api/weather/city/%s", city)

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
