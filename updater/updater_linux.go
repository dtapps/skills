//go:build linux

package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dtapp/skills/types"
)

// InstallType 安装类型
type InstallType int

const (
	InstallTypeBinary InstallType = iota // 二进制安装
	InstallTypeDeb                       // deb 包安装
	InstallTypeRpm                       // rpm 包安装
)

// binaryExt 返回二进制文件扩展名
func binaryExt() string {
	return ""
}

// updateWithPlatform Linux 平台特定的更新逻辑
func (u *Updater) updateWithPlatform(binaryURL, targetPath string) error {
	// 检测安装类型
	installType := detectInstallType(targetPath)

	// 如果是包管理器安装，使用包升级
	if installType == InstallTypeDeb || installType == InstallTypeRpm {
		return u.updateWithPackage(targetPath, installType)
	}

	// 二进制安装，使用标准升级流程
	return u.updateBinary(binaryURL, targetPath)
}

// detectInstallType 检测安装类型
func detectInstallType(targetPath string) InstallType {
	// 检查是否在系统目录
	if !strings.HasPrefix(targetPath, "/usr/bin/") && !strings.HasPrefix(targetPath, "/usr/local/bin/") {
		return InstallTypeBinary
	}

	// 检查是否有 dpkg（Debian/Ubuntu）
	if _, err := exec.LookPath("dpkg"); err == nil {
		// 检查当前包是否通过 dpkg 安装
		// #nosec G204 G702 - dpkg 是系统命令，targetPath 经过验证
		cmd := exec.Command("dpkg", "-S", targetPath)
		if err := cmd.Run(); err == nil {
			return InstallTypeDeb
		}
	}

	// 检查是否有 rpm（RHEL/CentOS/Fedora）
	if _, err := exec.LookPath("rpm"); err == nil {
		// 检查当前包是否通过 rpm 安装
		// #nosec G204 G702 - rpm 是系统命令，targetPath 经过验证
		cmd := exec.Command("rpm", "-qf", targetPath)
		if err := cmd.Run(); err == nil {
			return InstallTypeRpm
		}
	}

	return InstallTypeBinary
}

// buildPackageURL 构建包下载 URL
func buildPackageURL(repoURL, version string, installType InstallType) string {
	switch installType {
	case InstallTypeDeb:
		return fmt.Sprintf("%s/-/releases/download/%s/%s_amd64.deb", repoURL, version, types.ProjectName)
	case InstallTypeRpm:
		return fmt.Sprintf("%s/-/releases/download/%s/%s.x86_64.rpm", repoURL, version, types.ProjectName)
	default:
		return fmt.Sprintf("%s/-/releases/download/%s/%s-linux-amd64%s", repoURL, version, types.ProjectName, binaryExt())
	}
}

// needElevation 检查是否需要提升权限
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
func updateWithElevation(reader io.Reader, targetPath string) error {
	// 下载到临时文件
	tmpFile, err := os.CreateTemp("", "dtskills-update-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	fmt.Printf("下载位置: %s\n", tmpFile.Name())

	if _, err := io.Copy(tmpFile, reader); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 设置执行权限
	// #nosec G302 - 需要 0755 执行权限用于可执行文件
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("设置权限失败: %w", err)
	}

	// 使用 sudo 替换系统文件
	fmt.Println("需要 sudo 权限来更新系统安装的二进制文件...")
	// #nosec G204 G702 - sudo cp 是系统命令，路径经过验证
	cmd := exec.Command("sudo", "cp", tmpFile.Name(), targetPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo 更新失败: %w", err)
	}
	return nil
}

// updateWithPackage 使用包管理器升级
func (u *Updater) updateWithPackage(targetPath string, installType InstallType) error {
	// 获取最新版本
	latestVersion, _, err := u.getLatestVersionAndCommit()
	if err != nil {
		return fmt.Errorf("获取最新版本失败: %w", err)
	}

	packageURL := buildPackageURL(u.repoURL, latestVersion, installType)

	fmt.Printf("正在下载 %s 包...\n", map[InstallType]string{
		InstallTypeDeb: "deb",
		InstallTypeRpm: "rpm",
	}[installType])

	// #nosec G107 - packageURL 是程序内部构造的 URL，非用户输入
	resp, err := http.Get(packageURL)
	if err != nil {
		return fmt.Errorf("下载包失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载包失败，服务器返回: %d", resp.StatusCode)
	}

	// 下载到临时文件
	tmpFile, err := os.CreateTemp("", "dtskills-update-*."+map[InstallType]string{
		InstallTypeDeb: "deb",
		InstallTypeRpm: "rpm",
	}[installType])
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	fmt.Printf("下载位置: %s\n", tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 根据安装类型使用不同的升级方式
	switch installType {
	case InstallTypeDeb:
		return updateDeb(tmpFile.Name())
	case InstallTypeRpm:
		return updateRpm(tmpFile.Name())
	default:
		return fmt.Errorf("未知的安装类型")
	}
}

// updateDeb 使用 dpkg 升级 deb 包
func updateDeb(packagePath string) error {
	fmt.Println("检测到通过 deb 包安装，使用 dpkg 升级...")
	// #nosec G204 - sudo dpkg 是系统命令，packagePath 经过验证
	cmd := exec.Command("sudo", "dpkg", "-i", packagePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dpkg 升级失败: %w", err)
	}
	return nil
}

// updateRpm 使用 rpm 升级 rpm 包
func updateRpm(packagePath string) error {
	fmt.Println("检测到通过 rpm 包安装，使用 rpm 升级...")
	// #nosec G204 G702 - sudo rpm 是系统命令，packagePath 经过验证
	cmd := exec.Command("sudo", "rpm", "-U", packagePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rpm 升级失败: %w", err)
	}
	return nil
}
