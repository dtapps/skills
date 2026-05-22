package tencent_map

import (
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// TencentMapSkill 腾讯地图技能
type TencentMapSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*TencentMapSkill)(nil)

// NewTencentMapSkill 创建腾讯地图技能实例
func NewTencentMapSkill() *TencentMapSkill {
	skill := &TencentMapSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *TencentMapSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *TencentMapSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "tencent-map",
		Title:       "腾讯地图",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "腾讯地图工具，支持实时天气和天气预报查询。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdCurrentWeather  = "current-weather"
	CmdForecastWeather = "forecast-weather"
	CmdIpLocation      = "ip-location"
)

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewTencentMapSkill()
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
		{CommandName: CmdCurrentWeather, FieldType: nil},
		{CommandName: CmdForecastWeather, FieldType: nil},
		{CommandName: CmdIpLocation, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *TencentMapSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createCurrentWeatherCmd())
	cmd.AddCommand(s.createForecastWeatherCmd())
	cmd.AddCommand(s.createIpLocationCmd())

	return cmd
}

// requireAPIKey 检查并获取腾讯地图 API Key
func (s *TencentMapSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetTencentMapAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置腾讯地图 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.tencent_map <your-api-key>' 设置，或使用 '%s config add api_keys.tencent_map <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createCurrentWeatherCmd 创建当前天气子命令
func (s *TencentMapSkill) createCurrentWeatherCmd() *cobra.Command {
	var (
		adcode   string
		location string
	)

	cmd := &cobra.Command{
		Use:   CmdCurrentWeather,
		Short: "获取实时天气",
		Long:  "获取指定位置的实时天气状况",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 必须提供 adcode 或 location 之一
			if adcode == "" && location == "" {
				cmd.Println("错误：必须提供 --adcode 或 --location 参数之一")
				return
			}

			resp, err := GetCurrentWeather(adcode, location, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气失败: %v\n", err)
				return
			}

			// 语义化输出
			cmd.Println(*resp)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&adcode, "adcode", "", "行政区划代码，与location经纬度坐标必有其一（可通过调用地址解析获取，输入地址或城市名，返回结果中的adcode即为行政区划代码）")
	cmd.Flags().StringVar(&location, "location", "", "经纬度坐标，与adcode行政区划代码必有其一，结构是lat,lng")

	return cmd
}

// createForecastWeatherCmd 创建天气预报子命令
func (s *TencentMapSkill) createForecastWeatherCmd() *cobra.Command {
	var (
		adcode   string
		location string
	)

	cmd := &cobra.Command{
		Use:   CmdForecastWeather,
		Short: "获取天气预报",
		Long:  "获取指定位置的未来几天天气预报",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 必须提供 adcode 或 location 之一
			if adcode == "" && location == "" {
				cmd.Println("错误：必须提供 --adcode 或 --location 参数之一")
				return
			}

			resp, err := GetForecastWeather(adcode, location, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气预报失败: %v\n", err)
				return
			}

			// 语义化输出
			cmd.Println(*resp)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&adcode, "adcode", "", "行政区划代码，与location经纬度坐标必有其一（可通过调用地址解析获取，输入地址或城市名，返回结果中的adcode即为行政区划代码）")
	cmd.Flags().StringVar(&location, "location", "", "经纬度坐标，与adcode行政区划代码必有其一，结构是lat,lng")

	return cmd
}

// createIpLocationCmd 创建IP定位子命令
func (s *TencentMapSkill) createIpLocationCmd() *cobra.Command {
	var (
		ip string
	)

	cmd := &cobra.Command{
		Use:   CmdIpLocation,
		Short: "获取IP定位",
		Long:  "获取指定IP地址的定位信息",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 调用腾讯地图 API 获取 IP定位
			resp, err := GetIpLocation(ip, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取IP定位失败: %v\n", err)
				return
			}

			// 语义化输出
			cmd.Println(*resp)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&ip, "ip", "", "IP地址(IPv4或IPv6）")

	return cmd
}
