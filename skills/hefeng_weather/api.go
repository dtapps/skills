package hefeng_weather

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
)

// CurrentWeatherAPIResponse 实时天气 API 响应结构
type CurrentWeatherAPIResponse struct {
	BaseRsp
	Now   CurrentWeatherSimpleResult `json:"now"`   // 实时天气数据
	Refer map[string]any             `json:"refer"` // 参考信息
}

// GetCurrentWeather 获取实时天气
func GetCurrentWeather(location string) (*CurrentWeatherAPIResponse, error) {
	client := getHTTPClient()
	apiKey := config.GetHefengAPIKeyRandom()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("key", apiKey)

	resp, err := req.Get("/weather/now")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result CurrentWeatherAPIResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析实时天气数据失败: %w", err)
	}

	if result.Code != "200" {
		return nil, handleAPIError(result.Code, "和风天气 API 返回异常")
	}

	return &result, nil
}

// DailyForecastAPIResponse 每日天气预报 API 响应结构
type DailyForecastAPIResponse struct {
	BaseRsp
	Daily []DailyForecastSimpleItem `json:"daily"` // 每日预报数据
	Refer map[string]any            `json:"refer"`
}

// GetDailyForecast 获取每日天气预报
// days: 预报天数，支持 3d(3天)、7d(7天)、10d(10天)、15d(15天)、30d(30天)
func GetDailyForecast(location string, days string) (*DailyForecastAPIResponse, error) {
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

	var result DailyForecastAPIResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析天气预报数据失败: %w", err)
	}

	if result.Code != "200" {
		return nil, handleAPIError(result.Code, "和风天气 API 返回异常")
	}

	return &result, nil
}

// AirQualityAPIResponse 空气质量 API 响应结构
type AirQualityAPIResponse struct {
	BaseRsp
	Now   AirQualitySimpleResult `json:"now"` // 实时空气质量数据
	Refer map[string]any         `json:"refer"`
}

// GetAirQuality 获取实时空气质量
func GetAirQuality(location string) (*AirQualityAPIResponse, error) {
	client := getHTTPClient()
	apiKey := config.GetHefengAPIKeyRandom()

	req := client.R().
		SetQueryParam("location", location).
		SetQueryParam("key", apiKey)

	resp, err := req.Get("/air/now")
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result AirQualityAPIResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析空气质量数据失败: %w", err)
	}

	if result.Code != "200" {
		return nil, handleAPIError(result.Code, "和风天气 API 返回异常")
	}

	return &result, nil
}
