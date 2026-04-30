//go:build windows

package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/dtapp/skills/types"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// InstallType 安装类型
type InstallType int

const (
	InstallTypeBinary InstallType = iota // 二进制安装（便携版）
	InstallTypeMSI                       // MSI 安装包安装
)

// binaryExt 返回二进制文件扩展名
func binaryExt() string {
	return ".exe"
}

// archiveExt 返回压缩包扩展名
func archiveExt() string {
	return ".zip"
}

// msiExt 返回 MSI 安装包扩展名
func msiExt() string {
	return ".msi"
}

// detectInstallType 检测 Windows 安装类型
func detectInstallType(targetPath string) InstallType {
	// 通过查询 WMI 检测是否是通过 MSI 安装的
	// 查询 Win32_Product，检查是否有我们的产品
	productCode := getMSIProductCode()
	if productCode != "" {
		return InstallTypeMSI
	}

	// 检查是否在 Program Files 目录（通常是 MSI 安装的位置）
	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFiles != "" && strings.HasPrefix(targetPath, programFiles) {
		// 在 Program Files 中，可能是 MSI 安装的，再检查注册表
		if isInstalledByMSI() {
			return InstallTypeMSI
		}
	}
	if programFilesX86 != "" && strings.HasPrefix(targetPath, programFilesX86) {
		if isInstalledByMSI() {
			return InstallTypeMSI
		}
	}

	return InstallTypeBinary
}

// getMSIProductCode 通过 WMI 查询 MSI 产品代码
func getMSIProductCode() string {
	// 使用 wmic 查询产品
	// #nosec G204 - wmic 是系统命令
	cmd := exec.Command("wmic", "product", "where", fmt.Sprintf("Name='%s'", types.ProjectName), "get", "IdentifyingNumber", "/value")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// 解析输出，提取产品代码
	outputStr := string(output)
	if idx := strings.Index(outputStr, "IdentifyingNumber="); idx != -1 {
		code := strings.TrimSpace(outputStr[idx+18:])
		code = strings.Trim(code, "\r\n")
		return code
	}
	return ""
}

// isInstalledByMSI 检查注册表确认是否通过 MSI 安装
func isInstalledByMSI() bool {
	// 检查卸载注册表项
	// 这里简化处理，实际应该读取注册表
	// 如果程序在 Program Files 且不在 PATH 中，可能是 MSI 安装的
	return false
}

// buildPackageURL 构建包下载 URL
func buildPackageURL(repoURL, version string, installType InstallType) string {
	switch installType {
	case InstallTypeMSI:
		// 根据架构选择 MSI 文件
		arch := "amd64"
		if os.Getenv("PROCESSOR_ARCHITECTURE") == "ARM64" {
			arch = "arm64"
		}
		return fmt.Sprintf("%s/-/releases/download/%s/%s_%s.msi", repoURL, version, types.ProjectName, arch)
	default:
		return fmt.Sprintf("%s/-/releases/download/%s/%s-windows-amd64%s", repoURL, version, types.ProjectName, archiveExt())
	}
}

// updateWithPlatform Windows 平台特定的更新逻辑
func (u *Updater) updateWithPlatform(binaryURL, targetPath string) error {
	// 检测安装类型
	installType := detectInstallType(targetPath)

	// 如果是 MSI 安装，使用 MSI 升级
	if installType == InstallTypeMSI {
		return u.updateWithMSI(targetPath)
	}

	// 便携版安装，使用标准升级流程
	return u.updateBinary(binaryURL, targetPath)
}

// updateWithMSI 使用 MSI 安装包升级
func (u *Updater) updateWithMSI(targetPath string) error {
	// 获取最新版本
	latestVersion, _, err := u.getLatestVersionAndCommit()
	if err != nil {
		return fmt.Errorf("获取最新版本失败: %w", err)
	}

	msiURL := buildPackageURL(u.repoURL, latestVersion, InstallTypeMSI)

	fmt.Printf("检测到通过 MSI 安装，正在下载 MSI 安装包...\n")
	fmt.Printf("下载地址: %s\n", msiURL)

	// 下载 MSI 文件
	// #nosec G107 - msiURL 是程序内部构造的 URL，非用户输入
	resp, err := http.Get(msiURL)
	if err != nil {
		return fmt.Errorf("下载 MSI 失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载 MSI 失败，服务器返回: %d", resp.StatusCode)
	}

	// 保存到临时文件
	tmpFile, err := os.CreateTemp("", "dtskills-update-*.msi")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("写入 MSI 文件失败: %w", err)
	}
	_ = tmpFile.Close()

	fmt.Printf("MSI 下载完成: %s\n", tmpFile.Name())
	fmt.Println("正在启动 MSI 安装程序...")

	// 使用 msiexec 执行静默安装
	// #nosec G204 - msiexec 是系统命令，tmpFile.Name() 是程序内部生成的
	cmd := exec.Command("msiexec", "/i", tmpFile.Name(), "/qb", "/norestart")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("MSI 安装失败: %w", err)
	}

	return nil
}

// needElevation 检查是否需要提升权限（管理员权限）
func needElevation(targetPath string) bool {
	// 检查文件是否可写
	// #nosec G304 G703 - targetPath 经过验证的系统路径
	if _, err := os.OpenFile(targetPath, os.O_WRONLY, 0); err != nil {
		return true
	}
	// 检查所在目录是否可写
	dir := filepath.Dir(targetPath)
	testFile := filepath.Join(dir, fmt.Sprintf(".write_test_%d", time.Now().UnixNano()))
	// #nosec G304 G703 - testFile 是程序内部生成的临时测试文件路径
	if f, err := os.Create(testFile); err != nil {
		return true
	} else {
		_ = f.Close()
		// #nosec G703 - testFile 是程序内部生成的临时测试文件路径
		_ = os.Remove(testFile)
	}
	return false
}

// updateWithElevation 使用提升的权限更新
// Windows 特殊处理：将新文件下载到临时位置，然后通过批处理脚本替换
func updateWithElevation(reader io.Reader, targetPath string) error {
	// 使用系统临时目录
	tempDir := os.TempDir()

	// 创建固定的临时文件名
	timestamp := time.Now().Unix()
	newExeName := fmt.Sprintf("dtskills_new_%d.exe", timestamp)
	scriptName := fmt.Sprintf("dtskills_update_%d.bat", timestamp)

	newExePath := filepath.Join(tempDir, newExeName)
	scriptPath := filepath.Join(tempDir, scriptName)

	// 下载新文件到临时位置
	// #nosec G304 - newExePath 是程序内部生成的临时文件路径，非用户输入
	newExeFile, err := os.Create(newExePath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}

	fmt.Printf("下载位置: %s\n", newExePath)

	if _, err := io.Copy(newExeFile, reader); err != nil {
		_ = newExeFile.Close()
		_ = os.Remove(newExePath)
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := newExeFile.Close(); err != nil {
		_ = os.Remove(newExePath)
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 创建升级脚本
	// 使用 CRLF 换行符，Windows 批处理必需
	// chcp 936 设置 GBK 编码，确保中文正常显示
	scriptContent := "@echo off\r\n" +
		"chcp 936 >nul\r\n" +
		"echo 等待原程序退出...\r\n" +
		":wait_loop\r\n" +
		"timeout /t 1 /nobreak >nul\r\n" +
		"move /Y \"" + newExePath + "\" \"" + targetPath + "\" >nul 2>&1\r\n" +
		"if errorlevel 1 (\r\n" +
		"    echo 文件被占用，继续等待...\r\n" +
		"    goto :wait_loop\r\n" +
		")\r\n" +
		"echo 升级成功!\r\n" +
		"echo 正在启动新版本...\r\n" +
		"start \"\" \"" + targetPath + "\"\r\n" +
		"del \"" + scriptPath + "\"\r\n"

	// 将脚本内容转换为 GBK 编码，避免中文乱码
	encoder := simplifiedchinese.GBK.NewEncoder()
	scriptBytes, _, err := transform.Bytes(encoder, []byte(scriptContent))
	if err != nil {
		_ = os.Remove(newExePath)
		return fmt.Errorf("编码转换失败: %w", err)
	}

	// #nosec G304 G703 - scriptPath 是程序内部生成的临时脚本路径，非用户输入
	if err := os.WriteFile(scriptPath, scriptBytes, 0600); err != nil {
		_ = os.Remove(newExePath)
		return fmt.Errorf("创建升级脚本失败: %w", err)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("正在启动升级程序...")
	fmt.Println("========================================")

	// 启动升级脚本并退出当前程序
	// 使用 cmd /c 直接执行批处理脚本
	// #nosec G204 - cmd 是系统命令，scriptPath 是程序内部生成的
	cmd := exec.Command("cmd", "/c", scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: false, // 显示窗口，方便调试
	}
	if err := cmd.Start(); err != nil {
		_ = os.Remove(newExePath)
		_ = os.Remove(scriptPath)
		return fmt.Errorf("启动升级脚本失败: %w", err)
	}

	// 退出当前程序，让脚本完成替换
	fmt.Println("程序将在 2 秒后退出以完成升级...")
	time.Sleep(2 * time.Second)
	os.Exit(0)

	return nil
}
