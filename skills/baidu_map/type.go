package baidu_map

// Location 地理位置信息
type Location struct {
	Country  string `json:"country,omitempty" description:"国家名称"`
	Province string `json:"province,omitempty" description:"省份名称"`
	City     string `json:"city,omitempty" description:"城市名称"`
	Name     string `json:"name,omitempty" description:"区县名称"`
	ID       string `json:"id,omitempty" description:"区县id"`
}

// NowWeather 实况数据
type NowWeather struct {
	Text      string  `json:"text,omitempty" description:"天气现象"`
	Temp      int64   `json:"temp,omitempty" description:"温度（℃）"`
	FeelsLike int64   `json:"feels_like,omitempty" description:"体感温度(℃)"`
	Rh        int64   `json:"rh,omitempty" description:"相对湿度(%)"`
	WindClass string  `json:"wind_class,omitempty" description:"风力等级"`
	WindDir   string  `json:"wind_dir,omitempty" description:"风向描述"`
	Prec1h    float64 `json:"prec_1h,omitempty" description:"1小时累计降水量(mm)"`
	Clouds    int64   `json:"clouds,omitempty" description:"云量(%)"`
	Vis       int64   `json:"vis,omitempty" description:"能见度(m)"`
	Aqi       int64   `json:"aqi,omitempty" description:"空气质量指数数值"`
	Pm25      int64   `json:"pm25,omitempty" description:"pm2.5浓度(μg/m3)"`
	Pm10      int64   `json:"pm10,omitempty" description:"pm10浓度(μg/m3)"`
	No2       int64   `json:"no2,omitempty" description:"二氧化氮浓度(μg/m3)"`
	So2       int64   `json:"so2,omitempty" description:"二氧化硫浓度(μg/m3)"`
	O3        int64   `json:"o3,omitempty" description:"臭氧浓度(μg/m3)"`
	Co        int64   `json:"co,omitempty" description:"一氧化碳浓度(mg/m3)"`
	WindAngle int64   `json:"wind_angle,omitempty" description:"风向角度（°）"`
	Uvi       int64   `json:"uvi,omitempty" description:"紫外线指数"`
	Pressure  int64   `json:"pressure,omitempty" description:"气压（hPa）"`
	Dpt       int64   `json:"dpt,omitempty" description:"露点温度（℃）"`
}

// IndexItem 生活指数数据条目
type IndexItem struct {
	Name   string `json:"name,omitempty" description:"生活指数中文名称"`
	Brief  string `json:"brief,omitempty" description:"生活指数概要说明"`
	Detail string `json:"detail,omitempty" description:"生活指数详细说明"`
}

// ForecastItem 预报数据条目
type ForecastItem struct {
	TextDay      string `json:"text_day,omitempty" description:"白天天气现象"`
	TextNight    string `json:"text_night,omitempty" description:"晚上天气现象"`
	High         int64  `json:"high,omitempty" description:"最高温度(℃)"`
	Low          int64  `json:"low,omitempty" description:"最低温度(℃)"`
	WcDay        string `json:"wc_day,omitempty" description:"白天风力"`
	WdDay        string `json:"wd_day,omitempty" description:"白天风向"`
	WcNight      string `json:"wc_night,omitempty" description:"晚上风力"`
	WdNight      string `json:"wd_night,omitempty" description:"晚上风向"`
	WindAngleDay int64  `json:"wind_angle_day,omitempty" description:"白天风向角度（°）"`
	Uvi          int64  `json:"uvi,omitempty" description:"紫外线指数"`
	Pressure     int64  `json:"pressure,omitempty" description:"气压（hPa）"`
	Dpt          int64  `json:"dpt,omitempty" description:"露点温度（℃）"`
	Date         string `json:"date,omitempty" description:"日期，北京时区"`
	Week         string `json:"week,omitempty" description:"星期，北京时区"`
}

// WeatherSimpleResult 天气美化输出结构体
// https://lbsyun.baidu.com/docs/webapi?title=weatherinquiry/weather/base
type WeatherSimpleResult struct {
	Location Location       `json:"location" description:"地理位置信息"`
	Now      NowWeather     `json:"now" description:"实况数据"`
	Indexes  []IndexItem    `json:"indexes,omitempty" description:"生活指数数据"`
	Forecast []ForecastItem `json:"forecast,omitempty" description:"预报数据"`
}
