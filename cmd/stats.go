package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dtapp/skills/stats"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

// tableModel 表格模型
type tableModel struct {
	table table.Model
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() tea.View {
	return tea.NewView(m.table.View() + "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("按 q 退出"))
}

// NewStatsCmd 创建统计命令
func NewStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "查看技能使用统计",
		Long:  "查看和管理技能使用统计数据，包括调用次数、成功率、执行时间等",
	}

	cmd.AddCommand(newStatsListCmd())
	cmd.AddCommand(newStatsShowCmd())
	cmd.AddCommand(newStatsRecentCmd())
	cmd.AddCommand(newStatsResetCmd())
	cmd.AddCommand(newStatsPathCmd())

	return cmd
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

			// 简单文本输出
			cmd.Println()
			cmd.Println("┌─────────────────────┬──────────┬──────────┬──────────┬─────────────┬─────────────────────┐")
			cmd.Println("│ 技能名称            │ 总调用   │ 成功     │ 失败     │ 平均耗时    │ 最后使用            │")
			cmd.Println("├─────────────────────┼──────────┼──────────┼──────────┼─────────────┼─────────────────────┤")

			for _, s := range skills {
				name := padRight(s.name, 19)
				total := padRight(fmt.Sprintf("%d", s.stat.TotalCalls), 8)
				success := padRight(fmt.Sprintf("%d", s.stat.SuccessCalls), 8)
				failed := padRight(fmt.Sprintf("%d", s.stat.FailedCalls), 8)
				avgTime := padRight(fmt.Sprintf("%.0fms", s.stat.AvgDuration), 11)
				lastUsed := s.stat.LastUsedAt.Format("2006-01-02 15:04")

				cmd.Printf("│ %s │ %s │ %s │ %s │ %s │ %s │\n",
					name, total, success, failed, avgTime, lastUsed)
			}

			cmd.Println("└─────────────────────┴──────────┴──────────┴──────────┴─────────────┴─────────────────────┘")
			cmd.Printf("\n总计调用次数: %d\n", manager.GetTotalCalls())

			return nil
		},
	}

	cmd.Flags().BoolVarP(&useTable, "table", "t", false, "使用交互式表格显示")

	return cmd
}

type skillInfo struct {
	name string
	stat *stats.SkillStat
}

// renderSkillsTable 使用 bubbles table 渲染
func renderSkillsTable(skills []skillInfo, totalCalls int64) error {
	columns := []table.Column{
		{Title: "技能名称", Width: 20},
		{Title: "总调用", Width: 10},
		{Title: "成功", Width: 10},
		{Title: "失败", Width: 10},
		{Title: "平均耗时", Width: 12},
		{Title: "最后使用", Width: 20},
	}

	var rows []table.Row
	for _, s := range skills {
		rows = append(rows, table.Row{
			s.name,
			fmt.Sprintf("%d", s.stat.TotalCalls),
			fmt.Sprintf("%d", s.stat.SuccessCalls),
			fmt.Sprintf("%d", s.stat.FailedCalls),
			fmt.Sprintf("%.0fms", s.stat.AvgDuration),
			s.stat.LastUsedAt.Format("2006-01-02 15:04"),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{table: t}

	// 检测是否有 TTY
	if _, err := tea.NewProgram(m, tea.WithoutRenderer()).Run(); err != nil {
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			// 无 TTY 环境，使用简单输出
			fmt.Println("技能名称            总调用     成功       失败       平均耗时     最后使用")
			fmt.Println(strings.Repeat("-", 80))
			for _, s := range skills {
				fmt.Printf("%-20s %-10d %-10d %-10d %-12.0f %s\n",
					s.name, s.stat.TotalCalls, s.stat.SuccessCalls, s.stat.FailedCalls,
					s.stat.AvgDuration, s.stat.LastUsedAt.Format("2006-01-02 15:04"))
			}
			fmt.Printf("\n总计调用次数: %d\n", totalCalls)
			return nil
		}
		return err
	}

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

				// 使用 table 渲染命令统计
				if err := renderCommandsTable(stat.Commands); err != nil {
					// 回退到简单输出
					cmd.Println("  ┌─────────────────────┬──────────┬─────────────────────┐")
					cmd.Println("  │ 命令                │ 调用次数 │ 最后使用            │")
					cmd.Println("  ├─────────────────────┼──────────┼─────────────────────┤")

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

					for _, c := range cmds {
						name := padRight(c.name, 19)
						calls := padRight(fmt.Sprintf("%d", c.stat.Calls), 8)
						lastUsed := c.stat.LastUsed.Format("2006-01-02 15:04")
						cmd.Printf("  │ %s │ %s │ %s │\n", name, calls, lastUsed)
					}

					cmd.Println("  └─────────────────────┴──────────┴─────────────────────┘")
				}
			}

			cmd.Println()
			return nil
		},
	}
}

// renderCommandsTable 渲染命令表格
func renderCommandsTable(commands map[string]*stats.CommandStat) error {
	columns := []table.Column{
		{Title: "命令", Width: 20},
		{Title: "调用次数", Width: 12},
		{Title: "最后使用", Width: 20},
	}

	// 按调用次数排序
	type cmdInfo struct {
		name string
		stat *stats.CommandStat
	}
	var cmds []cmdInfo
	for name, cstat := range commands {
		cmds = append(cmds, cmdInfo{name: name, stat: cstat})
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].stat.Calls > cmds[j].stat.Calls
	})

	var rows []table.Row
	for _, c := range cmds {
		rows = append(rows, table.Row{
			c.name,
			fmt.Sprintf("%d", c.stat.Calls),
			c.stat.LastUsed.Format("2006-01-02 15:04"),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{table: t}

	// 检测是否有 TTY
	if _, err := tea.NewProgram(m, tea.WithoutRenderer()).Run(); err != nil {
		if !strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			return err
		}
		// 无 TTY 返回错误让调用方处理
		return fmt.Errorf("no tty")
	}

	return nil
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
			cmd.Println("┌─────────────────────┬─────────────────────┬──────────┬────────┬──────────────────────────┐")
			cmd.Println("│ 时间                │ 技能                │ 命令     │ 耗时   │ 结果                     │")
			cmd.Println("├─────────────────────┼─────────────────────┼──────────┼────────┼──────────────────────────┤")

			for _, r := range records {
				timeStr := r.StartTime.Format("01-02 15:04:05")
				skill := padRight(r.Skill, 19)
				command := padRight(r.Command, 8)
				duration := padRight(fmt.Sprintf("%dms", r.Duration.Milliseconds()), 6)

				result := "✓ 成功"
				if !r.Success {
					result = "✗ 失败"
					if r.ErrorMsg != "" {
						result = fmt.Sprintf("✗ %s", truncate(r.ErrorMsg, 20))
					}
				}
				result = padRight(result, 24)

				cmd.Printf("│ %s │ %s │ %s │ %s │ %s │\n",
					timeStr, skill, command, duration, result)
			}

			cmd.Println("└─────────────────────┴─────────────────────┴──────────┴────────┴──────────────────────────┘")
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
	columns := []table.Column{
		{Title: "时间", Width: 20},
		{Title: "技能", Width: 20},
		{Title: "命令", Width: 12},
		{Title: "耗时", Width: 10},
		{Title: "结果", Width: 25},
	}

	var rows []table.Row
	for _, r := range records {
		result := "✓ 成功"
		if !r.Success {
			result = "✗ 失败"
			if r.ErrorMsg != "" {
				result = fmt.Sprintf("✗ %s", truncate(r.ErrorMsg, 20))
			}
		}

		rows = append(rows, table.Row{
			r.StartTime.Format("01-02 15:04:05"),
			r.Skill,
			r.Command,
			fmt.Sprintf("%dms", r.Duration.Milliseconds()),
			result,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{table: t}

	// 检测是否有 TTY
	if _, err := tea.NewProgram(m, tea.WithoutRenderer()).Run(); err != nil {
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			// 无 TTY 环境，使用简单输出
			fmt.Println("时间                 技能                 命令         耗时       结果")
			fmt.Println(strings.Repeat("-", 90))
			for _, r := range records {
				result := "✓ 成功"
				if !r.Success {
					result = "✗ 失败"
				}
				fmt.Printf("%-20s %-20s %-12s %-10s %s\n",
					r.StartTime.Format("01-02 15:04:05"),
					r.Skill, r.Command,
					fmt.Sprintf("%dms", r.Duration.Milliseconds()),
					result)
			}
			return nil
		}
		return err
	}

	return nil
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
		Short: "显示统计数据文件路径",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := stats.GetManager()
			if err != nil {
				return fmt.Errorf("初始化统计管理器失败: %w", err)
			}

			cmd.Printf("统计数据文件: %s\n", manager.GetStatsFilePath())
			return nil
		},
	}
}

// padRight 右填充字符串
func padRight(s string, length int) string {
	if len(s) > length {
		return s[:length-3] + "..."
	}
	return s + strings.Repeat(" ", length-len(s))
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// isTTY 检查是否有 TTY
func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
