package baidu_search

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// BaiduSearchSkill 百度搜索技能
type BaiduSearchSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*BaiduSearchSkill)(nil)

// NewBaiduSearchSkill 创建百度搜索技能实例
func NewBaiduSearchSkill() *BaiduSearchSkill {
	skill := &BaiduSearchSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *BaiduSearchSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
// https://console.bce.baidu.com/qianfan/tools/toolsCenter/57d4e765-8af5-4ec0-8f9b-47075ec349e0/detail
func (s *BaiduSearchSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "baidu-search",
		Title:       "百度搜索",
		Version:     "1.1.4",
		Author:      "dtapp",
		Description: "百度搜索服务基于百度搜索能力与AI技术，专为生成式AI提供全网实时信息检索服务，具备海量内容站点、高权威性，强时效性，提供全新的搜索体验。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *BaiduSearchSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加搜索子命令
	cmd.AddCommand(s.createSearchCmd())

	return cmd
}

// 命令名称常量
const (
	CmdSearch = "search"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdSearch, FieldType: SearchSimpleResult{}},
	}
}

// requireAPIKey 检查并获取百度 API Key
func (s *BaiduSearchSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetBaiduAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置百度 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.baidu <your-api-key>' 设置，或使用 '%s config add api_keys.baidu <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createSearchCmd 创建搜索子命令
func (s *BaiduSearchSkill) createSearchCmd() *cobra.Command {
	var (
		query     string
		count     int
		freshness string
		verbose   bool
	)

	cmd := &cobra.Command{
		Use:   CmdSearch,
		Short: "执行百度搜索",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Printf("正在搜索: %s\n", query)

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 调用百度搜索 API
			opts := SearchOptions{
				Query:     query,
				Count:     count,
				Freshness: freshness,
			}

			resp, err := GetBaiduSearch(opts)
			if err != nil {
				cmd.Printf("搜索失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 References）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.References, "", "  "); err != nil {
					cmd.Println(string(resp.References))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果（使用结构体解析）
				if len(resp.References) > 0 {
					var simpleResults []SearchSimpleResult
					if err := json.Unmarshal(resp.References, &simpleResults); err != nil {
						// 如果解析失败，直接输出原始数据
						cmd.Println(string(resp.References))
					} else {
						output, err := json.MarshalIndent(simpleResults, "", "  ")
						if err != nil {
							cmd.Println(string(resp.References))
						} else {
							cmd.Println(string(output))
						}
					}
				} else {
					cmd.Println("[]")
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&query, "query", "q", "", "搜索关键词")
	cmd.Flags().IntVarP(&count, "count", "c", 10, "返回结果数量")
	cmd.Flags().StringVarP(&freshness, "freshness", "f", "", "时间筛选（如：day, week, month）")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewBaiduSearchSkill()
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

// NewBaiduSearchCmd 创建百度搜索子命令
func NewBaiduSearchCmd() *cobra.Command {
	skill := NewBaiduSearchSkill()
	return skill.GetCommand()
}
