package weibo_search

// HotSearchSimpleItem 热搜简洁条目
type HotSearchSimpleItem struct {
	ID          int    `json:"id,omitempty" description:"热搜排名"`
	Word        string `json:"word,omitempty" description:"热搜词"`
	Num         int    `json:"num,omitempty" description:"热度值"`
	H5QueryLink string `json:"h5_query_link,omitempty" description:"H5 页面链接"`
}

// SearchSimpleResult 搜索简洁结果
type SearchSimpleResult struct {
	ID        string `json:"id" description:"微博ID"`
	Text      string `json:"text" description:"微博内容"`
	User      string `json:"user,omitempty" description:"发布者昵称"`
	CreatedAt string `json:"created_at,omitempty" description:"发布时间"`
	Reposts   int    `json:"reposts,omitempty" description:"转发数"`
	Comments  int    `json:"comments,omitempty" description:"评论数"`
	Likes     int    `json:"likes,omitempty" description:"点赞数"`
}
