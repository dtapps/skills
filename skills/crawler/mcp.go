package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/dtapp/skills/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	// MCPTimeout MCP 请求超时时间
	MCPTimeout = 30 * time.Second
)

// MCPClient MCP 客户端封装
type MCPClient struct {
	session *mcp.ClientSession
	url     string
	token   string
}

// NewMCPClient 创建 MCP 客户端
func NewMCPClient(caller string) (*MCPClient, error) {
	// 从配置文件或环境变量获取 MCP URL 和 Token
	mcpURL := config.GetCrawlerMCPURL()
	if mcpURL == "" {
		return nil, fmt.Errorf("未设置 Crawler MCP URL，请在配置文件或环境变量中设置")
	}

	token := config.GetCrawlerMCPToken()

	// 创建 Streamable HTTP Transport
	transport := &mcp.StreamableClientTransport{
		Endpoint: mcpURL,
	}

	// 获取技能信息
	info := GetSkillInfo()

	// 创建 MCP 客户端
	client := mcp.NewClient(&mcp.Implementation{
		Name:    info.Name,
		Version: info.Version,
	}, nil)

	// 建立连接并创建会话
	ctx, cancel := context.WithTimeout(context.Background(), MCPTimeout)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("连接 MCP Server 失败: %w", err)
	}

	return &MCPClient{
		session: session,
		url:     mcpURL,
		token:   token,
	}, nil
}

// Close 关闭 MCP 客户端
func (c *MCPClient) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

// CallTool 调用 MCP Tool
func (c *MCPClient) CallTool(ctx context.Context, toolName string, params map[string]any) (*mcp.CallToolResult, error) {
	return c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: params,
	})
}

// ListTools 列出可用的 MCP Tools
func (c *MCPClient) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, nil)
}
