package tencent_news

// WeatherSimpleResult 天气简洁结果
type WeatherSimpleResult struct {
	City          string `json:"city" description:"城市"`
	Province      string `json:"province" description:"省份"`
	Weather       string `json:"weather" description:"天气"`
	MinDegree     string `json:"min_degree,omitempty" description:"最低温度"`
	MaxDegree     string `json:"max_degree,omitempty" description:"最高温度"`
	WindDirection string `json:"wind_direction,omitempty" description:"风向"`
	WindPower     string `json:"wind_power,omitempty" description:"风力"`
	Humidity      string `json:"humidity,omitempty" description:"湿度"`
}

// NewsSimpleItem 新闻简洁条目
type NewsSimpleItem struct {
	Title     string `json:"title" description:"标题"`
	Abstract  string `json:"abstract,omitempty" description:"摘要"`
	URL       string `json:"url" description:"链接"`
	Timestamp string `json:"timestamp,omitempty" description:"时间戳"`
	Source    string `json:"source,omitempty" description:"来源"`
}
