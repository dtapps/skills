package updater

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		currentCommit  string
		repoURL        string
		wantRepoURL    string
	}{
		{
			name:           "默认仓库地址",
			currentVersion: "v1.0.0",
			currentCommit:  "abc123",
			repoURL:        "",
			wantRepoURL:    "https://cnb.cool/dtapp/skills",
		},
		{
			name:           "自定义仓库地址",
			currentVersion: "v1.0.0",
			currentCommit:  "abc123",
			repoURL:        "https://example.com/repo",
			wantRepoURL:    "https://example.com/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.currentVersion, tt.currentCommit, tt.repoURL)
			if u.repoURL != tt.wantRepoURL {
				t.Errorf("New() repoURL = %v, want %v", u.repoURL, tt.wantRepoURL)
			}
			if u.currentVersion != tt.currentVersion {
				t.Errorf("New() currentVersion = %v, want %v", u.currentVersion, tt.currentVersion)
			}
		})
	}
}

func TestGetCurrentVersion(t *testing.T) {
	u := New("v1.2.3", "abc123", "")
	if got := u.GetCurrentVersion(); got != "v1.2.3" {
		t.Errorf("GetCurrentVersion() = %v, want %v", got, "v1.2.3")
	}
}

func TestGetPlatformInfo(t *testing.T) {
	info := GetPlatformInfo()
	if info == "" {
		t.Error("GetPlatformInfo() returned empty string")
	}
	t.Logf("Platform info: %s", info)
}

func TestBuildDownloadURL(t *testing.T) {
	u := New("v1.0.0", "abc123", "https://cnb.cool/dtapp/skills")

	tests := []struct {
		name    string
		version string
		wantExt string
	}{
		{
			name:    "构建下载URL",
			version: "v1.0.1",
			wantExt: binaryExt(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := u.buildDownloadURL(tt.version)
			if url == "" {
				t.Error("buildDownloadURL() returned empty string")
			}
			t.Logf("Download URL: %s", url)
		})
	}
}

func TestBinaryExt(t *testing.T) {
	ext := binaryExt()
	t.Logf("Binary extension: %q", ext)

	// 根据平台验证扩展名
	switch ext {
	case "":
		t.Log("Platform: Linux or macOS")
	case ".exe":
		t.Log("Platform: Windows")
	default:
		t.Errorf("Unexpected binary extension: %q", ext)
	}
}
