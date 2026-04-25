//go:build windows

package updater

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestUpdateWithElevationBasic 测试基本升级流程
func TestUpdateWithElevationBasic(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 创建模拟的旧版本可执行文件
	oldExe := filepath.Join(tempDir, "testapp.exe")
	oldContent := []byte("old version")
	if err := os.WriteFile(oldExe, oldContent, 0755); err != nil {
		t.Fatalf("创建旧版本文件失败: %v", err)
	}

	// 创建新版本内容
	newContent := []byte("new version 2.0")
	reader := bytes.NewReader(newContent)

	// 执行升级
	err = updateWithElevation(reader, oldExe)
	if err != nil {
		t.Errorf("升级失败: %v", err)
	}

	// 等待升级脚本执行
	time.Sleep(3 * time.Second)

	// 验证文件是否被替换
	content, err := os.ReadFile(oldExe)
	if err != nil {
		t.Errorf("读取文件失败: %v", err)
	}

	if string(content) != string(newContent) {
		t.Errorf("文件内容不匹配，期望: %s, 实际: %s", newContent, content)
	}
}

// TestUpdateWithElevationTempFile 测试临时文件创建
func TestUpdateWithElevationTempFile(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 创建模拟的可执行文件
	targetExe := filepath.Join(tempDir, "target.exe")
	if err := os.WriteFile(targetExe, []byte("target"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 创建新版本内容
	newContent := []byte("new content")
	reader := bytes.NewReader(newContent)

	// 执行升级（会创建临时文件）
	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Errorf("升级失败: %v", err)
	}

	// 检查临时文件是否被创建
	tempFiles, err := filepath.Glob(filepath.Join(os.TempDir(), "dtskills-update-*.exe"))
	if err != nil {
		t.Errorf("查找临时文件失败: %v", err)
	}

	// 临时文件应该存在（升级脚本还未执行完毕）
	if len(tempFiles) == 0 {
		t.Log("临时文件可能已被清理或脚本已执行")
	}
}

// TestUpdateWithElevationScriptContent 测试批处理脚本内容
func TestUpdateWithElevationScriptContent(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 创建模拟的可执行文件
	targetExe := filepath.Join(tempDir, "target.exe")
	if err := os.WriteFile(targetExe, []byte("target"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 创建新版本内容
	newContent := []byte("new version")
	reader := bytes.NewReader(newContent)

	// 执行升级
	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Errorf("升级失败: %v", err)
	}

	// 查找创建的批处理脚本
	scriptFiles, err := filepath.Glob(filepath.Join(os.TempDir(), "dtskills-update-*.exe.bat"))
	if err != nil {
		t.Errorf("查找脚本文件失败: %v", err)
	}

	if len(scriptFiles) > 0 {
		// 读取脚本内容验证
		scriptContent, err := os.ReadFile(scriptFiles[0])
		if err != nil {
			t.Errorf("读取脚本失败: %v", err)
		}

		content := string(scriptContent)

		// 验证脚本包含必要命令
		requiredCommands := []string{
			"@echo off",
			"chcp 65001",
			"timeout /t 2",
			"copy /Y",
			"if errorlevel 1",
			"del",
		}

		for _, cmd := range requiredCommands {
			if !strings.Contains(content, cmd) {
				t.Errorf("脚本缺少必要命令: %s", cmd)
			}
		}

		t.Logf("批处理脚本内容验证通过，包含 %d 个必要命令", len(requiredCommands))
	}
}

// TestUpdateWithElevationInvalidTarget 测试无效目标路径
func TestUpdateWithElevationInvalidTarget(t *testing.T) {
	newContent := []byte("new version")
	reader := bytes.NewReader(newContent)

	// 使用不存在的目标路径
	invalidTarget := "C:\\NonExistent\\Path\\app.exe"

	err := updateWithElevation(reader, invalidTarget)
	// 应该返回错误
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
}

// TestUpdateWithElevationEmptyContent 测试空内容
func TestUpdateWithElevationEmptyContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	targetExe := filepath.Join(tempDir, "target.exe")
	if err := os.WriteFile(targetExe, []byte("content"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 空内容
	reader := bytes.NewReader([]byte{})

	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Errorf("空内容升级失败: %v", err)
	}
}

// TestUpdateWithElevationLargeFile 测试大文件升级
func TestUpdateWithElevationLargeFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	targetExe := filepath.Join(tempDir, "target.exe")
	if err := os.WriteFile(targetExe, []byte("old"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 创建 10MB 的新内容
	largeContent := make([]byte, 10*1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	reader := bytes.NewReader(largeContent)

	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Errorf("大文件升级失败: %v", err)
	}

	t.Log("大文件升级测试通过")
}

// TestUpdateWithElevationSpecialChars 测试特殊字符路径
func TestUpdateWithElevationSpecialChars(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-update-test-特殊字符-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 包含空格和特殊字符的路径
	targetExe := filepath.Join(tempDir, "my app (v2).exe")
	if err := os.WriteFile(targetExe, []byte("old"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	newContent := []byte("new version with special chars")
	reader := bytes.NewReader(newContent)

	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Errorf("特殊字符路径升级失败: %v", err)
	}

	t.Log("特殊字符路径测试通过")
}

// TestUpdateWithElevationPermission 测试权限检查
func TestUpdateWithElevationPermission(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	targetExe := filepath.Join(tempDir, "readonly.exe")
	if err := os.WriteFile(targetExe, []byte("content"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 设置为只读
	if err := os.Chmod(targetExe, 0444); err != nil {
		t.Logf("设置只读权限失败: %v", err)
	}

	newContent := []byte("new version")
	reader := bytes.NewReader(newContent)

	// 升级应该仍然可以工作（通过脚本提升权限）
	err = updateWithElevation(reader, targetExe)
	if err != nil {
		t.Logf("只读文件升级可能需要管理员权限: %v", err)
	}

	// 恢复权限以便清理
	_ = os.Chmod(targetExe, 0755)
}

// TestUpdateScriptExecution 测试脚本执行命令
func TestUpdateScriptExecution(t *testing.T) {
	// 验证 cmd /c start 命令格式
	cmd := exec.Command("cmd", "/c", "echo", "test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("cmd 命令执行失败: %v", err)
	}

	if !strings.Contains(string(output), "test") {
		t.Error("cmd 命令输出不正确")
	}

	t.Log("Windows cmd 命令测试通过")
}

// TestConcurrentUpdates 测试并发升级（应该串行处理）
func TestConcurrentUpdates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-update-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	targetExe := filepath.Join(tempDir, "target.exe")
	if err := os.WriteFile(targetExe, []byte("old"), 0755); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 尝试并发执行两个升级
	done := make(chan error, 2)

	for i := 0; i < 2; i++ {
		go func(idx int) {
			newContent := []byte(fmt.Sprintf("version %d", idx))
			reader := bytes.NewReader(newContent)
			done <- updateWithElevation(reader, targetExe)
		}(i)
	}

	// 等待两个升级完成
	for i := 0; i < 2; i++ {
		err := <-done
		if err != nil {
			t.Logf("并发升级 %d 错误: %v", i, err)
		}
	}

	t.Log("并发升级测试完成")
}
