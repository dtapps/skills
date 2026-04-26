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
	currentCommit  string
	repoURL        string
	httpClient     *http.Client
}

// New 创建更新器
// currentVersion: 当前版本号
// currentCommit: 当前 Git Commit ID
// repoURL: 仓库地址，如 https://cnb.cool/dtapp/skills
func New(currentVersion, currentCommit, repoURL string) *Updater {
	if repoURL == "" {
		// 默认使用 cnb.cool 仓库
		repoURL = "https://cnb.cool/dtapp/skills"
	}

	return &Updater{
		currentVersion: currentVersion,
		currentCommit:  currentCommit,
		repoURL:        repoURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Check 检查是否有新版本
// force: 是否强制升级（即使版本号相同）
// 返回: 是否有新版本, 最新版本号, 最新Commit, 下载URL, 错误
func (u *Updater) Check(force bool) (bool, string, string, string, error) {
	// 从 release 页面获取最新版本和 Commit
	latestVersion, latestCommit, err := u.getLatestVersionAndCommit()
	if err != nil {
		return false, "", "", "", err
	}

	// 构建下载 URL
	downloadURL := u.buildDownloadURL(latestVersion)

	// 如果是 dev 版本，允许升级
	if u.currentVersion == "dev" {
		return true, latestVersion, latestCommit, downloadURL, nil
	}

	// 如果强制升级，直接返回
	if force {
		return true, latestVersion, latestCommit, downloadURL, nil
	}

	// 比较版本号
	if latestVersion != u.currentVersion {
		return true, latestVersion, latestCommit, downloadURL, nil
	}

	// 版本号相同，比较 Commit ID（如果都有的话）
	if u.currentCommit != "" && u.currentCommit != "unknown" &&
		latestCommit != "" && latestCommit != u.currentCommit {
		// Commit 不同，说明是同一版本的不同构建
		return true, latestVersion, latestCommit, downloadURL, nil
	}

	// 版本和 Commit 都相同，不需要升级
	return false, latestVersion, latestCommit, downloadURL, nil
}

// getLatestVersionAndCommit 从 CNB release 页面获取最新版本号和 Commit ID
func (u *Updater) getLatestVersionAndCommit() (string, string, error) {
	releasesURL := u.repoURL + "/-/releases"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("无法连接到 CNB: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("服务器返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	bodyStr := string(body)

	// 匹配 CNB HTML 里的标签路径，例如: /dtapp/skills/-/releases/tag/v1.0.1
	versionRe := regexp.MustCompile(`releases/tag/(v\d+\.\d+\.\d+)`)
	versionMatch := versionRe.FindStringSubmatch(bodyStr)
	if len(versionMatch) < 2 {
		return "", "", fmt.Errorf("无法从 CNB 页面解析最新版本号，请检查仓库是否公开")
	}
	latestVersion := versionMatch[1]

	// 尝试匹配 Commit ID（CNB release 页面通常会显示）
	// 匹配格式: commit/abc123 或 data-commit="abc123"
	commitRe := regexp.MustCompile(`commit/([a-f0-9]{7,40})`)
	commitMatch := commitRe.FindStringSubmatch(bodyStr)
	latestCommit := ""
	if len(commitMatch) >= 2 {
		latestCommit = commitMatch[1]
	}

	return latestVersion, latestCommit, nil
}

// buildDownloadURL 构建下载 URL
// 格式: https://cnb.cool/dtapp/skills/-/releases/download/{tag}/{binary-name}-{os}-{arch}
func (u *Updater) buildDownloadURL(version string) string {
	fileName := fmt.Sprintf("%s-%s-%s%s", types.ProjectName, runtime.GOOS, runtime.GOARCH, binaryExt())
	return fmt.Sprintf("%s/-/releases/download/%s/%s", u.repoURL, version, fileName)
}

// Update 执行更新
// binaryURL: 二进制文件下载地址
func (u *Updater) Update(binaryURL string) error {
	// 获取当前二进制路径
	targetPath, err := os.Executable()
	if err != nil {
		targetPath = os.Args[0]
	}

	// 平台特定的更新逻辑
	return u.updateWithPlatform(binaryURL, targetPath)
}

// updateBinary 使用二进制方式升级
func (u *Updater) updateBinary(binaryURL, targetPath string) error {
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

	// 创建临时文件存储下载内容
	tmpFile, err := os.CreateTemp("", "dtskills-download-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	}()

	// 使用进度条下载
	if resp.ContentLength > 0 {
		if err := downloadWithProgressBar(resp.Body, tmpFile, resp.ContentLength, "下载更新"); err != nil {
			return fmt.Errorf("下载失败: %w", err)
		}
	} else {
		// 未知大小，直接下载
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return fmt.Errorf("下载失败: %w", err)
		}
	}

	// 关闭临时文件以便读取
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 重新打开临时文件用于更新
	// #nosec G304 - tmpPath 是程序内部生成的临时文件路径，非用户输入
	updateFile, err := os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("打开临时文件失败: %w", err)
	}
	defer func() {
		_ = updateFile.Close()
	}()

	// 检查是否需要提升权限（系统目录如 /usr/bin 需要 root 权限）
	if needElevation(targetPath) {
		return updateWithElevation(updateFile, targetPath)
	}

	// 使用 selfupdate 应用更新
	opts := selfupdate.Options{
		TargetPath: targetPath,
	}

	if err := selfupdate.Apply(updateFile, opts); err != nil {
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
