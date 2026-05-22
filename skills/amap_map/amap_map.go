package amap_map

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// AmapMapSkill 高德地图技能
type AmapMapSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*AmapMapSkill)(nil)

// NewAmapMapSkill 创建高德地图技能实例
func NewAmapMapSkill() *AmapMapSkill {
	skill := &AmapMapSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *AmapMapSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *AmapMapSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "amap-map",
		Title:       "高德地图",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "高德地图工具，支持实时天气和天气预报查询。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdWeather    = "weather"
	CmdIpLocation = "ip-location"
)

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewAmapMapSkill()
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

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdWeather, FieldType: WeatherSimpleResult{}},
		{CommandName: CmdIpLocation, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *AmapMapSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createWeatherCmd())
	cmd.AddCommand(s.createIpLocationCmd())

	return cmd
}

// requireAPIKey 检查并获取高德地图 API Key
func (s *AmapMapSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetAmapMapAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置高德地图 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.amap_map <your-api-key>' 设置，或使用 '%s config add api_keys.amap_map <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createWeatherCmd 创建实时天气子命令
func (s *AmapMapSkill) createWeatherCmd() *cobra.Command {
	var (
		city   string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdWeather,
		Short: "获取天气信息",
		Long:  "获取指定位置的实时天气状况",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetWeather(city, s.GetInfo().Name)
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
	cmd.Flags().StringVar(&city, "city", "", "城市名称或者adcode")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("city")

	return cmd
}

// createIpLocationCmd 创建IP定位子命令
func (s *AmapMapSkill) createIpLocationCmd() *cobra.Command {
	var (
		ip     string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdIpLocation,
		Short: "获取IP定位信息",
		Long:  "获取指定IP的定位信息",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetIpLocation(ip, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取IP定位失败: %v\n", err)
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
	cmd.Flags().StringVar(&ip, "ip", "", "IP地址")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("ip")

	return cmd
}
