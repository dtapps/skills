package uapi

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// UapiSkill UAPI技能
type UapiSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*UapiSkill)(nil)

// NewUapiSkill 创建UAPI技能实例
func NewUapiSkill() *UapiSkill {
	skill := &UapiSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *UapiSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *UapiSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "uapi",
		Title:       "UAPI",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "UAPI工具，提供天气服务。",
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
func (s *UapiSkill) createCommand() *cobra.Command {
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
func (s *UapiSkill) createWeatherCmd() *cobra.Command {
	var (
		city   string
		adcode string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdWeather,
		Short: "获取天气信息",
		Long:  "出门前，查一下天气总是个好习惯。这个接口为你提供精准、实时的天气数据，支持国内和国际城市。",
		Run: func(cmd *cobra.Command, args []string) {

			// 必须提供 city 或 adcode 之一
			if city == "" && adcode == "" {
				cmd.Println("错误：必须提供 --city 或 --adcode 参数之一")
				return
			}

			resp, err := GetWeather(city, adcode)
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
	cmd.Flags().StringVar(&city, "city", "", "城市名称，支持中文（北京）和英文（Tokyo）。可选参数，不传时会尝试 IP 自动定位。")
	cmd.Flags().StringVar(&adcode, "adcode", "", "城市行政区划代码（如 110000），优先级高于 city。可选参数，不传时会尝试 IP 自动定位。")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewUapiSkill()
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

// NewUapiCmd 创建UAPI命令
func NewUapiCmd() *cobra.Command {
	skill := NewUapiSkill()
	return skill.GetCommand()
}
