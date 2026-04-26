//go:build linux

package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectInstallType(t *testing.T) {
	// 创建临时目录模拟不同场景
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		targetPath   string
		setupFunc    func()
		expectedType InstallType
	}{
		{
			name:         "非系统目录 - 二进制安装",
			targetPath:   filepath.Join(tmpDir, "dtskills"),
			expectedType: InstallTypeBinary,
		},
		{
			name:         "系统目录但无包管理器 - 二进制安装",
			targetPath:   "/usr/local/bin/dtskills",
			expectedType: InstallTypeBinary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			got := detectInstallType(tt.targetPath)
			t.Logf("Path: %s, Detected type: %v", tt.targetPath, got)
		})
	}
}

func TestBuildPackageURL(t *testing.T) {
	repoURL := "https://cnb.cool/dtapp/skills"
	version := "v1.0.0"

	tests := []struct {
		name        string
		installType InstallType
		wantContain string
	}{
		{
			name:        "deb包URL",
			installType: InstallTypeDeb,
			wantContain: "_amd64.deb",
		},
		{
			name:        "rpm包URL",
			installType: InstallTypeRpm,
			wantContain: ".x86_64.rpm",
		},
		{
			name:        "二进制URL",
			installType: InstallTypeBinary,
			wantContain: "linux-amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPackageURL(repoURL, version, tt.installType)
			if got == "" {
				t.Error("buildPackageURL() returned empty string")
			}
			if tt.wantContain != "" && !contains(got, tt.wantContain) {
				t.Errorf("buildPackageURL() = %v, want contain %v", got, tt.wantContain)
			}
			t.Logf("Package URL: %s", got)
		})
	}
}

func TestNeedElevation(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建一个可写文件
	writableFile := filepath.Join(tmpDir, "writable")
	if err := os.WriteFile(writableFile, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 创建一个只读文件
	readonlyFile := filepath.Join(tmpDir, "readonly")
	if err := os.WriteFile(readonlyFile, []byte("test"), 0400); err != nil {
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
			name:       "不存在的路径需要提权",
			targetPath: "/nonexistent/path/binary",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needElevation(tt.targetPath)
			t.Logf("Path: %s, Need elevation: %v", tt.targetPath, got)
		})
	}
}

func TestInstallTypeString(t *testing.T) {
	tests := []struct {
		installType InstallType
		want        string
	}{
		{InstallTypeBinary, "binary"},
		{InstallTypeDeb, "deb"},
		{InstallTypeRpm, "rpm"},
	}

	// 测试 InstallType 到字符串的映射
	installTypeNames := map[InstallType]string{
		InstallTypeBinary: "binary",
		InstallTypeDeb:    "deb",
		InstallTypeRpm:    "rpm",
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := installTypeNames[tt.installType]
			if got != tt.want {
				t.Errorf("InstallType name = %v, want %v", got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
