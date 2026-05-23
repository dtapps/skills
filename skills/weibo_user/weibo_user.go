package weibo_user

import (
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// WeiboUserSkill 微博用户技能
type WeiboUserSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*WeiboUserSkill)(nil)

// NewWeiboUserSkill 创建微博用户技能实例
func NewWeiboUserSkill() *WeiboUserSkill {
	skill := &WeiboUserSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *WeiboUserSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *WeiboUserSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "weibo-user",
		Title:       "微博用户",
		Version:     "1.0.1",
		Author:      "dtapp",
		Description: "微博用户技能，支持查询当前授权用户发布的微博列表和用户信息。需要 OAuth2 授权。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// createCommand 创建命令定义
func (s *WeiboUserSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createMyPostsCmd())
	cmd.AddCommand(s.createMyInfoCmd())

	return cmd
}

// 命令名称常量
const (
	CmdMyPosts = "my-posts"
	CmdMyInfo  = "my-info"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdMyPosts, FieldType: nil},
		{CommandName: CmdMyInfo, FieldType: nil},
	}
}

// requireAPIToken 检查并获取微博 API Token
func (s *WeiboUserSkill) requireAPIToken(cmd *cobra.Command) (string, bool) {
	apiToken := config.GetWeiboAPITokenRandom()
	if apiToken == "" {
		cmd.Println("错误：未设置微博 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_tokens.weibo <your-api-token>' 设置，或使用 '%s config add api_tokens.weibo <your-api-token>' 添加多个 Token\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiToken, true
}

// createMyPostsCmd 创建我的微博列表子命令
func (s *WeiboUserSkill) createMyPostsCmd() *cobra.Command {
	var (
		count  int
		page   int
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdMyPosts,
		Short: "获取当前用户发布的微博列表",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Token
			if _, ok := s.requireAPIToken(cmd); !ok {
				return
			}

			opts := UserTimelineOptions{
				Count:      count,
				Page:       page,
				TrimUser:   1,
				TrimStatus: 0,
			}

			resp, err := GetUserTimeline(opts)
			if err != nil {
				cmd.Printf("获取微博列表失败: %v\n", err)
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().IntVar(&count, "count", 20, "返回结果数量")
	cmd.Flags().IntVar(&page, "page", 1, "页码")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	return cmd
}

// createMyInfoCmd 创建我的信息子命令
func (s *WeiboUserSkill) createMyInfoCmd() *cobra.Command {
	var minify bool

	cmd := &cobra.Command{
		Use:   CmdMyInfo,
		Short: "获取当前授权用户信息",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Token
			if _, ok := s.requireAPIToken(cmd); !ok {
				return
			}

			resp, err := GetUserInfo()
			if err != nil {
				cmd.Printf("获取用户信息失败: %v\n", err)
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewWeiboUserSkill()
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

// NewWeiboUserCmd 创建微博用户子命令
func NewWeiboUserCmd() *cobra.Command {
	skill := NewWeiboUserSkill()
	return skill.GetCommand()
}
