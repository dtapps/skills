package hefeng_weather

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// HefengWeatherSkill 和风天气技能
type HefengWeatherSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*HefengWeatherSkill)(nil)

// NewHefengWeatherSkill 创建和风天气技能实例
func NewHefengWeatherSkill() *HefengWeatherSkill {
	skill := &HefengWeatherSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *HefengWeatherSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *HefengWeatherSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "hefeng-weather",
		Title:       "和风天气",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "和风天气信息查询工具，覆盖中国市级和区县级行政区。当用户查询实况天气、每日天气预报、实时空气质量、天气指数预报信息时使用。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdCurrent       = "current"
	CmdForecast      = "forecast"
	CmdIndices       = "indices"
	CmdSunriseSunset = "sunrise-sunset"
	CmdMoonPhase     = "moon-phase"
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
		{CommandName: CmdIndices, FieldType: nil},
		{CommandName: CmdSunriseSunset, FieldType: nil},
		{CommandName: CmdMoonPhase, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *HefengWeatherSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createCurrentCmd())
	cmd.AddCommand(s.createForecastCmd())
	cmd.AddCommand(s.createIndicesCmd())
	cmd.AddCommand(s.createSunriseSunsetCmd())
	cmd.AddCommand(s.createMoonPhaseCmd())

	return cmd
}

// requireAPIKey 检查并获取和风天气 API Key
func (s *HefengWeatherSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetHefengAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置和风天气 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.hefeng <your-api-key>' 设置，或使用 '%s config add api_keys.hefeng <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createCurrentCmd 创建实时天气子命令
func (s *HefengWeatherSkill) createCurrentCmd() *cobra.Command {
	var (
		location string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdCurrent,
		Short: "获取实时天气",
		Long:  "获取中国3000+市县区和海外20万个城市实时天气数据，包括实时温度、体感温度、风力风向、相对湿度、大气压强、降水量、能见度、露点温度、云量等。",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetCurrentWeather(location)
			if err != nil {
				cmd.Printf("获取实时天气失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				cmd.Println(string(resp))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&location, "location", "", "位置（经度,纬度 或 地区ID，如：116.41,39.92 或 101010100）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

// createForecastCmd 创建每日天气预报子命令
func (s *HefengWeatherSkill) createForecastCmd() *cobra.Command {
	var (
		location string
		days     string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdForecast,
		Short: "获取每日天气预报",
		Long:  "每日天气预报API，提供全球城市未来3-30天天气预报，包括：日出日落、月升月落、最高最低温度、天气白天和夜间状况、风力、风速、风向、相对湿度、大气压强、降水量、露点温度、紫外线强度、能见度等。",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetDailyForecast(location, days)
			if err != nil {
				cmd.Printf("获取天气预报失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				cmd.Println(string(resp))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&location, "location", "", "位置（经度,纬度 或 地区ID，如：116.41,39.92 或 101010100）")
	cmd.Flags().StringVar(&days, "days", "7d", "预报天数: 3d, 7d, 10d, 15d, 30d")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

// createIndicesCmd 创建天气指数预报子命令
func (s *HefengWeatherSkill) createIndicesCmd() *cobra.Command {
	var (
		location string
		days     string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdIndices,
		Short: "获取天气指数预报",
		Long:  "获取中国和全球城市天气生活指数预报数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetIndicesForecast(location, days)
			if err != nil {
				cmd.Printf("获取天气指数预报失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				cmd.Println(string(resp))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&location, "location", "", "位置（经度,纬度 或 地区ID，如：116.41,39.92 或 101010100）")
	cmd.Flags().StringVar(&days, "days", "3d", "预报天数: 1d, 3d")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

// createSunriseSunsetCmd 创建日出日落子命令
func (s *HefengWeatherSkill) createSunriseSunsetCmd() *cobra.Command {
	var (
		location string
		date     string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdSunriseSunset,
		Short: "获取日出日落",
		Long:  "获取未来60天全球任意地点日出日落时间",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			if date == "" {
				date = time.Now().Format("20060102")
			}

			resp, err := GetSunriseSunset(location, date)
			if err != nil {
				cmd.Printf("获取日出日落时间失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				cmd.Println(string(resp))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&location, "location", "", "位置（经度,纬度 或 地区ID，如：116.41,39.92 或 101010100）")
	cmd.Flags().StringVar(&date, "date", "", "日期，格式为 yyyyMMdd，如 20240101")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

// createMoonPhaseCmd 创建月升月落和月相子命令
func (s *HefengWeatherSkill) createMoonPhaseCmd() *cobra.Command {
	var (
		location string
		date     string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdMoonPhase,
		Short: "获取月升月落和月相",
		Long:  "获取未来60天全球城市月升月落和逐小时的月相数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			if date == "" {
				date = time.Now().Format("20060102")
			}

			resp, err := GetMoonPhase(location, date)
			if err != nil {
				cmd.Printf("获取月升月落和月相信息失败: %v\n", err)
				return
			}

			if minify {
				// 紧凑输出
				cmd.Println(string(resp))
			} else {
				// 美化输出
				output, _ := json.MarshalIndent(resp, "", "  ")
				cmd.Println(string(output))
			}
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&location, "location", "", "位置（经度,纬度 或 地区ID，如：116.41,39.92 或 101010100）")
	cmd.Flags().StringVar(&date, "date", "", "日期，格式为 yyyyMMdd，如 20240101")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewHefengWeatherSkill()
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

// NewHefengWeatherCmd 创建和风天气命令
func NewHefengWeatherCmd() *cobra.Command {
	skill := NewHefengWeatherSkill()
	return skill.GetCommand()
}
