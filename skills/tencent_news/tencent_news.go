package tencent_news

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"
	"github.com/tidwall/gjson"

	"github.com/spf13/cobra"
)

// TencentNewsSkill 腾讯新闻技能
type TencentNewsSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*TencentNewsSkill)(nil)

// NewTencentNewsSkill 创建腾讯新闻技能实例
func NewTencentNewsSkill() *TencentNewsSkill {
	skill := &TencentNewsSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *TencentNewsSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
// https://skillhub.cn/skills/tencent-news
func (s *TencentNewsSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "tencent-news",
		Title:       "腾讯新闻",
		Version:     "1.0.9",
		Author:      "dtapp",
		Description: "7×24 新闻资讯搜索工具，聚焦中国国内信息和国际热点。支持新闻搜索，包括热点新闻、早报晚报、实时资讯、领域新闻和天气信息查询。当用户需要搜索新闻、新闻热榜、新闻早晚报、订阅新闻推送、获取相关新闻资讯和查询天气信息时使用",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdSearch  = "search"
	CmdHot     = "hot"
	CmdMorning = "morning"
	CmdEvening = "evening"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdSearch, FieldType: NewsSimpleItem{}},
		{CommandName: CmdHot, FieldType: NewsSimpleItem{}},
		{CommandName: CmdMorning, FieldType: NewsSimpleItem{}},
		{CommandName: CmdEvening, FieldType: NewsSimpleItem{}},
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *TencentNewsSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createSearchCmd())
	cmd.AddCommand(s.createHotCmd())
	cmd.AddCommand(s.createMorningCmd())
	cmd.AddCommand(s.createEveningCmd())

	return cmd
}

// requireAPIKey 检查并获取腾讯新闻 API Key
func (s *TencentNewsSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetTencentNewsAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置腾讯新闻 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.tencent_news <your-api-key>' 设置，或使用 '%s config add api_keys.tencent_news <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createSearchCmd 创建搜索子命令
func (s *TencentNewsSkill) createSearchCmd() *cobra.Command {
	var (
		query   string
		count   int
		page    int
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   CmdSearch,
		Short: "搜索新闻",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("正在搜索新闻: %s\n", query)

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := SearchNews(query, count, page, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("搜索失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 NewsList）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.NewsList, "", "  "); err != nil {
					cmd.Println(string(resp.NewsList))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果
				simpleResults := make([]NewsSimpleItem, 0)
				if err := json.Unmarshal([]byte(resp.NewsList), &simpleResults); err != nil {
					cmd.Println(string(resp.NewsList))
				} else {
					output, err := json.MarshalIndent(simpleResults, "", "  ")
					if err != nil {
						cmd.Println(string(resp.NewsList))
					} else {
						cmd.Println(string(output))
					}
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&query, "query", "q", "", "搜索关键词")
	cmd.Flags().IntVarP(&count, "count", "c", 10, "返回结果数量")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "页码")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

// createHotCmd 创建热点新闻子命令
func (s *TencentNewsSkill) createHotCmd() *cobra.Command {
	var (
		count   int
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   CmdHot,
		Short: "获取热点新闻",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("正在获取热点新闻...")

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetHotNews(count, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取热点新闻失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 NewsList）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.NewsList, "", "  "); err != nil {
					cmd.Println(string(resp.NewsList))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果
				simpleResults := make([]NewsSimpleItem, 0)
				if err := json.Unmarshal([]byte(resp.NewsList), &simpleResults); err != nil {
					cmd.Println(string(resp.NewsList))
				} else {
					output, err := json.MarshalIndent(simpleResults, "", "  ")
					if err != nil {
						cmd.Println(string(resp.NewsList))
					} else {
						cmd.Println(string(output))
					}
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().IntVarP(&count, "count", "c", 10, "返回结果数量")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	return cmd
}

// createMorningCmd 创建新闻早报子命令
func (s *TencentNewsSkill) createMorningCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   CmdMorning,
		Short: "获取新闻早报",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("正在获取新闻早报...")

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetBroadcast(1, s.GetInfo().Name) // 1=早报
			if err != nil {
				cmd.Printf("获取新闻早报失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 Data）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.Data, "", "  "); err != nil {
					cmd.Println(string(resp.Data))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果
				simpleResults := make([]NewsSimpleItem, 0)
				observe := gjson.GetBytes(resp.Data, "data.news_list").String()
				if err := json.Unmarshal([]byte(observe), &simpleResults); err != nil {
					cmd.Println(observe)
				} else {
					output, err := json.MarshalIndent(simpleResults, "", "  ")
					if err != nil {
						cmd.Println(observe)
					} else {
						cmd.Println(string(output))
					}
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	return cmd
}

// createEveningCmd 创建新闻晚报子命令
func (s *TencentNewsSkill) createEveningCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   CmdEvening,
		Short: "获取新闻晚报",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("正在获取新闻晚报...")

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetBroadcast(2, s.GetInfo().Name) // 2=晚报
			if err != nil {
				cmd.Printf("获取新闻晚报失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 Data）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.Data, "", "  "); err != nil {
					cmd.Println(string(resp.Data))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果
				simpleResults := make([]NewsSimpleItem, 0)
				observe := gjson.GetBytes(resp.Data, "data.news_list").String()
				if err := json.Unmarshal([]byte(observe), &simpleResults); err != nil {
					cmd.Println(observe)
				} else {
					output, err := json.MarshalIndent(simpleResults, "", "  ")
					if err != nil {
						cmd.Println(observe)
					} else {
						cmd.Println(string(output))
					}
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewTencentNewsSkill()
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

// NewTencentNewsCmd 创建腾讯新闻子命令
func NewTencentNewsCmd() *cobra.Command {
	skill := NewTencentNewsSkill()
	return skill.GetCommand()
}
