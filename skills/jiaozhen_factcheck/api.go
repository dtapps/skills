package jiaozhen_factcheck

import (
	"encoding/json"
	"fmt"
)

// SourceInfo 来源信息
type SourceInfo struct {
	Title string `json:"title"` // 来源标题
	URL   string `json:"url"`   // 来源链接
}

// FactCheckResult 事实查证结果
type FactCheckResult struct {
	Conclusion      string       `json:"conclusion"`       // 查证结论
	ConclusionLevel string       `json:"conclusion_level"` // 结论等级
	Explanation     string       `json:"explanation"`      // 解释说明
	Process         string       `json:"process"`          // 查证过程
	Confidence      string       `json:"confidence"`       // 信心评估
	Sources         []SourceInfo `json:"sources"`          // 来源列表
}

// FactCheckData 事实查证数据
type FactCheckData struct {
	QueryID string          `json:"query_id"` // 查询ID
	Query   string          `json:"query"`    // 查询内容
	Result  FactCheckResult `json:"result"`   // 查证结果
}

// FactCheckResponse 事实查证响应
type FactCheckResponse struct {
	BaseRsp BaseRsp         `json:"base_rsp"` // 基础响应
	Data    json.RawMessage `json:"data"`     // 查证数据
}

// ChatCompletionsRequest AI对话请求
type ChatCompletionsRequest struct {
	Query      string `json:"query"`       // 查询内容
	Stream     bool   `json:"stream"`      // 是否流式返回
	AgentAppID string `json:"agent_appid"` // 应用ID
	Caller     string `json:"caller"`      // Skill名称
}

// ChatCompletionsResponse AI对话响应
type ChatCompletionsResponse struct {
	BaseRsp BaseRsp `json:"base_rsp"` // 基础响应
	Data    struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"` // 回复内容
			} `json:"message"`
		} `json:"choices"`
	} `json:"data"`
}

// FactCheck 事实查证 - 使用AI对话接口
func FactCheck(query string, caller string) (*FactCheckResponse, error) {
	client := getHTTPClient(caller)

	// 构建AI对话请求
	reqBody := ChatCompletionsRequest{
		Query:      query,
		Stream:     false,
		AgentAppID: "tencent-news-cli",
		Caller:     caller,
	}

	resp, err := client.R().
		SetBody(reqBody).
		Post("/v1/stream/agent/chat/completions")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	// 解析响应
	var chatResp ChatCompletionsResponse
	if err := json.Unmarshal(resp.Bytes(), &chatResp); err != nil {
		return nil, fmt.Errorf("解析查证数据失败: %w", err)
	}

	// 检查 API 返回码
	if chatResp.BaseRsp.Code != 0 {
		return nil, handleAPIError(chatResp.BaseRsp.Code, chatResp.BaseRsp.Msg)
	}

	// 检查是否是限流或权限错误 (code 在 data 中)
	if len(chatResp.Data.Choices) == 0 {
		// 可能是限流或 appid 错误
		return nil, fmt.Errorf("AI对话接口暂不可用: 可能需要有效的 agent_appid 授权，或已达到每日请求限制")
	}

	// 将AI对话响应转换为事实查证响应格式
	factCheckResp := &FactCheckResponse{
		BaseRsp: chatResp.BaseRsp,
	}

	// 如果有回复内容，封装到Data中
	if len(chatResp.Data.Choices) > 0 {
		content := chatResp.Data.Choices[0].Message.Content
		factCheckData := FactCheckData{
			Query: query,
			Result: FactCheckResult{
				Conclusion:  "AI查证结果",
				Explanation: content,
			},
		}
		dataBytes, _ := json.Marshal(factCheckData)
		factCheckResp.Data = dataBytes
	}

	return factCheckResp, nil
}
