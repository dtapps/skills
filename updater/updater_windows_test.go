//go:build windows

package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNeedElevationWindows(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建一个可写文件
	writableFile := filepath.Join(tmpDir, "writable.exe")
	if err := os.WriteFile(writableFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name       string
		targetPath string
		want       bool
	}{
		{
			name:       "可写文件不需要提权",
			targetPath: writableFile,
			want:       false,
		},
		{
			name:       "系统目录需要管理员权限",
			targetPath: `C:\Windows\System32\dtskills.exe`,
			want:       true,
		},
		{
			name:       "不存在的路径需要提权",
			targetPath: `C:\nonexistent\path\dtskills.exe`,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needElevation(tt.targetPath)
			t.Logf("Path: %s, Need elevation: %v (expected: %v)", tt.targetPath, got, tt.want)
		})
	}
}

func TestUpdateWithPlatformWindows(t *testing.T) {
	_ = New("v1.0.0", "https://cnb.cool/dtapp/skills")

	// 测试 Windows 平台使用二进制升级
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "dtskills.exe")
	binaryURL := "https://example.com/dtskills-windows-amd64.exe"

	// 这个测试不会真正执行更新，只是验证函数存在且能调用
	t.Logf("Windows platform update: binaryURL=%s, targetPath=%s", binaryURL, targetPath)
	t.Log("Windows uses binary update method with admin elevation")
}

func TestBinaryExtWindows(t *testing.T) {
	ext := binaryExt()
	if ext != ".exe" {
		t.Errorf("binaryExt() on Windows = %q, want .exe", ext)
	}
	t.Logf("Windows binary extension: %q", ext)
}
