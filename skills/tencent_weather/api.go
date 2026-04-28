package tencent_weather

import (
	"encoding/json"
	"fmt"
)

// CurrentWeatherSimpleResult 实时天气简洁输出结构体
type CurrentWeatherSimpleResult struct {
	Weather           string `json:"weather,omitempty" description:"天气"`
	WeatherShort      string `json:"weather_short,omitempty" description:"天气简写"`
	Degree            string `json:"degree,omitempty" description:"温度"`
	Humidity          string `json:"humidity,omitempty" description:"湿度"`
	WindDirectionName string `json:"wind_direction_name,omitempty" description:"风向"`
	WindPower         string `json:"wind_power,omitempty" description:"风力"`
	Pressure          string `json:"pressure,omitempty" description:"气压"`
	Precipitation     string `json:"precipitation,omitempty" description:"降水量"`
}

// CurrentWeatherResponse 实时天气响应
type CurrentWeatherResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 天气详情数据
}

// GetCurrentWeather 获取实时天气
func GetCurrentWeather(adcode string, caller string) (*CurrentWeatherResponse, error) {
	client := getHTTPClient(caller)

	req := client.R().
		SetQueryParam("caller", caller).
		SetQueryParam("weather_type", "observe")

	if adcode != "" {
		req.SetQueryParam("adcode", adcode)
	}

	resp, err := req.Get("/api/v1/agent/weather/comprehensive_weather_info")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result CurrentWeatherResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析天气数据失败: %w", err)
	}

	// 检查 API 返回码
	if result.BaseRsp.Code != 0 {
		return nil, handleAPIError(result.BaseRsp.Code, result.BaseRsp.Msg)
	}

	return &result, nil
}

// ForecastWeatherResponse 预报天气响应
type ForecastWeatherResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 天气详情数据
}

// ForecastWeatherSimpleItem 预报天气简洁条目
type ForecastWeatherSimpleItem struct {
	DayWeather         string `json:"day_weather,omitempty" description:"白天天气"`
	DayWeatherShort    string `json:"day_weather_short,omitempty" description:"白天天气简写"`
	DayWindDirection   string `json:"day_wind_direction,omitempty" description:"白天风向"`
	DayWindPower       string `json:"day_wind_power,omitempty" description:"白天风力"`
	MinDegree          string `json:"min_degree,omitempty" description:"最低温度"`
	MaxDegree          string `json:"max_degree,omitempty" description:"最高温度"`
	NightWeather       string `json:"night_weather,omitempty" description:"夜晚天气"`
	NightWeatherShort  string `json:"night_weather_short,omitempty" description:"夜晚天气简写"`
	NightWindDirection string `json:"night_wind_direction,omitempty" description:"夜晚风向"`
	NightWindPower     string `json:"night_wind_power,omitempty" description:"夜晚风力"`
	Time               string `json:"time,omitempty" description:"时间"`
	Aqi                int    `json:"aqi,omitempty" description:"空气质量"`
	AqiLevel           int    `json:"aqi_level,omitempty" description:"空气质量等级"`
	AqiName            string `json:"aqi_name,omitempty" description:"空气质量名称"`
}

// GetForecastWeather 获取天气预报
func GetForecastWeather(adcode string, weatherType string, caller string) (*ForecastWeatherResponse, error) {
	client := getHTTPClient(caller)

	req := client.R().
		SetQueryParam("caller", caller).
		SetQueryParam("weather_type", weatherType)

	if adcode != "" {
		req.SetQueryParam("adcode", adcode)
	}

	resp, err := req.Get("/api/v1/agent/weather/comprehensive_weather_info")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result ForecastWeatherResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析天气数据失败: %w", err)
	}

	// 检查 API 返回码
	if result.BaseRsp.Code != 0 {
		return nil, handleAPIError(result.BaseRsp.Code, result.BaseRsp.Msg)
	}

	return &result, nil
}

// AirQualitySimpleItem 空气质量简洁条目
type AirQualitySimpleItem struct {
	AQI        int    `json:"aqi,omitempty" description:"AQI 指数"`      // AQI 指数
	AqiName    string `json:"aqi_name,omitempty" description:"空气质量名称"` // 空气质量名称
	CO         string `json:"co,omitempty" description:"CO"`           // CO
	NO2        string `json:"no2,omitempty" description:"NO2"`         // NO2
	O3         string `json:"o3,omitempty" description:"O3"`           // O3
	PM10       string `json:"pm10,omitempty" description:"PM10"`       // PM10
	PM25       string `json:"pm25,omitempty" description:"PM2.5"`      // PM2.5
	SO2        string `json:"so2,omitempty" description:"SO2"`         // SO2
	UpdateTime string `json:"update_time,omitempty"`                   // 更新时间
}

// AirQualityResponse 空气质量响应
type AirQualityResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 空气质量详情数据
}

// GetAirQuality 获取天气空气质量
func GetAirQuality(adcode string, caller string) (*AirQualityResponse, error) {
	client := getHTTPClient(caller)

	req := client.R().
		SetQueryParam("caller", caller).
		SetQueryParam("weather_type", "air")

	if adcode != "" {
		req.SetQueryParam("adcode", adcode)
	}

	resp, err := req.Get("/api/v1/agent/weather/comprehensive_weather_info")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result AirQualityResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析空气质量数据失败: %w", err)
	}

	// 检查 API 返回码
	if result.BaseRsp.Code != 0 {
		return nil, handleAPIError(result.BaseRsp.Code, result.BaseRsp.Msg)
	}

	return &result, nil
}

// IP2CitySimpleResult IP 定位简洁输出结构体
type IP2CitySimpleResult struct {
	IP           string `json:"ip,omitempty" description:"IP 地址"`
	Country      string `json:"country,omitempty" description:"国家"`  // 国家
	Province     string `json:"province,omitempty" description:"省份"` // 省份
	City         string `json:"city,omitempty" description:"城市"`
	District     string `json:"district,omitempty" description:"区县"`
	DistrictCode string `json:"district_code,omitempty" description:"区县编码"`
}

// IP2CityResponse IP 定位响应
type IP2CityResponse struct {
	Ret          int    `json:"ret,omitempty"`          // 返回码
	Msg          string `json:"msg,omitempty"`          // 返回消息
	IP           string `json:"ip,omitempty"`           // IP 地址
	Provcode     string `json:"provcode,omitempty"`     // 省份编码
	Citycode     string `json:"citycode,omitempty"`     // 城市编码
	Country      string `json:"country,omitempty"`      // 国家
	Province     string `json:"province,omitempty"`     // 省份
	City         string `json:"city,omitempty"`         // 城市
	District     string `json:"district,omitempty"`     // 区县
	Isp          string `json:"isp,omitempty"`          // ISP 名称
	DistrictCode string `json:"districtCode,omitempty"` // 区县编码
}

// GetIP2City IP 定位
func GetIP2City(caller string) (*IP2CityResponse, error) {
	client := getHTTPClient(caller)

	resp, err := client.R().
		SetQueryParam("caller", caller).
		Get("/api/v1/agent/ip2city")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var result IP2CityResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("解析 IP 定位数据失败: %w", err)
	}

	// 检查 API 返回码
	if result.Ret != 0 {
		return nil, handleAPIError(result.Ret, result.Msg)
	}

	return &result, nil
}
