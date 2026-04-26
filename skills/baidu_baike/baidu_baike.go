package baidu_baike

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// SkillInfo 技能配置信息
var SkillInfo = types.SkillInfo{
	Name:        "baidu-baike",
	Version:     "1.0.0",
	Author:      "dtapp",
	Description: "通过百度百科 API 查询百科内容",
	Homepage:    "https://cnb.cool/dtapp/skills",
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	info := SkillInfo
	info.Commands = extractCommands()
	return info
}

// extractCommands 从 cobra.Command 提取 Commands 信息
func extractCommands() []types.Command {
	cmd := NewBaiduBaikeCmd()
	var commands []types.Command

	for _, subCmd := range cmd.Commands() {
		command := types.Command{
			Name:        subCmd.Name(),
			Description: subCmd.Short,
			Usage:       fmt.Sprintf("%s %s %s", types.ProjectName, cmd.Name(), subCmd.Use),
			Flags:       extractFlags(subCmd),
		}
		commands = append(commands, command)
	}

	return commands
}

// extractFlags 从 cobra.Command 提取 Flags 信息
func extractFlags(cmd *cobra.Command) []types.Flag {
	var flags []types.Flag

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flag := types.Flag{
			Name:        f.Name,
			Shorthand:   f.Shorthand,
			Description: f.Usage,
			Default:     f.DefValue,
		}
		// 检查是否是必填参数 (通过 Annotations 判断)
		if f.Annotations != nil {
			if _, ok := f.Annotations["cobra_annotation_bash_completion_one_required_flag"]; ok {
				flag.Required = true
			}
		}
		flags = append(flags, flag)
	})

	return flags
}

// NewBaiduBaikeCmd 创建百度百科子命令
func NewBaiduBaikeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   SkillInfo.Name,
		Short: SkillInfo.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", SkillInfo.Description, SkillInfo.Version, SkillInfo.Author, SkillInfo.Homepage),
	}

	// 添加搜索子命令
	cmd.AddCommand(NewSearchCmd())
	cmd.AddCommand(NewVersionCmd())

	return cmd
}

// NewVersionCmd 创建版本信息命令
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("技能名称: %s\n", SkillInfo.Name)
			cmd.Printf("版本号: %s\n", SkillInfo.Version)
			cmd.Printf("作者: %s\n", SkillInfo.Author)
			cmd.Printf("描述: %s\n", SkillInfo.Description)
			cmd.Printf("主页: %s\n", SkillInfo.Homepage)
		},
	}
}

// NewSearchCmd 创建百度百科搜索命令
func NewSearchCmd() *cobra.Command {
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
