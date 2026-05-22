package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dtapp/skills/ansi"
	"github.com/dtapp/skills/cmd"
	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/generator"
	"github.com/dtapp/skills/skills/amap_map"
	"github.com/dtapp/skills/skills/baidu_baike"
	"github.com/dtapp/skills/skills/baidu_map"
	"github.com/dtapp/skills/skills/baidu_search"
	"github.com/dtapp/skills/skills/cma_weather"
	"github.com/dtapp/skills/skills/hefeng_weather"
	"github.com/dtapp/skills/skills/sojson"
	"github.com/dtapp/skills/skills/tencent_map"
	"github.com/dtapp/skills/skills/tencent_news"
	"github.com/dtapp/skills/skills/tencent_weather"
	"github.com/dtapp/skills/skills/uapi"
	"github.com/dtapp/skills/skills/weibo_search"
	"github.com/dtapp/skills/stats"
	"github.com/dtapp/skills/types"
	"github.com/dtapp/skills/ui"
	"github.com/dtapp/skills/updater"

	"github.com/spf13/cobra"
)

func init() {
	// 启用 ANSI 颜色支持（Windows 需要特殊处理）
	ansi.EnableANSI()
}

// 打包时注入
var (
	Version   = "dev"     // 版本号，线上为vX.X.X
	GitCommit = "unknown" // Git 提交 ID
	BuildTime = "unknown" // 构建时间
)

// runGenerate 执行生成 skill 的逻辑
func runGenerate(cmd *cobra.Command, args []string) {
	// 预定义的 skill 模板列表（从已注册的 skill 中获取）
	// 使用 GetSkillInfo() 获取包含自动提取 Commands 的完整信息
	availableSkills := []types.SkillInfo{
		baidu_baike.GetSkillInfo(),
		baidu_search.GetSkillInfo(),
		tencent_news.GetSkillInfo(),
		tencent_weather.GetSkillInfo(),
		// jiaozhen_factcheck.GetSkillInfo(),
		weibo_search.GetSkillInfo(),
		// weibo_user.GetSkillInfo(),
		hefeng_weather.GetSkillInfo(),
		baidu_map.GetSkillInfo(),
		amap_map.GetSkillInfo(),
		tencent_map.GetSkillInfo(),
		cma_weather.GetSkillInfo(),
		uapi.GetSkillInfo(),
		sojson.GetSkillInfo(),
	}

	// 构建选项列表（显示中文名称、原名字和版本号）
	var options []ui.Option
	for _, skill := range availableSkills {
		display := fmt.Sprintf("%s [%s] (v%s)", skill.Title, skill.Name, skill.Version)
		options = append(options, ui.Option{
			Display: display,
			Value:   skill.Name,
		})
	}

	// 多选交互
	selectedSkillNames, err := ui.RunMultiSelect("选择要生成的 skill", options)
	if err != nil {
		cmd.Printf("错误：%v\n", err)
		return
	}

	if len(selectedSkillNames) == 0 {
		cmd.Println("未选择任何 skill")
		return
	}

	// 获取输出目录（从配置文件）
	outputDir := config.GetSkillsPath()
	if outputDir == "" {
		cmd.Println("错误: 未配置 skills.path，请使用 'dtskills config set skills.path <路径>' 设置")
		return
	}

	// 创建打印机
	p := ui.NewPrinter(cmd.OutOrStdout())

	// 先显示已选择的技能列表
	p.Println()
	p.Title("📋 已选择的技能")
	p.Println()
	for _, skill := range selectedSkillNames {
		p.Printf("  ✓ %s\n", skill)
	}
	p.Println()

	// 询问是否是 picoclaw
	isPicoclaw := ui.RunConfirm("是否为 PicoClaw 生成 .skill-origin.json 文件？", false)

	// 显示摘要并确认
	p.Println()
	p.Title("📋 生成摘要")
	p.Println()
	p.Println("将要生成以下 skill：")
	for _, skill := range selectedSkillNames {
		p.Printf("  ✓ %s\n", skill)
	}
	p.Println()
	p.Printf("输出目录: %s\n", outputDir)
	p.Printf("PicoClaw 模式: %s\n", map[bool]string{true: "是", false: "否"}[isPicoclaw])

	if !ui.RunConfirm("确认生成？", true) {
		cmd.Println("已取消生成")
		return
	}

	// 生成结果记录
	results := make(map[string]error)

	// 为每个选择的 skill 生成文件
	for _, skillName := range selectedSkillNames {
		// 查找对应的 skillInfo
		var skillInfo types.SkillInfo
		for _, skill := range availableSkills {
			if skill.Name == skillName {
				skillInfo = skill
				break
			}
		}

		// 生成 skill，传入 CLI 名称
		err := generator.GenerateSkill(skillInfo, outputDir, isPicoclaw, types.ProjectName)
		results[skillName] = err
	}

	// 显示结果
	p.Println()
	successCount := 0
	failCount := 0
	for skill, err := range results {
		if err != nil {
			p.Printf("  ✗ %s: %s\n", skill, err)
			failCount++
		} else {
			p.Printf("  ✓ %s: 生成成功\n", skill)
			successCount++
		}
	}
	p.Println()
	if failCount == 0 {
		p.Title(fmt.Sprintf("✅ 全部 %d 个 skill 生成完成！", successCount))
	} else {
		p.Printf("%s (成功: %d, 失败: %d)\n",
			ui.SubtleStyle.Render("⚠️ 部分 skill 生成失败"),
			successCount, failCount)
	}
}

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Skill 管理工具",
	Long:  "用于下载、执行和管理技能包",
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有子命令，默认执行 generate
		runGenerate(cmd, args)
	},
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "生成新的 skill",
	Long:    "通过模板生成新的 skill 目录和文件",
	Run:     runGenerate,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出已有的 skill",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("已有的 skill：")

		// 遍历当前目录，查找 skill 目录
		dir, err := os.ReadDir("..")
		if err != nil {
			cmd.Printf("错误：%v\n", err)
			return
		}

		skillsFound := false
		for _, entry := range dir {
			if entry.IsDir() {
				// 检查目录是否包含 SKILL.md 文件
				skillPath := filepath.Join("..", entry.Name())
				skillMdPath := filepath.Join(skillPath, "SKILL.md")
				if _, err := os.Stat(skillMdPath); !os.IsNotExist(err) {
					cmd.Printf("- %s\n", entry.Name())
					skillsFound = true
				}
			}
		}

		if !skillsFound {
			cmd.Println("未找到任何 skill")
		}
	},
}

// versionCmd 显示版本信息
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		p := ui.NewPrinter(cmd.OutOrStdout())
		p.Title("📦 " + types.ProjectName)
		p.Printf("版本: %s\n", Version)
		p.Printf("Git 提交: %s\n", GitCommit)
		p.Printf("构建时间: %s\n", BuildTime)
		p.Printf("平台: %s\n", updater.GetPlatformInfo())
	},
}

// upgradeCmd 检查并执行更新
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"u", "up"},
	Short:   "检查并更新到最新版本",
	Long:    "检查是否有新版本可用，并下载更新",
	Run: func(cmd *cobra.Command, args []string) {
		p := ui.NewPrinter(cmd.OutOrStdout())

		// 获取参数
		repoURL, _ := cmd.Flags().GetString("repo")
		force, _ := cmd.Flags().GetBool("force")

		// 创建更新器
		u := updater.New(Version, GitCommit, repoURL)

		p.Title("🔍 检查更新")
		p.Printf("当前版本: %s\n", Version)
		if GitCommit != "unknown" {
			p.Printf("当前 Commit: %s\n", GitCommit)
		}
		p.Printf("平台: %s\n", updater.GetPlatformInfo())
		p.Println()

		// 检查更新
		hasUpdate, latestVersion, latestCommit, downloadURL, err := u.Check(force)
		if err != nil {
			p.Printf("检查更新失败: %v\n", err)
			return
		}

		if !hasUpdate {
			p.Title("✅ 已是最新版本")
			return
		}

		p.Printf("发现新版本: %s\n", latestVersion)
		if latestCommit != "" {
			p.Printf("最新 Commit: %s\n", latestCommit)
		}
		p.Printf("下载地址: %s\n", downloadURL)
		p.Println()

		// 确认是否更新
		if !ui.RunConfirm("是否立即更新？", true) {
			p.Println("已取消更新")
			return
		}

		// 执行更新
		p.Println()
		p.Title("⬇️ 正在下载更新...")

		if err := u.Update(downloadURL); err != nil {
			p.Printf("更新失败: %v\n", err)
			return
		}

		p.Println()
		p.Title("✅ 更新成功！")
		p.Println("请重新运行程序以使用新版本。")
	},
}

func init() {
	// 为 upgrade 命令添加 flag
	upgradeCmd.Flags().String("repo", "", "仓库地址（可选，默认: https://cnb.cool/dtapp/skills）")
	upgradeCmd.Flags().Bool("force", false, "强制升级，即使版本号相同")
}

func main() {
	// 加载配置文件
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 加载配置文件失败: %v\n", err)
	}

	// 添加子命令
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(cmd.NewConfigCmd())
	rootCmd.AddCommand(cmd.NewStatsCmd())

	// 添加 skill 子命令（包装统计）
	rootCmd.AddCommand(stats.WrapSkillCommands(baidu_baike.NewBaiduBaikeCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(baidu_search.NewBaiduSearchCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(tencent_news.NewTencentNewsCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(tencent_weather.NewTencentWeatherCmd()))
	// rootCmd.AddCommand(stats.WrapSkillCommands(jiaozhen_factcheck.NewJiaozhenFactCheckCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(weibo_search.NewWeiboSearchCmd()))
	// rootCmd.AddCommand(stats.WrapSkillCommands(weibo_user.NewWeiboUserCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(hefeng_weather.NewHefengWeatherCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(baidu_map.NewBaiduMapSkill().GetCommand()))
	rootCmd.AddCommand(stats.WrapSkillCommands(amap_map.NewAmapMapSkill().GetCommand()))
	rootCmd.AddCommand(stats.WrapSkillCommands(tencent_map.NewTencentMapSkill().GetCommand()))
	rootCmd.AddCommand(stats.WrapSkillCommands(cma_weather.NewCmaWeatherCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(uapi.NewUapiCmd()))
	rootCmd.AddCommand(stats.WrapSkillCommands(sojson.NewSojsonCmd()))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
