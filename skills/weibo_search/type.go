package weibo_search

// HotSearchSimpleItem 热搜简洁条目
type HotSearchSimpleItem struct {
	ID          int    `json:"id,omitempty" description:"热搜排名"`
	Word        string `json:"word,omitempty" description:"热搜词"`
	Num         int    `json:"num,omitempty" description:"热度值"`
	H5QueryLink string `json:"h5_query_link,omitempty" description:"H5 页面链接"`
}
