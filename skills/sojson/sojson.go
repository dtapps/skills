package sojson

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// SojsonSkill SoJSON技能
type SojsonSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*SojsonSkill)(nil)

// NewSojsonSkill 创建SoJSON技能实例
func NewSojsonSkill() *SojsonSkill {
	skill := &SojsonSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *SojsonSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *SojsonSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "sojson",
		Title:       "SoJSON",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "SoJSON工具，提供天气服务。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdWeather = "weather"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdWeather, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *SojsonSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createWeatherCmd())

	return cmd
}

// createWeatherCmd 创建获取天气子命令
func (s *SojsonSkill) createWeatherCmd() *cobra.Command {
	var (
		city   string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdWeather,
		Short: "获取天气信息",
		Long:  "获取指定城市的实时天气和未来5天预报",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := GetWeather(city)
			if err != nil {
				cmd.Printf("获取天气失败: %v\n", err)
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
	cmd.Flags().StringVar(&city, "city", "", "城市编码（如：101010100表示北京）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("city")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewSojsonSkill()
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

// NewSojsonCmd 创建SoJSON命令
func NewSojsonCmd() *cobra.Command {
	skill := NewSojsonSkill()
	return skill.GetCommand()
}
