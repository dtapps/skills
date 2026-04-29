package baidu_baike

// BaikeSimpleResult 百科简洁输出结构体
type BaikeSimpleResult struct {
	LemmaTitle    string   `json:"lemma_title,omitempty" description:"词条标题"`
	Summary       string   `json:"summary,omitempty" description:"词条摘要"`
	AbstractPlain string   `json:"abstract_plain,omitempty" description:"词条概述"`
	URL           string   `json:"url,omitempty" description:"词条页URL"`
	Classify      []string `json:"classify,omitempty" description:"词条分类"`
}
