package ui

import (
	"fmt"
	"io"

	"charm.land/lipgloss/v2"
)

var (
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	SubtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// Printer 输出封装
type Printer struct {
	w io.Writer
}

// NewPrinter 创建打印机
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print 打印
func (p *Printer) Print(a ...any) {
	_, _ = fmt.Fprint(p.w, a...)
}

// Println 打印并换行
func (p *Printer) Println(a ...any) {
	_, _ = fmt.Fprintln(p.w, a...)
}

// Printf 格式化打印
func (p *Printer) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(p.w, format, a...)
}

// Title 打印标题
func (p *Printer) Title(a ...any) {
	_, _ = fmt.Fprintln(p.w, TitleStyle.Render(fmt.Sprint(a...)))
}

// Subtle 打印次要信息
func (p *Printer) Subtle(a ...any) {
	_, _ = fmt.Fprintln(p.w, SubtleStyle.Render(fmt.Sprint(a...)))
}
