package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ProgressOptions 进度条选项
type ProgressOptions struct {
	TotalSize   int64
	Description string
	Width       int
	Reader      io.Reader
	Writer      io.Writer
}

// progressModel 进度条模型
type progressModel struct {
	progress  progress.Model
	totalSize int64
	completed int64
	reader    io.Reader
	writer    io.Writer
	desc      string
	err       error
	done      bool
	startTime time.Time
	buffer    []byte
}

// progressTickMsg 进度消息
type progressTickMsg struct {
	n    int
	err  error
	done bool
}

// Init 初始化
func (m progressModel) Init() tea.Cmd {
	return tea.Batch(
		m.progress.Init(),
		func() tea.Msg {
			return m.readChunk()()
		},
	)
}

// Update 更新
func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressTickMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		if msg.done {
			m.done = true
			m.completed = m.totalSize
			m.progress.SetPercent(1.0)
			return m, tea.Quit
		}
		m.completed += int64(msg.n)
		if m.totalSize > 0 {
			percent := float64(m.completed) / float64(m.totalSize)
			m.progress.SetPercent(percent)
		}
		return m, m.readChunk()

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

// View 视图
func (m progressModel) View() tea.View {
	if m.err != nil {
		return tea.NewView(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("❌ %s 失败: %v\n", m.desc, m.err)))
	}

	if m.done {
		elapsed := time.Since(m.startTime)
		speed := float64(m.totalSize) / elapsed.Seconds() / 1024 / 1024
		return tea.NewView(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).
			Render(fmt.Sprintf("✓ %s 完成! (%.2f MB/s)\n", m.desc, speed)))
	}

	elapsed := time.Since(m.startTime)
	var speed float64
	if elapsed > 0 {
		speed = float64(m.completed) / elapsed.Seconds() / 1024 / 1024
	}

	percent := 0.0
	if m.totalSize > 0 {
		percent = float64(m.completed) / float64(m.totalSize) * 100
	}

	return tea.NewView(fmt.Sprintf(
		"%s\n%s %.1f%% (%.2f MB/s)\n",
		m.progress.View(),
		m.desc,
		percent,
		speed,
	))
}

// readChunk 读取数据块
func (m *progressModel) readChunk() tea.Cmd {
	return func() tea.Msg {
		n, err := m.reader.Read(m.buffer)
		if n > 0 {
			if _, werr := m.writer.Write(m.buffer[:n]); werr != nil {
				return progressTickMsg{err: werr}
			}
		}
		if err == io.EOF {
			return progressTickMsg{n: n, done: true}
		}
		if err != nil {
			return progressTickMsg{err: err}
		}
		return progressTickMsg{n: n}
	}
}

// CopyWithProgress 带进度条的复制操作
func CopyWithProgress(opts ProgressOptions) error {
	width := opts.Width
	if width <= 0 {
		width = 50
	}

	p := progress.New(
		progress.WithWidth(width),
	)

	m := progressModel{
		progress:  p,
		totalSize: opts.TotalSize,
		reader:    opts.Reader,
		writer:    opts.Writer,
		desc:      opts.Description,
		startTime: time.Now(),
		buffer:    make([]byte, 64*1024), // 64KB 缓冲区
	}

	// 使用 WithoutRenderer 选项以支持无 TTY 环境
	program := tea.NewProgram(m, tea.WithoutRenderer())
	if _, err := program.Run(); err != nil {
		// 如果是因为没有 TTY，则回退到简单复制
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			return fallbackCopy(opts)
		}
		return err
	}

	return m.err
}

// fallbackCopy 无 TTY 环境下的简单复制
func fallbackCopy(opts ProgressOptions) error {
	const chunkSize = 64 * 1024 // 64KB
	buf := make([]byte, chunkSize)
	var completed int64
	startTime := time.Now()

	for {
		n, err := opts.Reader.Read(buf)
		if n > 0 {
			if _, werr := opts.Writer.Write(buf[:n]); werr != nil {
				return fmt.Errorf("写入失败: %w", werr)
			}
			completed += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取失败: %w", err)
		}
	}

	elapsed := time.Since(startTime)
	speed := float64(completed) / elapsed.Seconds() / 1024 / 1024
	fmt.Printf("%s 完成! (%.2f MB/s)\n", opts.Description, speed)
	return nil
}

// ProgressReader 带进度报告的 Reader
type ProgressReader struct {
	Reader     io.Reader
	TotalSize  int64
	OnProgress func(completed, total int64)
	completed  int64
}

// Read 实现 io.Reader
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if n > 0 {
		pr.completed += int64(n)
		if pr.OnProgress != nil {
			pr.OnProgress(pr.completed, pr.TotalSize)
		}
	}
	return n, err
}

// NewProgressReader 创建带进度报告的 Reader
func NewProgressReader(reader io.Reader, totalSize int64, onProgress func(completed, total int64)) *ProgressReader {
	return &ProgressReader{
		Reader:     reader,
		TotalSize:  totalSize,
		OnProgress: onProgress,
	}
}

// SimpleProgress 简单文本进度条
func SimpleProgress(current, total int64, width int) string {
	if total <= 0 {
		return strings.Repeat("█", width)
	}
	percent := float64(current) / float64(total)
	filled := min(int(percent*float64(width)), width)
	empty := width - filled

	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// FormatBytes 格式化字节大小
func FormatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatDuration 格式化持续时间
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", float64(d)/float64(time.Second))
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", float64(d)/float64(time.Minute))
	}
	return fmt.Sprintf("%.1fh", float64(d)/float64(time.Hour))
}
