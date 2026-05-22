package hefeng_weather

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
)

// GetCurrentWeather 获取实时天气
func GetCurrentWeather(location string) (json.RawMessage, error) {
	client := getHTTPClient()
	apiKey := config.GetHefengAPIKeyRandom()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("key", apiKey)

	resp, err := req.Get("/weather/now")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析实时天气数据失败: %w", err)
	}

	return data, nil
}

// GetDailyForecast 获取每日天气预报
// days: 预报天数，支持 3d(3天)、7d(7天)、10d(10天)、15d(15天)、30d(30天)
func GetDailyForecast(location string, days string) (json.RawMessage, error) {
	client := getHTTPClient()
	apiKey := config.GetHefengAPIKeyRandom()

	// 构建 API 路径: /weather/3d, /weather/7d 等
	apiPath := fmt.Sprintf("/weather/%s", days)

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("key", apiKey)

	resp, err := req.Get(apiPath)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析天气预报数据失败: %w", err)
	}

	return data, nil
}

// GetAirQuality 获取实时空气质量
func GetAirQuality(location string) (json.RawMessage, error) {
	client := getHTTPClient()
	apiKey := config.GetHefengAPIKeyRandom()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("key", apiKey)

	resp, err := req.Get("/air/now")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析空气质量数据失败: %w", err)
	}

	return data, nil
}
