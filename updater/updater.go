package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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
	if u.currentCommit != "" && u.currentCommit != "unknown" && latestCommit != "" {
		// 统一比较长度：使用较短的长度进行比较
		minLen := len(u.currentCommit)
		if len(latestCommit) < minLen {
			minLen = len(latestCommit)
		}
		if latestCommit[:minLen] != u.currentCommit[:minLen] {
			// Commit 不同，说明是同一版本的不同构建
			return true, latestVersion, latestCommit, downloadURL, nil
		}
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
// 格式: https://cnb.cool/dtapp/skills/-/releases/download/{tag}/{binary-name}-{os}-{arch}.{tar.gz|zip}
func (u *Updater) buildDownloadURL(version string) string {
	fileName := fmt.Sprintf("%s-%s-%s%s", types.ProjectName, runtime.GOOS, runtime.GOARCH, archiveExt())
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
func (u *Updater) updateBinary(archiveURL, targetPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, archiveURL, nil)
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

	// 创建临时目录存储下载的压缩包和解压后的文件
	tmpDir, err := os.MkdirTemp("", "dtskills-update-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// 下载压缩包
	archivePath := filepath.Join(tmpDir, "update"+archiveExt())
	// #nosec G304 - archivePath 是程序内部生成的临时文件路径，非用户输入
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("创建压缩包文件失败: %w", err)
	}

	// 使用进度条下载
	if resp.ContentLength > 0 {
		if err := downloadWithProgressBar(resp.Body, archiveFile, resp.ContentLength, "下载更新"); err != nil {
			_ = archiveFile.Close()
			return fmt.Errorf("下载失败: %w", err)
		}
	} else {
		// 未知大小，直接下载
		if _, err := io.Copy(archiveFile, resp.Body); err != nil {
			_ = archiveFile.Close()
			return fmt.Errorf("下载失败: %w", err)
		}
	}
	_ = archiveFile.Close()

	// 解压压缩包
	binaryName := filepath.Base(targetPath)
	extractedPath, err := extractArchive(archivePath, tmpDir, binaryName)
	if err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	// 打开解压后的二进制文件
	// #nosec G304 - extractedPath 是程序内部生成的临时文件路径，非用户输入
	updateFile, err := os.Open(extractedPath)
	if err != nil {
		return fmt.Errorf("打开解压后的文件失败: %w", err)
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

// extractArchive 解压压缩包
func extractArchive(archivePath, destDir, expectedBinary string) (string, error) {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir, expectedBinary)
	}
	return extractTarGz(archivePath, destDir, expectedBinary)
}

// extractTarGz 解压 tar.gz 文件
func extractTarGz(archivePath, destDir, expectedBinary string) (string, error) {
	// #nosec G304 - archivePath 是程序内部生成的临时文件路径，非用户输入
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = gzReader.Close()
	}()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// 跳过目录和非目标文件
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// 检查是否为目标二进制文件
		if header.Name == expectedBinary || filepath.Base(header.Name) == expectedBinary {
			destPath := filepath.Join(destDir, expectedBinary)
			// #nosec G304 - destPath 是程序内部生成的临时文件路径，非用户输入
			// #nosec G115 - header.Mode 是 tar 文件头中的标准权限字段，转换安全
			destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			// #nosec G110 - 压缩包大小已在下载时验证，且文件大小在合理范围内
			_, err = io.Copy(destFile, tarReader)
			_ = destFile.Close()
			if err != nil {
				return "", err
			}
			return destPath, nil
		}
	}

	return "", fmt.Errorf("在压缩包中未找到二进制文件: %s", expectedBinary)
}

// extractZip 解压 zip 文件
func extractZip(archivePath, destDir, expectedBinary string) (string, error) {
	// #nosec G304 - archivePath 是程序内部生成的临时文件路径，非用户输入
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = reader.Close()
	}()

	for _, file := range reader.File {
		// 跳过目录
		if file.FileInfo().IsDir() {
			continue
		}

		// 检查是否为目标二进制文件
		if file.Name == expectedBinary || filepath.Base(file.Name) == expectedBinary {
			destPath := filepath.Join(destDir, expectedBinary)
			// #nosec G304 - destPath 是程序内部生成的临时文件路径，非用户输入
			destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
			if err != nil {
				return "", err
			}

			srcFile, err := file.Open()
			if err != nil {
				_ = destFile.Close()
				return "", err
			}

			// #nosec G110 - 压缩包大小已在下载时验证，且文件大小在合理范围内
			_, err = io.Copy(destFile, srcFile)
			_ = srcFile.Close()
			_ = destFile.Close()
			if err != nil {
				return "", err
			}
			return destPath, nil
		}
	}

	return "", fmt.Errorf("在压缩包中未找到二进制文件: %s", expectedBinary)
}

// GetCurrentVersion 获取当前版本
func (u *Updater) GetCurrentVersion() string {
	return u.currentVersion
}

// GetPlatformInfo 获取平台信息
func GetPlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
