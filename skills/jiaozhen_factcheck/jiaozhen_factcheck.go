package jiaozhen_factcheck

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// JiaozhenFactCheckSkill 较真事实查证技能
type JiaozhenFactCheckSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*JiaozhenFactCheckSkill)(nil)

// NewJiaozhenFactCheckSkill 创建较真事实查证技能实例
func NewJiaozhenFactCheckSkill() *JiaozhenFactCheckSkill {
	skill := &JiaozhenFactCheckSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *JiaozhenFactCheckSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
// https://skillhub.cn/skills/jiaozhen-factcheck
func (s *JiaozhenFactCheckSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "jiaozhen-factcheck",
		Title:       "腾讯较真",
		Version:     "1.0.1",
		Author:      "TencentNews",
		Description: "事实查证工具，对输入内容的具体说法、资讯、事件或常识进行真实性、准确性、可靠性判断。当用户需要较真一下，查证问题或判断信息真伪、识别谣言、询问真假，是真的吗，真的假的，能否xxx，可不可以，是谣言吗...等场景时调用。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *JiaozhenFactCheckSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createCheckCmd())

	return cmd
}

// requireAPIKey 检查并获取腾讯较真 API Key
func (s *JiaozhenFactCheckSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetTencentNewsAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置腾讯较真 API Key，请在配置文件或环境变量中设置")
		cmd.Println("提示：使用 'dtskills config set api_keys.tencent_news <your-api-key>' 设置，或使用 'dtskills config add api_keys.tencent_news <your-api-key>' 添加多个 Key")
		return "", false
	}
	return apiKey, true
}

// createCheckCmd 创建查证子命令
func (s *JiaozhenFactCheckSkill) createCheckCmd() *cobra.Command {
	var (
		query   string
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "事实查证",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			cmd.Printf("正在查证: %s\n", query)

			resp, err := FactCheck(query, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("查证失败: %v\n", err)
				return
			}

			// 输出结果
			output, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				cmd.Printf("格式化结果失败: %v\n", err)
				return
			}

			cmd.Println(string(output))
		},
	}

	// 定义 flags
	cmd.Flags().StringVarP(&query, "query", "q", "", "需要查证的内容")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "输出完整结果")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewJiaozhenFactCheckSkill()
	return skills.ExtractSkillInfo(skill)
}

// NewJiaozhenFactCheckCmd 创建较真事实查证命令
func NewJiaozhenFactCheckCmd() *cobra.Command {
	skill := NewJiaozhenFactCheckSkill()
	return skill.GetCommand()
}
