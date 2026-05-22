package baidu_map

import (
	"encoding/json"
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"

	"github.com/spf13/cobra"
)

// BaiduMapSkill 百度地图技能
type BaiduMapSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*BaiduMapSkill)(nil)

// NewBaiduMapSkill 创建百度地图技能实例
func NewBaiduMapSkill() *BaiduMapSkill {
	skill := &BaiduMapSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *BaiduMapSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *BaiduMapSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "baidu-map",
		Title:       "百度地图",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "百度地图工具，支持实时天气和天气预报查询。",
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
	skill := NewBaiduMapSkill()
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
		{CommandName: CmdWeather, FieldType: nil},
		{CommandName: CmdIpLocation, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *BaiduMapSkill) createCommand() *cobra.Command {
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

// requireAPIKey 检查并获取百度地图 API Key
func (s *BaiduMapSkill) requireAPIKey(cmd *cobra.Command) (string, bool) {
	apiKey := config.GetBaiduMapAPIKeyRandom()
	if apiKey == "" {
		cmd.Println("错误：未设置百度地图 API Key，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set api_keys.baidu_map <your-api-key>' 设置，或使用 '%s config add api_keys.baidu_map <your-api-key>' 添加多个 Key\n", types.ProjectName, types.ProjectName)
		return "", false
	}
	return apiKey, true
}

// createWeatherCmd 创建获取天气子命令
func (s *BaiduMapSkill) createWeatherCmd() *cobra.Command {
	var (
		location   string
		districtID string
		isChina    bool
		minify     bool
	)

	cmd := &cobra.Command{
		Use:   CmdWeather,
		Short: "获取天气信息",
		Long:  "获取指定位置的天气信息，支持实时天气和未来5天预报",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			// 必须提供 district_id 或 location 之一
			if districtID == "" && location == "" {
				cmd.Println("错误：必须提供 --district_id 或 --location 参数之一")
				return
			}

			resp, err := GetWeather(districtID, location, isChina, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取天气失败: %v\n", err)
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
	cmd.Flags().StringVar(&location, "location", "", "经纬度坐标, 经度在前纬度在后, 逗号分隔")
	cmd.Flags().StringVar(&districtID, "district_id", "", "行政区划代码, 需保证为6位无符号整数")
	cmd.Flags().BoolVar(&isChina, "is_china", true, "查询地是否在中国境内(包括香港、澳门、台湾地区), 可选值为`true`或`false`, 默认为`true`; 若请求其他国家或地区, 则需传`false`")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	return cmd
}

// createIpLocationCmd 创建IP定位子命令
func (s *BaiduMapSkill) createIpLocationCmd() *cobra.Command {
	var (
		ip     string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdIpLocation,
		Short: "获取IP定位信息",
		Long:  "获取指定IP的定位信息",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查并获取 API Key
			if _, ok := s.requireAPIKey(cmd); !ok {
				return
			}

			resp, err := GetIpLocation(ip, s.GetInfo().Name)
			if err != nil {
				cmd.Printf("获取定位信息失败: %v\n", err)
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
	cmd.Flags().StringVar(&ip, "ip", "", "需要定位的IP地址, 如果为空则获取用户当前IP地址(支持IPv4和IPv6)")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	return cmd
}
