//go:build darwin

package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// needElevation 检查是否需要提升权限
func needElevation(targetPath string) bool {
	// 检查文件是否可写
	if _, err := os.OpenFile(targetPath, os.O_WRONLY, 0); err != nil {
		return true
	}
	// 检查所在目录是否可写
	dir := filepath.Dir(targetPath)
	testFile := filepath.Join(dir, fmt.Sprintf(".write_test_%d", time.Now().UnixNano()))
	if f, err := os.Create(testFile); err != nil {
		return true
	} else {
		_ = f.Close()
		_ = os.Remove(testFile)
	}
	return false
}

// updateWithElevation 使用提升的权限更新
func updateWithElevation(reader io.Reader, targetPath string) error {
	// 下载到临时文件
	tmpFile, err := os.CreateTemp("", "dtskills-update-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := io.Copy(tmpFile, reader); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 设置执行权限
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("设置权限失败: %w", err)
	}

	// 使用 sudo 替换系统文件
	fmt.Println("需要 sudo 权限来更新系统安装的二进制文件...")
	cmd := exec.Command("sudo", "cp", tmpFile.Name(), targetPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo 更新失败: %w", err)
	}
	return nil
}
