package updater

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dtapp/skills/ui"
)

// downloadWithProgressBar 使用 progress bar 下载
func downloadWithProgressBar(reader io.Reader, writer io.Writer, totalSize int64, desc string) error {
	opts := ui.ProgressOptions{
		TotalSize:   totalSize,
		Description: desc,
		Width:       50,
		Reader:      reader,
		Writer:      writer,
	}

	// 尝试使用进度条
	if err := ui.CopyWithProgress(opts); err != nil {
		// 如果是因为没有 TTY，则回退到简单下载
		if strings.Contains(err.Error(), "open /dev/tty: device not configured") {
			return fallbackDownload(reader, writer, totalSize, desc)
		}
		return err
	}

	return nil
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
