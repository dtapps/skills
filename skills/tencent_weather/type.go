package tencent_weather

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
