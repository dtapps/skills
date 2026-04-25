package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/dtapp/skills/types"
	"github.com/minio/selfupdate"
)

// Updater 更新器
type Updater struct {
	currentVersion string
	repoURL        string
	httpClient     *http.Client
}

// New 创建更新器
// repoURL: 仓库地址，如 https://cnb.cool/dtapp/skills
func New(currentVersion, repoURL string) *Updater {
	if repoURL == "" {
		// 默认使用 cnb.cool 仓库
		repoURL = "https://cnb.cool/dtapp/skills"
	}

	return &Updater{
		currentVersion: currentVersion,
		repoURL:        repoURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Check 检查是否有新版本
// 返回: 是否有新版本, 最新版本号, 下载URL, 错误
func (u *Updater) Check() (bool, string, string, error) {
	// 从 release 页面获取最新版本
	latestVersion, err := u.getLatestVersion()
	if err != nil {
		return false, "", "", err
	}

	// 比较版本
	if latestVersion == u.currentVersion && u.currentVersion != "dev" {
		return false, latestVersion, "", nil
	}

	// 构建下载 URL
	downloadURL := u.buildDownloadURL(latestVersion)

	return true, latestVersion, downloadURL, nil
}

// getLatestVersion 从 CNB release 页面获取最新版本号
func (u *Updater) getLatestVersion() (string, error) {
	releasesURL := u.repoURL + "/-/releases"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("无法连接到 CNB: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("服务器返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 匹配 CNB HTML 里的标签路径，例如: /dtapp/skills/-/releases/tag/v1.0.1
	re := regexp.MustCompile(`releases/tag/(v\d+\.\d+\.\d+)`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return "", fmt.Errorf("无法从 CNB 页面解析最新版本号，请检查仓库是否公开")
	}

	return match[1], nil
}

// buildDownloadURL 构建下载 URL
// 格式: https://cnb.cool/dtapp/skills/-/releases/download/{tag}/{binary-name}-{os}-{arch}
func (u *Updater) buildDownloadURL(version string) string {
	fileName := fmt.Sprintf("%s-%s-%s", types.ProjectName, runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf("%s/-/releases/download/%s/%s", u.repoURL, version, fileName)
}

// Update 执行更新
// binaryURL: 二进制文件下载地址
func (u *Updater) Update(binaryURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, binaryURL, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，服务器返回: %d", resp.StatusCode)
	}

	// 获取当前二进制路径
	targetPath, err := os.Executable()
	if err != nil {
		targetPath = os.Args[0]
	}

	// 检查是否需要提升权限（系统目录如 /usr/bin 需要 root 权限）
	if needElevation(targetPath) {
		return updateWithElevation(resp.Body, targetPath)
	}

	// 使用 selfupdate 应用更新
	opts := selfupdate.Options{
		TargetPath: targetPath,
	}

	if err := selfupdate.Apply(resp.Body, opts); err != nil {
		return fmt.Errorf("应用更新失败: %w", err)
	}

	return nil
}

// GetCurrentVersion 获取当前版本
func (u *Updater) GetCurrentVersion() string {
	return u.currentVersion
}

// GetPlatformInfo 获取平台信息
func GetPlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
