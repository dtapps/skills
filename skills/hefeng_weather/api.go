package hefeng_weather

import (
	"encoding/json"
	"fmt"
)

// GetCurrentWeather 获取实时天气
// https://dev.qweather.com/docs/api/weather/weather-now/
func GetCurrentWeather(location string) (json.RawMessage, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("location", location)

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
// https://dev.qweather.com/docs/api/weather/weather-daily-forecast/
// days: 预报天数，支持 3d(3天)、7d(7天)、10d(10天)、15d(15天)、30d(30天)
func GetDailyForecast(location string, days string) (json.RawMessage, error) {
	client := getHTTPClient()

	// 构建 API 路径: /weather/3d, /weather/7d 等
	apiPath := fmt.Sprintf("/weather/%s", days)

	req := client.R().
		SetQueryParam("location", location)

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

// GetIndicesForecast 获取天气指数预报
// https://dev.qweather.com/docs/api/indices/indices-forecast/
// days: 预报天数，支持 1d(1天)、3d(3天)
func GetIndicesForecast(location string, days string) (json.RawMessage, error) {
	client := getHTTPClient()

	// 构建 API 路径: /indices/1d, /indices/3d 等
	apiPath := fmt.Sprintf("/indices/%s", days)

	req := client.R().
		SetQueryParam("type", "0").
		SetQueryParam("location", location)

	resp, err := req.Get(apiPath)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析天气指数预报数据失败: %w", err)
	}

	return data, nil
}

// GetSunriseSunset 获取日出日落
// https://dev.qweather.com/docs/api/astronomy/sunrise-sunset/
// date: 日期，格式为 yyyyMMdd，如 20240101
func GetSunriseSunset(location string, date string) (json.RawMessage, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("date", date)

	resp, err := req.Get("/astronomy/sun")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析日出日落数据失败: %w", err)
	}

	return data, nil
}

// GetMoonPhase 获取月升月落和月相
// https://dev.qweather.com/docs/api/astronomy/moon-and-moon-phase/
// date: 日期，格式为 yyyyMMdd，如 20240101
func GetMoonPhase(location string, date string) (json.RawMessage, error) {
	client := getHTTPClient()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("date", date)

	resp, err := req.Get("/astronomy/moon")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var data json.RawMessage
	if err := json.Unmarshal(resp.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("解析月升月落和月相数据失败: %w", err)
	}

	return data, nil
}
