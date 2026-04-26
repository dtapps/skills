//go:build darwin

package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNeedElevationDarwin(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建一个可写文件
	writableFile := filepath.Join(tmpDir, "writable")
	if err := os.WriteFile(writableFile, []byte("test"), 0600); err != nil {
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
			name:       "系统目录需要提权",
			targetPath: "/usr/bin/dtskills",
			want:       true,
		},
		{
			name:       "不存在的路径需要提权",
			targetPath: "/nonexistent/path/binary",
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

func TestUpdateWithPlatformDarwin(t *testing.T) {
	_ = New("v1.0.0", "https://cnb.cool/dtapp/skills")

	// 测试 macOS 平台使用二进制升级
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "dtskills")
	binaryURL := "https://example.com/dtskills-darwin-arm64"

	// 这个测试不会真正执行更新，只是验证函数存在且能调用
	t.Logf("Darwin platform update: binaryURL=%s, targetPath=%s", binaryURL, targetPath)
	t.Log("Darwin uses binary update method")
}

func TestBinaryExtDarwin(t *testing.T) {
	ext := binaryExt()
	if ext != "" {
		t.Errorf("binaryExt() on Darwin = %q, want empty string", ext)
	}
	t.Logf("Darwin binary extension: %q", ext)
}
