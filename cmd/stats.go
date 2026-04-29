package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/stats"
	"github.com/dtapp/skills/ui"

	"github.com/spf13/cobra"
)

// NewStatsCmd 创建统计命令
func NewStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stats",
		Aliases: []string{"s", "stat"},
		Short:   "查看技能使用统计",
		Long:    "查看和管理技能使用统计数据，包括调用次数、成功率、执行时间等",
	}

	cmd.AddCommand(newStatsListCmd())
	cmd.AddCommand(newStatsShowCmd())
	cmd.AddCommand(newStatsRecentCmd())
	cmd.AddCommand(newStatsResetCmd())
	cmd.AddCommand(newStatsPathCmd())

	return cmd
}

type skillInfo struct {
	name string
	stat *stats.SkillStat
}

// newStatsListCmd 列出所有技能统计
func newStatsListCmd() *cobra.Command {
	var useTable bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有技能的使用统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := stats.GetManager()
			if err != nil {
				return fmt.Errorf("初始化统计管理器失败: %w", err)
			}

			allStats := manager.GetAllStats()
			if len(allStats) == 0 {
				cmd.Println("暂无使用统计数据")
				return nil
			}

			// 按调用次数排序
			var skills []skillInfo
			for name, stat := range allStats {
				skills = append(skills, skillInfo{name: name, stat: stat})
			}
			sort.Slice(skills, func(i, j int) bool {
				return skills[i].stat.TotalCalls > skills[j].stat.TotalCalls
			})

			// 使用 bubbles table 渲染
			if useTable {
				return renderSkillsTable(skills, manager.GetTotalCalls())
			}

			// 使用 BoxTable 输出
			cmd.Println()
			var rows [][]string
			for _, s := range skills {
				rows = append(rows, []string{
					s.name,
					fmt.Sprintf("%d", s.stat.TotalCalls),
					fmt.Sprintf("%d", s.stat.SuccessCalls),
					fmt.Sprintf("%d", s.stat.FailedCalls),
					fmt.Sprintf("%.0fms", s.stat.AvgDuration),
					s.stat.LastUsedAt.Format("2006-01-02 15:04"),
				})
			}
			ui.BoxTable(
				[]string{"技能名称", "总调用", "成功", "失败", "平均耗时", "最后使用"},
				rows,
			)
			cmd.Printf("\n总计调用次数: %d\n", manager.GetTotalCalls())

			return nil
		},
	}

	cmd.Flags().BoolVarP(&useTable, "table", "t", false, "使用交互式表格显示")

	return cmd
}

// renderSkillsTable 使用 bubbles table 渲染
func renderSkillsTable(skills []skillInfo, totalCalls int64) error {
	columns := []ui.TableColumn{
		{Title: "技能名称", Width: 20},
		{Title: "总调用", Width: 10},
		{Title: "成功", Width: 10},
		{Title: "失败", Width: 10},
		{Title: "平均耗时", Width: 12},
		{Title: "最后使用", Width: 20},
	}

	var rows []ui.TableRow
	for _, s := range skills {
		rows = append(rows, ui.TableRow{
			s.name,
			fmt.Sprintf("%d", s.stat.TotalCalls),
			fmt.Sprintf("%d", s.stat.SuccessCalls),
			fmt.Sprintf("%d", s.stat.FailedCalls),
			fmt.Sprintf("%.0fms", s.stat.AvgDuration),
			s.stat.LastUsedAt.Format("2006-01-02 15:04"),
		})
	}

	opts := ui.TableOptions{
		Columns:      columns,
		Rows:         rows,
		Height:       len(rows) + 1,
		Focused:      true,
		ShowHelp:     true,
		FallbackText: true,
	}

	if err := ui.RenderTable(opts); err != nil {
		return err
	}

	fmt.Printf("\n总计调用次数: %d\n", totalCalls)
	return nil
}

// newStatsShowCmd 显示单个技能详细统计
func newStatsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <skill-name>",
		Short: "显示指定技能的详细统计",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			skillName := args[0]

			manager, err := stats.GetManager()
			if err != nil {
				return fmt.Errorf("初始化统计管理器失败: %w", err)
			}

			stat, exists := manager.GetSkillStats(skillName)
			if !exists {
				return fmt.Errorf("技能 '%s' 暂无使用统计数据", skillName)
			}

			// 计算成功率
			successRate := 0.0
			if stat.TotalCalls > 0 {
				successRate = float64(stat.SuccessCalls) / float64(stat.TotalCalls) * 100
			}

			cmd.Println()
			cmd.Printf("📊 技能: %s\n", skillName)
			cmd.Println(strings.Repeat("─", 50))
			cmd.Printf("  总调用次数:   %d\n", stat.TotalCalls)
			cmd.Printf("  成功次数:     %d (%.1f%%)\n", stat.SuccessCalls, successRate)
			cmd.Printf("  失败次数:     %d\n", stat.FailedCalls)
			cmd.Printf("  平均执行时间: %.2f ms\n", stat.AvgDuration)
			cmd.Printf("  首次使用:     %s\n", stat.FirstUsedAt.Format("2006-01-02 15:04:05"))
			cmd.Printf("  最后使用:     %s\n", stat.LastUsedAt.Format("2006-01-02 15:04:05"))

			if len(stat.Commands) > 0 {
				cmd.Println()
				cmd.Println("  命令使用情况:")

				// 按调用次数排序
				type cmdInfo struct {
					name string
					stat *stats.CommandStat
				}
				var cmds []cmdInfo
				for name, cstat := range stat.Commands {
					cmds = append(cmds, cmdInfo{name: name, stat: cstat})
				}
				sort.Slice(cmds, func(i, j int) bool {
					return cmds[i].stat.Calls > cmds[j].stat.Calls
				})

				// 使用 BoxTable 输出
				var rows [][]string
				for _, c := range cmds {
					rows = append(rows, []string{
						c.name,
						fmt.Sprintf("%d", c.stat.Calls),
						c.stat.LastUsed.Format("2006-01-02 15:04"),
					})
				}
				ui.BoxTable(
					[]string{"命令", "调用次数", "最后使用"},
					rows,
				)
			}

			cmd.Println()
			return nil
		},
	}
}

// newStatsRecentCmd 显示最近使用记录
func newStatsRecentCmd() *cobra.Command {
	var limit int
	var useTable bool

	cmd := &cobra.Command{
		Use:   "recent",
		Short: "显示最近的使用记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := stats.GetManager()
			if err != nil {
				return fmt.Errorf("初始化统计管理器失败: %w", err)
			}

			records := manager.GetRecentUsage(limit)
			if len(records) == 0 {
				cmd.Println("暂无使用记录")
				return nil
			}

			if useTable {
				return renderRecentTable(records)
			}

			cmd.Println()
			cmd.Printf("最近 %d 条使用记录:\n", len(records))

			var rows [][]string
			for _, r := range records {
				result := "✓ 成功"
				if !r.Success {
					result = "✗ 失败"
					if r.ErrorMsg != "" {
						result = fmt.Sprintf("✗ %s", truncate(r.ErrorMsg, 20))
					}
				}

				rows = append(rows, []string{
					r.StartTime.Format("01-02 15:04:05"),
					r.Skill,
					r.Command,
					fmt.Sprintf("%dms", r.Duration.Milliseconds()),
					result,
				})
			}
			ui.BoxTable(
				[]string{"时间", "技能", "命令", "耗时", "结果"},
				rows,
			)
			cmd.Println()
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "显示的记录数量")
	cmd.Flags().BoolVarP(&useTable, "table", "t", false, "使用交互式表格显示")

	return cmd
}

// renderRecentTable 渲染最近使用记录表格
func renderRecentTable(records []stats.UsageRecord) error {
	columns := []ui.TableColumn{
		{Title: "时间", Width: 20},
		{Title: "技能", Width: 20},
		{Title: "命令", Width: 12},
		{Title: "耗时", Width: 10},
		{Title: "结果", Width: 25},
	}

	var rows []ui.TableRow
	for _, r := range records {
		result := "✓ 成功"
		if !r.Success {
			result = "✗ 失败"
			if r.ErrorMsg != "" {
				result = fmt.Sprintf("✗ %s", truncate(r.ErrorMsg, 20))
			}
		}

		rows = append(rows, ui.TableRow{
			r.StartTime.Format("01-02 15:04:05"),
			r.Skill,
			r.Command,
			fmt.Sprintf("%dms", r.Duration.Milliseconds()),
			result,
		})
	}

	opts := ui.TableOptions{
		Columns:      columns,
		Rows:         rows,
		Height:       len(rows) + 1,
		Focused:      true,
		ShowHelp:     true,
		FallbackText: true,
	}

	return ui.RenderTable(opts)
}

// newStatsResetCmd 重置统计数据
func newStatsResetCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "重置所有统计数据",
		Long:  "重置所有技能使用统计数据，此操作不可恢复！",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				cmd.Println("⚠️  警告: 此操作将删除所有统计数据且不可恢复！")
				cmd.Println("请使用 --force 或 -f 参数确认重置")
				return nil
			}

			manager, err := stats.GetManager()
			if err != nil {
				return fmt.Errorf("初始化统计管理器失败: %w", err)
			}

			if err := manager.ResetStats(); err != nil {
				return fmt.Errorf("重置统计数据失败: %w", err)
			}

			cmd.Println("✓ 统计数据已重置")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制重置，无需确认")

	return cmd
}

// newStatsPathCmd 显示统计文件路径
func newStatsPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "显示统计存储路径",
		Long:  "显示统计数据文件路径和工作目录等存储相关信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := ui.NewPrinter(cmd.OutOrStdout())
			p.Title("📊 统计存储路径")
			p.Println()

			// 显示工作目录
			workDir := config.GetWorkDir()
			p.Printf("  工作目录: %s\n", workDir)

			// 显示统计数据文件路径
			manager, err := stats.GetManager()
			if err != nil {
				p.Printf("  统计文件: (初始化失败: %v)\n", err)
			} else {
				p.Printf("  统计文件: %s\n", manager.GetStatsFilePath())
			}

			// 显示环境变量（如果设置了）
			envVar := config.GetWorkDirEnv()
			if envValue := os.Getenv(envVar); envValue != "" {
				p.Printf("  环境变量: %s=%s\n", envVar, envValue)
			}

			p.Println()
			return nil
		},
	}
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
