package tencent_weather

import (
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// TencentWeatherSkill 腾讯天气技能
type TencentWeatherSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*TencentWeatherSkill)(nil)

// NewTencentWeatherSkill 创建腾讯天气技能实例
func NewTencentWeatherSkill() *TencentWeatherSkill {
	skill := &TencentWeatherSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *TencentWeatherSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
// https://skillhub.cn/skills/tencent-weather
func (s *TencentWeatherSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "tencent-weather",
		Title:       "腾讯天气",
		Version:     "1.0.0",
		Author:      "TencentNews",
		Description: "天气信息查询工具，覆盖中国市级和区县级行政区。当用户查询实况天气、天气预报信息时使用。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdCurrent    = "current"
	CmdForecast   = "forecast"
	CmdAirQuality = "air-quality"
	CmdIP2City    = "ip2city"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdCurrent, FieldType: nil},
		{CommandName: CmdForecast, FieldType: nil},
		{CommandName: CmdAirQuality, FieldType: nil},
		{CommandName: CmdIP2City, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *TencentWeatherSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createCurrentCmd())
	cmd.AddCommand(s.createForecastCmd())
	cmd.AddCommand(s.createAirQualityCmd())
	cmd.AddCommand(s.createIP2CityCmd())

	return cmd
}

// requireAPIKey 检查并获取腾讯天气 API Key
func (s *TencentWeatherSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetTencentNewsAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置腾讯天气 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.tencent_news <your-api-key>' 设置，或使用 '%s config add api_keys.tencent_news <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createCurrentCmd 创建实时天气子命令
func (s *TencentWeatherSkill) createCurrentCmd() *cobra.Command {
	var (
		adcode string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdCurrent,
		Short: "获取实时天气",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetCurrentWeather(adcode, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气失败: %v\n", err)
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&adcode, "adcode", "", "地区编码（如北京：110000）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("adcode")

	return cmd
}

// createForecastCmd 创建天气预报子命令
func (s *TencentWeatherSkill) createForecastCmd() *cobra.Command {
	var (
		adcode      string
		weatherType string
		minify      bool
	)

	cmd := &cobra.Command{
		Use:   CmdForecast,
		Short: "获取天气预报",
		Long:  "获取天气预报，支持 forecast_1h(1小时预报) 或 forecast_24h(24小时预报)",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetForecastWeather(adcode, weatherType, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气预报失败: %v\n", err)
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&adcode, "adcode", "", "地区编码（如北京：110000）")
	cmd.Flags().StringVar(&weatherType, "type", "forecast_24h", "预报类型: forecast_1h 或 forecast_24h")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("adcode")

	return cmd
}

// createAirQualityCmd 创建天气空气质量子命令
func (s *TencentWeatherSkill) createAirQualityCmd() *cobra.Command {
	var (
		adcode string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdAirQuality,
		Short: "获取天气空气质量",
		Long:  "获取天气空气质量",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetAirQuality(adcode, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气空气质量失败: %v\n", err)
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&adcode, "adcode", "", "地区编码（如北京：110000）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("adcode")

	return cmd
}

// createIP2CityCmd 创建IP定位子命令
func (s *TencentWeatherSkill) createIP2CityCmd() *cobra.Command {
	var (
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdIP2City,
		Short: "IP定位获取城市",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetIP2City(s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取IP定位失败: %v\n", err)
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
	skill := NewTencentWeatherSkill()
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

// NewTencentWeatherCmd 创建腾讯天气命令
func NewTencentWeatherCmd() *cobra.Command {
	skill := NewTencentWeatherSkill()
	return skill.GetCommand()
}
