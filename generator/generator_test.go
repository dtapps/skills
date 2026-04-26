package generator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dtapp/skills/types"
)

func TestContentsEqualIgnoreDate(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		new      string
		expected bool
	}{
		{
			name: "相同内容，不同 date",
			existing: `---
name: "test"
description: "test desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 10:00:00"
---

# Test`,
			new: `---
name: "test"
description: "test desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 14:00:00"
---

# Test`,
			expected: true,
		},
		{
			name: "不同内容",
			existing: `---
name: "test"
description: "test desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 10:00:00"
---

# Test`,
			new: `---
name: "test"
description: "different desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 14:00:00"
---

# Test`,
			expected: false,
		},
		{
			name: "相同内容，相同 date",
			existing: `---
name: "test"
description: "test desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 10:00:00"
---

# Test`,
			new: `---
name: "test"
description: "test desc"
metadata:
  version: "1.0.0"
  author: "test"
  date: "2026-04-26 10:00:00"
---

# Test`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contentsEqualIgnoreDate(tt.existing, tt.new)
			if result != tt.expected {
				t.Errorf("contentsEqualIgnoreDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateSkillMd_SkipWhenContentUnchanged(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skill-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	skillInfo := types.SkillInfo{
		Name:        "test-skill",
		Version:     "1.0.0",
		Author:      "test",
		Description: "Test skill description",
	}

	// 第一次生成
	err = generateSkillMd(skillInfo, tempDir, "dtskills")
	if err != nil {
		t.Fatalf("第一次生成失败: %v", err)
	}

	// 获取第一次生成后的文件信息
	skillMdPath := filepath.Join(tempDir, "SKILL.md")
	firstStat, err := os.Stat(skillMdPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}
	firstModTime := firstStat.ModTime()

	// 等待一小段时间确保时间戳不同
	time.Sleep(100 * time.Millisecond)

	// 第二次生成（内容相同）
	err = generateSkillMd(skillInfo, tempDir, "dtskills")
	if err != nil {
		t.Fatalf("第二次生成失败: %v", err)
	}

	// 获取第二次生成后的文件信息
	secondStat, err := os.Stat(skillMdPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}
	secondModTime := secondStat.ModTime()

	// 修改时间应该相同（文件没有被重新写入）
	if !firstModTime.Equal(secondModTime) {
		t.Errorf("文件被重复写入了，修改时间发生变化: %v -> %v", firstModTime, secondModTime)
	}
}

func TestGenerateSkillMd_UpdateWhenContentChanged(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skill-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	skillInfo := types.SkillInfo{
		Name:        "test-skill",
		Version:     "1.0.0",
		Author:      "test",
		Description: "Test skill description",
	}

	// 第一次生成
	err = generateSkillMd(skillInfo, tempDir, "dtskills")
	if err != nil {
		t.Fatalf("第一次生成失败: %v", err)
	}

	// 获取第一次生成后的文件信息
	skillMdPath := filepath.Join(tempDir, "SKILL.md")
	firstStat, err := os.Stat(skillMdPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}
	firstModTime := firstStat.ModTime()

	// 等待一小段时间确保时间戳不同
	time.Sleep(100 * time.Millisecond)

	// 修改 skill 信息
	skillInfo.Description = "Changed description"

	// 第二次生成（内容不同）
	err = generateSkillMd(skillInfo, tempDir, "dtskills")
	if err != nil {
		t.Fatalf("第二次生成失败: %v", err)
	}

	// 获取第二次生成后的文件信息
	secondStat, err := os.Stat(skillMdPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}
	secondModTime := secondStat.ModTime()

	// 修改时间应该不同（文件被重新写入了）
	if !secondModTime.After(firstModTime) {
		t.Errorf("文件应该被更新，但修改时间没有变化: %v -> %v", firstModTime, secondModTime)
	}
}

func TestGenerateSkillOrigin_SkipWhenFileExists(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skill-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	originPath := filepath.Join(tempDir, ".skill-origin.json")

	// 创建已存在的 origin 文件
	originalContent := `{
  "version": 1,
  "origin_kind": "manual",
  "installed_at": 1234567890000
}`
	err = os.WriteFile(originPath, []byte(originalContent), 0600)
	if err != nil {
		t.Fatalf("创建原始文件失败: %v", err)
	}

	// 调用生成函数
	err = generateSkillOrigin(tempDir)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	// 读取文件内容
	// #nosec G304 - originPath 是测试中的临时文件路径
	content, err := os.ReadFile(originPath)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	// 内容应该保持不变
	if string(content) != originalContent {
		t.Errorf("文件内容被修改了，应该保持原有 installed_at\n原始: %s\n现在: %s", originalContent, string(content))
	}
}

func TestGenerateSkillOrigin_CreateWhenFileNotExists(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skill-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	originPath := filepath.Join(tempDir, ".skill-origin.json")

	// 调用生成函数
	err = generateSkillOrigin(tempDir)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	// 检查文件是否存在
	_, err = os.Stat(originPath)
	if err != nil {
		t.Errorf("文件应该被创建，但没有找到: %v", err)
	}
}
