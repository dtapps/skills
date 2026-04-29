package weibo_user

import (
	"encoding/json"
	"fmt"
)

// BaseRsp 基础响应结构
type BaseRsp struct {
	Code int    `json:"code"` // 0 表示成功
	Msg  string `json:"msg"`  // 消息
}

// UserTimelineResponse 用户微博时间线响应
type UserTimelineResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 微博列表数据
	HasMore int             `json:"has_more"` // 是否有更多 (0/1)
	Total   int             `json:"total"`    // 总条数
	TraceID string          `json:"trace_id"` // 追踪 ID
}

// UserTimelineOptions 用户微博时间线选项
type UserTimelineOptions struct {
	Count      int // 返回结果数量
	Page       int // 页码
	TrimUser   int // 是否过滤用户字段 (0/1)
	TrimStatus int // 是否过滤微博字段 (0/1)
}

// GetUserTimeline 获取当前授权用户发布的微博列表
func GetUserTimeline(opts UserTimelineOptions) (*UserTimelineResponse, error) {
	client := getHTTPClient()

	resp, err := client.R().
		SetQueryParam("count", fmt.Sprintf("%d", opts.Count)).
		SetQueryParam("page", fmt.Sprintf("%d", opts.Page)).
		SetQueryParam("trim_user", fmt.Sprintf("%d", opts.TrimUser)).
		SetQueryParam("trim_status", fmt.Sprintf("%d", opts.TrimStatus)).
		Get("/2/statuses/user_timeline.json")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var timelineResp UserTimelineResponse
	if err := json.Unmarshal(resp.Bytes(), &timelineResp); err != nil {
		return nil, fmt.Errorf("解析微博列表数据失败: %w", err)
	}

	if timelineResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(timelineResp.BaseRsp.Code, timelineResp.BaseRsp.Msg)
	}

	return &timelineResp, nil
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 用户数据
	TraceID string          `json:"trace_id"` // 追踪 ID
}

// GetUserInfo 获取当前授权用户信息
func GetUserInfo() (*UserInfoResponse, error) {
	client := getHTTPClient()

	resp, err := client.R().
		Get("/2/account/get_uid.json")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	var userResp UserInfoResponse
	if err := json.Unmarshal(resp.Bytes(), &userResp); err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %w", err)
	}

	if userResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(userResp.BaseRsp.Code, userResp.BaseRsp.Msg)
	}

	return &userResp, nil
}

// WeiboSimpleItem 微博简洁条目
type WeiboSimpleItem struct {
	ID        string `json:"id" description:"微博ID"`
	Text      string `json:"text" description:"微博内容"`
	CreatedAt string `json:"created_at,omitempty" description:"发布时间"`
	Reposts   int    `json:"reposts,omitempty" description:"转发数"`
	Comments  int    `json:"comments,omitempty" description:"评论数"`
	Likes     int    `json:"likes,omitempty" description:"点赞数"`
	Source    string `json:"source,omitempty" description:"来源"`
}

// UserSimpleInfo 用户简洁信息
type UserSimpleInfo struct {
	ID        string `json:"id" description:"用户ID"`
	Name      string `json:"name" description:"昵称"`
	Avatar    string `json:"avatar,omitempty" description:"头像地址"`
	Followers int    `json:"followers,omitempty" description:"粉丝数"`
	Friends   int    `json:"friends,omitempty" description:"关注数"`
}
