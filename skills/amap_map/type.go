package amap_map

// ForecastItem 预报天气简洁条目
type ForecastItem struct {
	Date           string `json:"date,omitempty" description:"日期"`
	Week           string `json:"week,omitempty" description:"星期"`
	DayWeather     string `json:"dayweather,omitempty" description:"天气现象"`
	NightWeather   string `json:"nightweather,omitempty" description:"天气现象"`
	DayTemp        string `json:"daytemp,omitempty" description:"白天温度"`
	NightTemp      string `json:"nighttemp,omitempty" description:"夜晚温度"`
	DayWind        string `json:"daywind,omitempty" description:"白天风"`
	NightWind      string `json:"nightwind,omitempty" description:"夜晚风"`
	DayPower       string `json:"daypower,omitempty" description:"白天风力"`
	NightPower     string `json:"nightpower,omitempty" description:"夜晚风力"`
	DayTempFloat   string `json:"daytemp_float,omitempty" description:"白天温度"`
	NighttempFloat string `json:"nighttemp_float,omitempty" description:"温度范围"`
}

// WeatherSimpleResult 天气美化输出结构体
// https://lbs.amap.com/api/webservice/guide/api-advanced/weatherinfo
type WeatherSimpleResult struct {
	City      string         `json:"city,omitempty" description:"城市名称"`
	Forecasts []ForecastItem `json:"forecasts,omitempty" description:"预报天气"`
}
