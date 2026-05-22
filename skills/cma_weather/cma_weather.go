package cma_weather

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// CmaWeatherSkill 中国气象局天气技能
type CmaWeatherSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*CmaWeatherSkill)(nil)

// NewCmaWeatherSkill 创建中国气象局天气技能实例
func NewCmaWeatherSkill() *CmaWeatherSkill {
	skill := &CmaWeatherSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *CmaWeatherSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *CmaWeatherSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "cma-weather",
		Title:       "中国气象局天气",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "中国气象局天气查询工具，提供全国城市实时天气信息。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdNow = "now"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdNow, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *CmaWeatherSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createNowCmd())

	return cmd
}

// createNowCmd 创建实时天气子命令
func (s *CmaWeatherSkill) createNowCmd() *cobra.Command {
	var (
		cityID string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdNow,
		Short: "获取实时天气",
		Long:  "获取指定城市的实时天气状况",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := GetCurrentWeather(cityID)
			if err != nil {
				cmd.Printf("获取实时天气失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				output, _ := json.Marshal(resp)
				cmd.Println(string(output))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&cityID, "city", "", "城市ID（如：54511表示北京）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("city")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewCmaWeatherSkill()
	info := skills.ExtractSkillInfo(skill)

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

// NewCmaWeatherCmd 创建中国气象局天气命令
func NewCmaWeatherCmd() *cobra.Command {
	skill := NewCmaWeatherSkill()
	return skill.GetCommand()
}
