//go:build windows

package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// binaryExt 返回二进制文件扩展名
func binaryExt() string {
	return ".exe"
}

// updateWithPlatform Windows 平台特定的更新逻辑
func (u *Updater) updateWithPlatform(binaryURL, targetPath string) error {
	// Windows 使用标准二进制升级
	return u.updateBinary(binaryURL, targetPath)
}

// needElevation 检查是否需要提升权限（管理员权限）
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
	tmpFile, err := os.CreateTemp("", "dtskills-update-*.exe")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}

	fmt.Printf("下载位置: %s\n", tmpFile.Name())

	if _, err := io.Copy(tmpFile, reader); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 创建升级脚本
	scriptPath := tmpFile.Name() + ".bat"
	scriptContent := fmt.Sprintf(`@echo off
chcp 65001 >nul
echo 等待原程序退出...
timeout /t 2 /nobreak >nul
echo 正在替换文件...
copy /Y "%s" "%s"
if errorlevel 1 (
    echo 替换失败，请手动执行:
    echo copy /Y "%s" "%s"
    pause
    exit /b 1
)
echo 升级成功!
del "%%~f0"
"%s"
`, tmpFile.Name(), targetPath, tmpFile.Name(), targetPath, targetPath)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		_ = os.Remove(tmpFile.Name())
		return fmt.Errorf("创建升级脚本失败: %w", err)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("正在启动升级程序...")
	fmt.Println("========================================")

	// 启动升级脚本并退出当前程序
	cmd := exec.Command("cmd", "/c", "start", "", scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	if err := cmd.Start(); err != nil {
		_ = os.Remove(tmpFile.Name())
		_ = os.Remove(scriptPath)
		return fmt.Errorf("启动升级脚本失败: %w", err)
	}

	// 退出当前程序，让脚本完成替换
	fmt.Println("程序将在 2 秒后退出以完成升级...")
	time.Sleep(2 * time.Second)
	os.Exit(0)

	return nil
}
