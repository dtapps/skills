package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// TableColumn 表格列定义
type TableColumn struct {
	Title string
	Width int
}

// TableRow 表格行数据
type TableRow []string

// TableOptions 表格选项
type TableOptions struct {
	Columns      []TableColumn
	Rows         []TableRow
	Height       int
	Focused      bool
	ShowHelp     bool
	HelpText     string
	FallbackText bool // 无 TTY 时是否使用文本输出
}

// TableModel 表格模型
type TableModel struct {
	table    table.Model
	options  TableOptions
	quitting bool
}

// Init 初始化
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View 视图
func (m TableModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	view := m.table.View()
	if m.options.ShowHelp {
		helpText := m.options.HelpText
		if helpText == "" {
			helpText = "按 q 退出"
		}
		view += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(helpText)
	}
	return tea.NewView(view)
}

// RenderTable 渲染表格
// 有 TTY 时使用 bubbles table，无 TTY 时回退到文本输出
func RenderTable(opts TableOptions) error {
	// 转换为 bubbles table 的列定义
	columns := make([]table.Column, len(opts.Columns))
	for i, col := range opts.Columns {
		columns[i] = table.Column{Title: col.Title, Width: col.Width}
	}

	// 转换为 bubbles table 的行数据
	rows := make([]table.Row, len(opts.Rows))
	for i, row := range opts.Rows {
		rows[i] = table.Row(row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(opts.Focused),
		table.WithHeight(opts.Height),
	)

	// 设置样式
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

	m := TableModel{table: t, options: opts}

	// 尝试运行程序
	if _, err := tea.NewProgram(m, tea.WithoutRenderer()).Run(); err != nil {
		// 无 TTY 环境
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			if opts.FallbackText {
				return renderTextTable(opts)
			}
			return fmt.Errorf("需要交互式终端才能显示表格")
		}
		return err
	}

	return nil
}

// renderTextTable 文本表格输出（无 TTY 回退）
func renderTextTable(opts TableOptions) error {
	// 打印表头
	var headers []string
	for _, col := range opts.Columns {
		headers = append(headers, padRight(col.Title, col.Width))
	}
	fmt.Println(strings.Join(headers, " "))

	// 打印分隔线
	var sep []string
	for _, col := range opts.Columns {
		sep = append(sep, strings.Repeat("-", col.Width))
	}
	fmt.Println(strings.Join(sep, " "))

	// 打印数据行
	for _, row := range opts.Rows {
		var cells []string
		for i, cell := range row {
			if i < len(opts.Columns) {
				cells = append(cells, padRight(cell, opts.Columns[i].Width))
			}
		}
		fmt.Println(strings.Join(cells, " "))
	}

	return nil
}

// padRight 右填充字符串
func padRight(s string, length int) string {
	if len(s) > length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// SimpleTable 简单表格（纯文本，无需 TTY）
func SimpleTable(headers []string, rows [][]string) {
	// 计算列宽
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 打印表头
	printRow(headers, colWidths)

	// 打印分隔线
	var sep []string
	for i, width := range colWidths {
		if i < len(headers) {
			sep = append(sep, strings.Repeat("─", width+2))
		}
	}
	fmt.Println("├" + strings.Join(sep, "┼") + "┤")

	// 打印数据行
	for _, row := range rows {
		printRow(row, colWidths)
	}
}

// printRow 打印一行
func printRow(cells []string, widths []int) {
	var parts []string
	for i, cell := range cells {
		if i < len(widths) {
			parts = append(parts, " "+padRight(cell, widths[i])+" ")
		}
	}
	fmt.Println("│" + strings.Join(parts, "│") + "│")
}

// BoxTable 带边框的表格
func BoxTable(headers []string, rows [][]string) {
	// 计算列宽
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 打印顶部边框
	printBorder(colWidths, "┌", "┬", "┐")

	// 打印表头
	printRow(headers, colWidths)

	// 打印分隔线
	printBorder(colWidths, "├", "┼", "┤")

	// 打印数据行
	for i, row := range rows {
		printRow(row, colWidths)
		if i < len(rows)-1 {
			printBorder(colWidths, "├", "┼", "┤")
		}
	}

	// 打印底部边框
	printBorder(colWidths, "└", "┴", "┘")
}

// printBorder 打印边框
func printBorder(widths []int, left, mid, right string) {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat("─", w+2))
	}
	fmt.Println(left + strings.Join(parts, mid) + right)
}
