package baidu_baike

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// BaiduBaikeSkill 百度百科技能
type BaiduBaikeSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*BaiduBaikeSkill)(nil)

// NewBaiduBaikeSkill 创建百度百科技能实例
func NewBaiduBaikeSkill() *BaiduBaikeSkill {
	skill := &BaiduBaikeSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *BaiduBaikeSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
// https://console.bce.baidu.com/qianfan/tools/toolsCenter/eb39b6d6-4ce4-498a-91f1-831a8dd38cf9/detail
func (s *BaiduBaikeSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "baidu-baike",
		Title:       "百度百科",
		Version:     "1.1.0",
		Author:      "dtapp",
		Description: "百度百科提供百科官方内容检索服务，帮助用户快速、全面地获取权威的百科信息。支持百科词条义项查询、词条详情查询以及秒懂百科视频获取。适用于AI问答、内容创作、教育教学、知识检索等场景，为模型提供可信的知识支撑。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *BaiduBaikeSkill) createCommand() *cobra.Command {
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
		{CommandName: CmdSearch, FieldType: BaikeSimpleResult{}},
	}
}

// requireAPIKey 检查并获取百度 API Key
func (s *BaiduBaikeSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetBaiduAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置百度 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.baidu <your-api-key>' 设置，或使用 '%s config add api_keys.baidu <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createSearchCmd 创建搜索子命令
func (s *BaiduBaikeSkill) createSearchCmd() *cobra.Command {
	var (
		searchType string
		searchKey  string
		verbose    bool
	)

	cmd := &cobra.Command{
		Use:   CmdSearch,
		Short: "查询百度百科词条",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("正在查询百度百科：%s\n", searchKey)

			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 调用百度百科 API
			opts := BaikeOptions{
				SearchType: searchType,
				SearchKey:  searchKey,
			}

			resp, err := GetBaiduBaike(opts)
			if err != nil {
				cmd.Printf("查询失败: %v\n", err)
				return
			}

			// 输出结果
			if verbose {
				// 输出完整结果（格式化后的 Result）
				var buf bytes.Buffer
				if err := json.Indent(&buf, resp.Result, "", "  "); err != nil {
					cmd.Println(string(resp.Result))
				} else {
					cmd.Println(buf.String())
				}
			} else {
				// 输出简洁结果（直接使用结构体解析）
				if len(resp.Result) > 0 {
					var simpleResult BaikeSimpleResult
					if err := json.Unmarshal(resp.Result, &simpleResult); err != nil {
						cmd.Println(string(resp.Result))
					} else {
						output, err := json.MarshalIndent(simpleResult, "", "  ")
						if err != nil {
							cmd.Println(string(resp.Result))
						} else {
							cmd.Println(string(output))
						}
					}
				} else {
					cmd.Println("{}")
				}
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&searchType, "type", "t", "lemmaTitle", "检索类型：lemmaTitle(按词条名)、lemmaId(按词条ID)")
	cmd.Flags().StringVarP(&searchKey, "key", "k", "", "检索关键字")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewBaiduBaikeSkill()
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

// NewBaiduBaikeCmd 创建百度百科子命令
func NewBaiduBaikeCmd() *cobra.Command {
	skill := NewBaiduBaikeSkill()
	return skill.GetCommand()
}
