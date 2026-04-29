package weibo_search

import (
	"bytes"
	"encoding/json"
	"fmt"

	"charm.land/glamour/v2"
	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"
	"github.com/tidwall/gjson"

	"github.com/spf13/cobra"
)

// WeiboSearchSkill 微博搜索技能
type WeiboSearchSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*WeiboSearchSkill)(nil)

// NewWeiboSearchSkill 创建微博搜索技能实例
func NewWeiboSearchSkill() *WeiboSearchSkill {
	skill := &WeiboSearchSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *WeiboSearchSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *WeiboSearchSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "weibo-search",
		Title:       "微博搜索",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "微博搜索技能，支持微博热搜榜查询和微博智搜（搜索公开微博内容）。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *WeiboSearchSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createHotSearchCmd())
	cmd.AddCommand(s.createSmartSearchCmd())

	return cmd
}

// 命令名称常量
const (
	CmdHotSearch   = "hot-search"
	CmdSmartSearch = "smart-search"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdHotSearch, FieldType: HotSearchSimpleItem{}},
		{CommandName: CmdSmartSearch, FieldType: SearchSimpleResult{}},
	}
}

// requireAPIToken 检查并获取微博 API Token
func (s *WeiboSearchSkill) requireAPIToken(cmd *cobra.Command) (string, bool) {
	apiToken := config.GetWeiboAPITokenRandom()
	if apiToken == "" {
		cmd.Println("错误：未设置微博 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_tokens.weibo <your-api-token>' 设置，或使用 '%s config add api_tokens.weibo <your-api-token>' 添加多个 Token\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiToken, true
}

// createHotSearchCmd 创建热搜榜子命令
func (s *WeiboSearchSkill) createHotSearchCmd() *cobra.Command {
	var (
		category string
		limit    int
		verbose  bool
	)

	cmd := &cobra.Command{
		Use:   CmdHotSearch,
		Short: "获取微博热搜榜",
		Long:  "微博热搜榜工具。当用户需要查看微博热搜、获取热门话题、了解当前热点新闻时激活。",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("正在获取微博热搜榜...")

			// 检查并获取 API Token
			if _, ok := s.requireAPIToken(cmd); !ok {
				return
			}

			opts := GetHotSearchOptions{
				Category: category,
				Limit:    limit,
			}

			resp, err := GetHotSearch(opts)
			if err != nil {
				cmd.Printf("获取热搜榜失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.Data, "", "  "); err != nil {
					cmd.Println(string(resp.Data))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				simpleResults := make([]HotSearchSimpleItem, 0)
				data := gjson.GetBytes(resp.Data, "data").String()
				if err := json.Unmarshal([]byte(data), &simpleResults); err != nil {
					cmd.Println(data)
				} else {
					output, err := json.MarshalIndent(simpleResults, "", "  ")
					if err != nil {
						cmd.Println(data)
					} else {
						cmd.Println(string(output))
					}
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&category, "category", "t", "主榜", "榜单类型（中文）：主榜、文娱榜、社会榜、生活榜、acg榜、科技榜、体育榜")
	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "返回结果数量")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("category")

	return cmd
}

// createSmartSearchCmd 创建智搜子命令
func (s *WeiboSearchSkill) createSmartSearchCmd() *cobra.Command {
	var (
		query   string
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   CmdSmartSearch,
		Short: "执行微博智搜",
		Long:  "微博智搜工具，通过关键词获取微博智搜内容。当用户需要搜索微博内容、查找微博话题、搜索微博用户发布的内容时激活。",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("正在搜索微博: %s\n", query)

			// 检查并获取 API Token
			if _, ok := s.requireAPIToken(cmd); !ok {
				return
			}

			opts := SmartSearchOptions{
				Query: query,
			}

			resp, err := SmartSearch(opts)
			if err != nil {
				cmd.Printf("搜索失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.Data, "", "  "); err != nil {
					cmd.Println(string(resp.Data))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				msg := gjson.GetBytes(resp.Data, "msg").String()
				r, _ := glamour.NewTermRenderer()
				out, err := r.Render(msg)
				if err != nil {
					cmd.Println(msg)
					return
				}
				cmd.Println(out)
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&query, "query", "q", "", "搜索关键词")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewWeiboSearchSkill()
	info := skills.ExtractSkillInfo(skill)

	// 使用统一的映射为每个子命令设置对应的输出字段
	cmdFieldMap := GetCommandOutputFields()
	for i := range info.Commands {
		for _, mapping := range cmdFieldMap {
			if info.Commands[i].Name == mapping.CommandName {
				info.Commands[i].OutputFields = skills.ExtractOutputFieldsFromStruct(mapping.FieldType)
				break
			}
		}
	}

	return info
}

// NewWeiboSearchCmd 创建微博搜索子命令
func NewWeiboSearchCmd() *cobra.Command {
	skill := NewWeiboSearchSkill()
	return skill.GetCommand()
}
