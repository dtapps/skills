//go:build windows

package updater

import (
	"fmt"
	"io"
	"os"
)

// needElevation 检查是否需要提升权限（管理员权限）
func needElevation(targetPath string) bool {
	// 尝试写入文件测试权限
	if _, err := os.OpenFile(targetPath, os.O_WRONLY, 0); err != nil {
		return true
	}
	return false
}

// updateWithElevation 使用提升的权限更新
func updateWithElevation(reader io.Reader, targetPath string) error {
	// Windows 系统需要管理员权限才能更新系统目录的文件
	fmt.Println("需要管理员权限来更新系统安装的二进制文件...")
	fmt.Println("请右键以管理员身份运行此命令，或手动下载替换文件")
	return fmt.Errorf("Windows 系统请以管理员身份运行升级命令")
}
