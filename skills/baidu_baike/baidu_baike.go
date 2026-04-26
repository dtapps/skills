package baidu_baike

import (
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
func (s *BaiduBaikeSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "baidu-baike",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "通过百度百科 API 查询百科内容",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *BaiduBaikeSkill) createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "baidu-baike",
		Short: "通过百度百科 API 查询百科内容",
		Long:  fmt.Sprintf("通过百度百科 API 查询百科内容 v%s\n作者: %s\n主页: %s", "1.0.0", "dtapp", "https://cnb.cool/dtapp/skills"),
	}

	// 添加搜索子命令
	cmd.AddCommand(s.createSearchCmd())
	cmd.AddCommand(s.createVersionCmd())

	return cmd
}

// createSearchCmd 创建搜索子命令
func (s *BaiduBaikeSkill) createSearchCmd() *cobra.Command {
	var (
		searchType string
		searchKey  string
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "查询百度百科词条",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("正在查询百度百科：%s\n", searchKey)

			// 获取百度 API Key
			apiKey := config.GetBaiduAPIKey()
			if apiKey == "" {
				cmd.Println("错误：未设置百度 API Key，请在配置文件或环境变量中设置")
				cmd.Println("提示：使用 'dtskills config set api_keys.baidu <your-api-key>' 设置")
				return
			}

			// 调用百度百科 API
			opts := BaikeOptions{
				SearchType: searchType,
				SearchKey:  searchKey,
			}

			resp, err := Search(apiKey, opts)
			if err != nil {
				cmd.Printf("查询失败: %v\n", err)
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
	cmd.Flags().StringVarP(&searchType, "type", "t", "lemmaTitle", "检索类型：lemmaTitle(按词条名)、lemmaId(按词条ID)")
	cmd.Flags().StringVarP(&searchKey, "key", "k", "", "检索关键字")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

// createVersionCmd 创建版本信息命令
func (s *BaiduBaikeSkill) createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("技能名称: baidu-baike\n")
			cmd.Printf("版本号: 1.0.0\n")
			cmd.Printf("作者: dtapp\n")
			cmd.Printf("描述: 通过百度百科 API 查询百科内容\n")
			cmd.Printf("主页: https://cnb.cool/dtapp/skills\n")
		},
	}
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
// 这是为了兼容现有代码的辅助函数
func GetSkillInfo() types.SkillInfo {
	skill := NewBaiduBaikeSkill()
	return skills.ExtractSkillInfo(skill)
}

// NewBaiduBaikeCmd 创建百度百科子命令（为了兼容现有代码）
func NewBaiduBaikeCmd() *cobra.Command {
	skill := NewBaiduBaikeSkill()
	return skill.GetCommand()
}
