package crawler

import (
	"fmt"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/skills"
	"github.com/dtapp/skills/types"
	"github.com/tidwall/gjson"

	"github.com/spf13/cobra"
)

// CrawlerSkill 爬虫数据技能
type CrawlerSkill struct {
	cmd *cobra.Command
}

// 确保实现了 Skill 接口
var _ skills.Skill = (*CrawlerSkill)(nil)

// NewCrawlerSkill 创建爬虫数据技能实例
func NewCrawlerSkill() *CrawlerSkill {
	skill := &CrawlerSkill{}
	skill.cmd = skill.createCommand()
	return skill
}

// GetCommand 获取技能的 cobra.Command
func (s *CrawlerSkill) GetCommand() *cobra.Command {
	return s.cmd
}

// GetInfo 获取技能元数据
func (s *CrawlerSkill) GetInfo() types.SkillInfo {
	return types.SkillInfo{
		Name:        "crawler",
		Title:       "爬虫数据",
		Version:     "1.0.0",
		Author:      "dtapp",
		Description: "爬虫数据查询工具，提供黄金价格、热搜、天气等数据。",
		Homepage:    "https://cnb.cool/dtapp/skills",
	}
}

// 命令名称常量
const (
	CmdGoldPriceData       = "gold-price-data"
	CmdLatestGoldPrice     = "latest-gold-price"
	CmdJDFinanceGoldStatus = "jd-finance-gold-status"
	CmdHotData             = "hot-data"
	CmdWeatherData         = "weather-data"
)

// CommandOutputField 命令与输出字段的映射
type CommandOutputField struct {
	CommandName string
	FieldType   any
}

// GetCommandOutputFields 获取所有命令的输出字段映射
func GetCommandOutputFields() []CommandOutputField {
	return []CommandOutputField{
		{CommandName: CmdGoldPriceData, FieldType: nil},
		{CommandName: CmdLatestGoldPrice, FieldType: nil},
		{CommandName: CmdHotData, FieldType: nil},
		{CommandName: CmdWeatherData, FieldType: nil},
	}
}

// createCommand 创建命令定义
func (s *CrawlerSkill) createCommand() *cobra.Command {
	info := s.GetInfo()
	cmd := &cobra.Command{
		Use:   info.Name,
		Short: info.Description,
		Long:  fmt.Sprintf("%s v%s\n作者: %s\n主页: %s", info.Description, info.Version, info.Author, info.Homepage),
	}

	// 添加子命令
	cmd.AddCommand(s.createGoldPriceDataCmd())
	cmd.AddCommand(s.createLatestGoldPriceCmd())
	cmd.AddCommand(s.createJDFinanceGoldStatusCmd())
	cmd.AddCommand(s.createHotDataCmd())
	cmd.AddCommand(s.createWeatherDataCmd())

	return cmd
}

// requireMCPURL 检查并获取 MCP URL
func (s *CrawlerSkill) requireMCPURL(cmd *cobra.Command) (string, bool) {
	mcpURL := config.GetCrawlerMCPURL()
	if mcpURL == "" {
		cmd.Println("错误：未设置 Crawler MCP URL，请在配置文件或环境变量中设置")
		cmd.Printf("提示：使用 '%s config set mcp.crawler.url <url>' 设置\n", types.ProjectName)
		cmd.Printf("或设置环境变量 %s\n", config.CrawlerMCPURLEnv)
		return "", false
	}
	return mcpURL, true
}

// createGoldPriceDataCmd 创建获取黄金价格全部数据子命令
func (s *CrawlerSkill) createGoldPriceDataCmd() *cobra.Command {
	var (
		date   string
		source string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdGoldPriceData,
		Short: "获取黄金价格全部数据",
		Long:  "获取指定日期黄金价格的全部数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireMCPURL(cmd); !ok {
				return
			}

			resp, err := GetGoldPriceData(date, source)
			if err != nil {
				cmd.Printf("获取黄金价格数据失败: %v\n", err)
				return
			}

			// 打印 JSON 响应
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&date, "date", "", "日期（格式：YYYY-MM-DD）")
	cmd.Flags().StringVar(&source, "source", "jd_jr", "数据源，可选值：jd_jr,cngold,eastmoney")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

// createLatestGoldPriceCmd 创建获取最新黄金价格子命令
func (s *CrawlerSkill) createLatestGoldPriceCmd() *cobra.Command {
	var (
		date   string
		source string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdLatestGoldPrice,
		Short: "获取最新黄金价格",
		Long:  "获取指定日期最新的一条黄金价格数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireMCPURL(cmd); !ok {
				return
			}

			resp, err := GetLatestGoldPrice(date, source)
			if err != nil {
				cmd.Printf("获取最新黄金价格失败: %v\n", err)
				return
			}

			// 打印 JSON 响应
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&date, "date", "", "日期（格式：YYYY-MM-DD）")
	cmd.Flags().StringVar(&source, "source", "jd_jr", "数据源，可选值：jd_jr,cngold,eastmoney")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

// createJDFinanceGoldStatusCmd 创建京东金融黄金市场状态子命令
func (s *CrawlerSkill) createJDFinanceGoldStatusCmd() *cobra.Command {
	var (
		datetime string
		format   string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdJDFinanceGoldStatus,
		Short: "获取京东金融黄金市场状态",
		Long:  "获取指定日期京东金融黄金市场状态，支持格式化输出",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireMCPURL(cmd); !ok {
				return
			}

			resp, err := GetJDFinanceGoldStatus(datetime)
			if err != nil {
				cmd.Printf("获取京东金融黄金市场状态失败: %v\n", err)
				return
			}

			// 如果指定了 format，使用 gjson 提取指定字段
			if format != "" {
				result := gjson.GetBytes(resp, format)
				if result.Exists() {
					cmd.Println(result.String())
				} else {
					cmd.Printf("字段 '%s' 不存在\n", format)
				}
				return
			}

			// 默认完整输出
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&datetime, "datetime", "", "日期时间（格式：YYYY-MM-DD HH:MM:SS）")
	cmd.Flags().StringVar(&format, "format", "", "gjson 路径，用于提取指定字段（如：code, message, data.status）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出（仅在 format 为空时有效）")

	return cmd
}

// createHotDataCmd 创建获取热搜数据子命令
func (s *CrawlerSkill) createHotDataCmd() *cobra.Command {
	var (
		date   string
		source string
		minify bool
	)

	cmd := &cobra.Command{
		Use:   CmdHotData,
		Short: "获取热搜数据",
		Long:  "查询指定日期的热搜数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireMCPURL(cmd); !ok {
				return
			}

			resp, err := GetHotData(date, source)
			if err != nil {
				cmd.Printf("获取热搜数据失败: %v\n", err)
				return
			}

			// 打印 JSON 响应
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&date, "date", "", "日期（格式：YYYY-MM-DD）")
	cmd.Flags().StringVar(&source, "source", "", "数据源（可选）")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

// createWeatherDataCmd 创建获取天气数据子命令
func (s *CrawlerSkill) createWeatherDataCmd() *cobra.Command {
	var (
		cityName string
		date     string
		source   string
		minify   bool
	)

	cmd := &cobra.Command{
		Use:   CmdWeatherData,
		Short: "获取天气数据",
		Long:  "查询指定城市的天气数据",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := s.requireMCPURL(cmd); !ok {
				return
			}

			resp, err := GetWeatherData(cityName, date, source)
			if err != nil {
				cmd.Printf("获取天气数据失败: %v\n", err)
				return
			}

			// 打印 JSON 响应
			skills.PrintJSONResponse(cmd, resp, minify)
		},
	}

	// 定义 flags
	cmd.Flags().StringVar(&cityName, "city", "", "城市名称")
	cmd.Flags().StringVar(&date, "date", "", "日期（格式：YYYY-MM-DD，可选）")
	cmd.Flags().StringVar(&source, "source", "", "数据源，可选值：cma,itboy,sm_search,so_search,toutiao_search,baidu_search,bing_search")
	cmd.Flags().BoolVar(&minify, "minify", false, "紧凑输出")

	// 标记必填参数
	_ = cmd.MarkFlagRequired("city")

	return cmd
}

// GetSkillInfo 获取完整的 SkillInfo（包含自动提取的 Commands）
func GetSkillInfo() types.SkillInfo {
	skill := NewCrawlerSkill()
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

// NewCrawlerCmd 创建爬虫数据命令
func NewCrawlerCmd() *cobra.Command {
	skill := NewCrawlerSkill()
	return skill.GetCommand()
}
