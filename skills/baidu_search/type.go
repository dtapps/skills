package baidu_search

// SearchSimpleResult 搜索简洁输出结构体
type SearchSimpleResult struct {
	Title   string `json:"title,omitempty" description:"网页标题"`
	URL     string `json:"url,omitempty" description:"网页地址"`
	Website string `json:"website,omitempty" description:"站点名称"`
	Site    string `json:"site,omitempty" description:"站点"`
	Content string `json:"content,omitempty" description:"网页内容，显示2000字以内的相关信息原文片段"`
	Date    string `json:"date,omitempty" description:"网页日期"`
	Type    string `json:"type,omitempty" description:"检索资源类型"`
}
