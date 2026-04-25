package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/dtapp/skills/types"
	"github.com/minio/selfupdate"
)

// VersionInfo 版本信息
type VersionInfo struct {
	Version string `json:"version"`
	URL     string `json:"url,omitempty"`
}

// Updater 更新器
type Updater struct {
	currentVersion string
	updateURL      string
	httpClient     *http.Client
}

// New 创建更新器
// updateURL: 版本检查 URL，返回 JSON 格式 {"version": "x.x.x", "url": "下载链接"}
func New(currentVersion, updateURL string) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		updateURL:      updateURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Check 检查是否有新版本
// 返回: 是否有新版本, 最新版本号, 下载URL, 错误
func (u *Updater) Check() (bool, string, string, error) {
	if u.updateURL == "" {
		return false, "", "", fmt.Errorf("更新 URL 未配置")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.updateURL, nil)
	if err != nil {
		return false, "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return false, "", "", fmt.Errorf("请求失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return false, "", "", fmt.Errorf("服务器返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return false, "", "", fmt.Errorf("解析版本信息失败: %w", err)
	}

	if info.Version == "" {
		return false, "", "", fmt.Errorf("版本信息为空")
	}

	// 比较版本
	if info.Version == u.currentVersion {
		return false, info.Version, "", nil
	}

	return true, info.Version, info.URL, nil
}

// Update 执行更新
// binaryURL: 二进制文件下载地址
func (u *Updater) Update(binaryURL string) error {
	if binaryURL == "" {
		// 如果没有提供 URL，构建默认 URL
		binaryURL = u.buildDefaultURL()
	}

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

	// 使用 selfupdate 应用更新
	opts := selfupdate.Options{
		TargetPath: os.Args[0],
	}

	if err := selfupdate.Apply(resp.Body, opts); err != nil {
		return fmt.Errorf("应用更新失败: %w", err)
	}

	return nil
}

// buildDefaultURL 构建默认下载 URL
// 格式: https://cnb.cool/dtapp/skills/-/releases/latest/download/{binary-name}-{os}-{arch}
func (u *Updater) buildDefaultURL() string {
	baseURL := fmt.Sprintf("https://cnb.cool/dtapp/skills/-/releases/latest/download/%s", types.ProjectName)
	return fmt.Sprintf("%s-%s-%s", baseURL, runtime.GOOS, runtime.GOARCH)
}

// GetCurrentVersion 获取当前版本
func (u *Updater) GetCurrentVersion() string {
	return u.currentVersion
}

// GetPlatformInfo 获取平台信息
func GetPlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
