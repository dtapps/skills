package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
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

func TestExtractArchive(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 测试 tar.gz 解压（Linux/macOS）
	if archiveExt() == ".tar.gz" {
		testExtractTarGz(t, tmpDir)
	}

	// 测试 zip 解压（Windows）
	if archiveExt() == ".zip" {
		testExtractZip(t, tmpDir)
	}
}

func testExtractTarGz(t *testing.T, tmpDir string) {
	// 创建一个测试用的 tar.gz 文件
	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	if err := createTestTarGz(archivePath); err != nil {
		t.Fatalf("创建测试 tar.gz 失败: %v", err)
	}

	// 测试解压
	extractedPath, err := extractArchive(archivePath, tmpDir, "testbin")
	if err != nil {
		t.Fatalf("解压 tar.gz 失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(extractedPath); os.IsNotExist(err) {
		t.Errorf("解压后的文件不存在: %s", extractedPath)
	}

	t.Logf("Successfully extracted tar.gz to: %s", extractedPath)
}

func testExtractZip(t *testing.T, tmpDir string) {
	// 创建一个测试用的 zip 文件
	archivePath := filepath.Join(tmpDir, "test.zip")
	if err := createTestZip(archivePath); err != nil {
		t.Fatalf("创建测试 zip 失败: %v", err)
	}

	// 测试解压
	extractedPath, err := extractArchive(archivePath, tmpDir, "testbin.exe")
	if err != nil {
		t.Fatalf("解压 zip 失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(extractedPath); os.IsNotExist(err) {
		t.Errorf("解压后的文件不存在: %s", extractedPath)
	}

	t.Logf("Successfully extracted zip to: %s", extractedPath)
}

// createTestTarGz 创建测试用的 tar.gz 文件
func createTestTarGz(archivePath string) error {
	// #nosec G304 - archivePath 是测试生成的临时文件路径，非用户输入
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	gzWriter := gzip.NewWriter(file)
	defer func() {
		_ = gzWriter.Close()
	}()

	tarWriter := tar.NewWriter(gzWriter)
	defer func() {
		_ = tarWriter.Close()
	}()

	// 添加一个测试文件
	content := []byte("test binary content")
	header := &tar.Header{
		Name: "testbin",
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tarWriter.Write(content); err != nil {
		return err
	}

	return nil
}

// createTestZip 创建测试用的 zip 文件
func createTestZip(archivePath string) error {
	// #nosec G304 - archivePath 是测试生成的临时文件路径，非用户输入
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	zipWriter := zip.NewWriter(file)
	defer func() {
		_ = zipWriter.Close()
	}()

	// 添加一个测试文件
	content := []byte("test binary content")
	writer, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   "testbin.exe",
		Method: zip.Deflate,
	})
	if err != nil {
		return err
	}
	if _, err := writer.Write(content); err != nil {
		return err
	}

	return nil
}
