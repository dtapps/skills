package updater

import (
	"fmt"
	"io"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// downloadProgressModel 下载进度模型
type downloadProgressModel struct {
	progress   progress.Model
	totalSize  int64
	downloaded int64
	reader     io.Reader
	writer     io.Writer
	desc       string
	err        error
	done       bool
	startTime  time.Time
	buffer     []byte
}

// downloadTickMsg 下载进度消息
type downloadTickMsg struct {
	n    int
	err  error
	done bool
}

// Init 初始化
func (m downloadProgressModel) Init() tea.Cmd {
	return tea.Batch(
		m.progress.Init(),
		func() tea.Msg {
			return m.readChunk()()
		},
	)
}

// Update 更新
func (m downloadProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case downloadTickMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		if msg.done {
			m.done = true
			m.downloaded = m.totalSize
			m.progress.SetPercent(1.0)
			return m, tea.Quit
		}
		m.downloaded += int64(msg.n)
		if m.totalSize > 0 {
			percent := float64(m.downloaded) / float64(m.totalSize)
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

// View 视图 - bubbletea v2 返回 tea.View
func (m downloadProgressModel) View() tea.View {
	if m.err != nil {
		return tea.NewView(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("❌ 下载失败: %v\n", m.err)))
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
		speed = float64(m.downloaded) / elapsed.Seconds() / 1024 / 1024
	}

	percent := 0.0
	if m.totalSize > 0 {
		percent = float64(m.downloaded) / float64(m.totalSize) * 100
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
func (m *downloadProgressModel) readChunk() tea.Cmd {
	return func() tea.Msg {
		n, err := m.reader.Read(m.buffer)
		if n > 0 {
			if _, werr := m.writer.Write(m.buffer[:n]); werr != nil {
				return downloadTickMsg{err: werr}
			}
		}
		if err == io.EOF {
			return downloadTickMsg{n: n, done: true}
		}
		if err != nil {
			return downloadTickMsg{err: err}
		}
		return downloadTickMsg{n: n}
	}
}

// downloadWithProgressBar 使用 progress bar 下载
func downloadWithProgressBar(reader io.Reader, writer io.Writer, totalSize int64, desc string) error {
	p := progress.New(
		progress.WithWidth(50),
	)

	m := downloadProgressModel{
		progress:  p,
		totalSize: totalSize,
		reader:    reader,
		writer:    writer,
		desc:      desc,
		startTime: time.Now(),
		buffer:    make([]byte, 64*1024), // 64KB 缓冲区
	}

	// 使用 WithoutRenderer 选项以支持无 TTY 环境（如测试）
	program := tea.NewProgram(m, tea.WithoutRenderer())
	if _, err := program.Run(); err != nil {
		// 如果是因为没有 TTY，则回退到简单下载
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			return fallbackDownload(reader, writer, totalSize, desc)
		}
		return err
	}

	return m.err
}

// fallbackDownload 无 TTY 环境下的简单下载
func fallbackDownload(reader io.Reader, writer io.Writer, totalSize int64, desc string) error {
	const chunkSize = 64 * 1024 // 64KB
	buf := make([]byte, chunkSize)
	var downloaded int64
	startTime := time.Now()

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, werr := writer.Write(buf[:n]); werr != nil {
				return fmt.Errorf("写入失败: %w", werr)
			}
			downloaded += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取失败: %w", err)
		}
	}

	elapsed := time.Since(startTime)
	speed := float64(downloaded) / elapsed.Seconds() / 1024 / 1024
	fmt.Printf("%s 完成! (%.2f MB/s)\n", desc, speed)
	return nil
}
