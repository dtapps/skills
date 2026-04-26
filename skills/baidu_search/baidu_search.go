package baidu_search

import (
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
func (s *BaiduSearchSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "baidu-search",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "通过百度搜索 API 进行网络搜索",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义（只定义一次）
func (s *BaiduSearchSkill) createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "baidu-search",
		Short: "通过百度搜索 API 进行网络搜索",
		Long:  fmt.Sprintf("通过百度搜索 API 进行网络搜索 v%s\n作者: %s\n主页: %s", "1.0.0", "dtapp", "https://cnb.cool/dtapp/skills"),
	}

	// 添加搜索子命令
	cmd.AddCommand(s.createSearchCmd())
	cmd.AddCommand(s.createVersionCmd())

	return cmd
}

// createSearchCmd 创建搜索子命令
func (s *BaiduSearchSkill) createSearchCmd() *cobra.Command {
	var (
		query     string
		count     int
		freshness string
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "执行百度搜索",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("正在搜索: %s\n", query)

			// 获取百度 API Key
			apiKey := config.GetBaiduAPIKey()
			if apiKey == "" {
				cmd.Println("错误：未设置百度 API Key，请在配置文件或环境变量中设置")
				cmd.Println("提示：使用 'dtskills config set api_keys.baidu <your-api-key>' 设置")
				return
			}

			// 调用百度搜索 API
			opts := SearchOptions{
				Query:     query,
				Count:     count,
				Freshness: freshness,
			}

			resp, err := Search(apiKey, opts)
			if err != nil {
				cmd.Printf("搜索失败: %v\n", err)
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
	cmd.Flags().StringVarP(&query, "query", "q", "", "搜索关键词")
	cmd.Flags().IntVarP(&count, "count", "c", 10, "返回结果数量")
	cmd.Flags().StringVarP(&freshness, "freshness", "f", "", "时间筛选（如：day, week, month）")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("query")

	return cmd
}

// createVersionCmd 创建版本信息命令
func (s *BaiduSearchSkill) createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("技能名称: baidu-search\n")
			cmd.Printf("版本号: 1.0.0\n")
			cmd.Printf("作者: dtapp\n")
			cmd.Printf("描述: 通过百度搜索 API 进行网络搜索\n")
			cmd.Printf("主页: https://cnb.cool/dtapp/skills\n")
		},
	}
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
// 这是为了兼容现有代码的辅助函数
func GetSkillInfo() types.SkillInfo {
	skill := NewBaiduSearchSkill()
	return skills.ExtractSkillInfo(skill)
}

// NewBaiduSearchCmd 创建百度搜索子命令（为了兼容现有代码）
func NewBaiduSearchCmd() *cobra.Command {
	skill := NewBaiduSearchSkill()
	return skill.GetCommand()
}
