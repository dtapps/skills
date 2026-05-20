package hefeng_weather

// CurrentWeatherSimpleResult 实时天气简洁输出结构体
type CurrentWeatherSimpleResult struct {
	ObsTime   string `json:"obsTime,omitempty" description:"观测时间"`
	Temp      string `json:"temp,omitempty" description:"温度（℃）"`
	FeelsLike string `json:"feelsLike,omitempty" description:"体感温度（℃）"`
	Icon      string `json:"icon,omitempty" description:"天气图标代码"`
	Text      string `json:"text,omitempty" description:"天气状况"`
	WindDir   string `json:"windDir,omitempty" description:"风向"`
	WindScale string `json:"windScale,omitempty" description:"风力等级"`
	WindSpeed string `json:"windSpeed,omitempty" description:"风速（公里/小时）"`
	Humidity  string `json:"humidity,omitempty" description:"湿度（%）"`
	Precip    string `json:"precip,omitempty" description:"降水量（毫米）"`
	Pressure  string `json:"pressure,omitempty" description:"气压（百帕）"`
	Vis       string `json:"vis,omitempty" description:"能见度（公里）"`
	Cloud     string `json:"cloud,omitempty" description:"云量（%）"`
	Dew       string `json:"dew,omitempty" description:"露点温度（℃）"`
}

// DailyForecastSimpleItem 每日天气预报简洁条目
type DailyForecastSimpleItem struct {
	FxDate         string `json:"fxDate,omitempty" description:"预报日期"`
	Sunrise        string `json:"sunrise,omitempty" description:"日出时间"`
	Sunset         string `json:"sunset,omitempty" description:"日落时间"`
	Moonrise       string `json:"moonrise,omitempty" description:"月出时间"`
	Moonset        string `json:"moonset,omitempty" description:"月落时间"`
	MoonPhase      string `json:"moonPhase,omitempty" description:"月相名称"`
	TempMax        string `json:"tempMax,omitempty" description:"最高温度（℃）"`
	TempMin        string `json:"tempMin,omitempty" description:"最低温度（℃）"`
	IconDay        string `json:"iconDay,omitempty" description:"白天天气图标代码"`
	TextDay        string `json:"textDay,omitempty" description:"白天天气状况"`
	IconNight      string `json:"iconNight,omitempty" description:"夜晚天气图标代码"`
	TextNight      string `json:"textNight,omitempty" description:"夜晚天气状况"`
	WindDirDay     string `json:"windDirDay,omitempty" description:"白天风向"`
	WindScaleDay   string `json:"windScaleDay,omitempty" description:"白天风力等级"`
	WindDirNight   string `json:"windDirNight,omitempty" description:"夜晚风向"`
	WindScaleNight string `json:"windScaleNight,omitempty" description:"夜晚风力等级"`
	Humidity       string `json:"humidity,omitempty" description:"湿度（%）"`
	Precip         string `json:"precip,omitempty" description:"降水量（毫米）"`
	Pressure       string `json:"pressure,omitempty" description:"气压（百帕）"`
	Vis            string `json:"vis,omitempty" description:"能见度（公里）"`
	Cloud          string `json:"cloud,omitempty" description:"云量（%）"`
	UVI            string `json:"uvi,omitempty" description:"紫外线强度指数"`
}

// AirQualitySimpleResult 空气质量简洁输出结构体
type AirQualitySimpleResult struct {
	UpdateTime string `json:"updateTime,omitempty" description:"更新时间"`
	Aqi        string `json:"aqi,omitempty" description:"AQI 指数"`
	Level      string `json:"level,omitempty" description:"空气质量等级"`
	Category   string `json:"category,omitempty" description:"空气质量类别"`
	Primary    string `json:"primary,omitempty" description:"首要污染物"`
	PM10       string `json:"pm10,omitempty" description:"PM10"`
	PM25       string `json:"pm2p5,omitempty" description:"PM2.5"`
	NO2        string `json:"no2,omitempty" description:"二氧化氮"`
	SO2        string `json:"so2,omitempty" description:"二氧化硫"`
	CO         string `json:"co,omitempty" description:"一氧化碳"`
	O3         string `json:"o3,omitempty" description:"臭氧"`
}
