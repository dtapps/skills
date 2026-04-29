package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// ExternalConfigManager 外部配置文件管理器
// 用于读取通过环境变量指定的目录下的所有 JSON 配置文件
type ExternalConfigManager struct {
	v       *viper.Viper
	envVar  string
	loaded  bool
	mu      sync.Mutex
	sources []string // 记录加载的配置文件路径
}

// NewExternalConfigManager 创建外部配置管理器
// envVar: 环境变量名，该环境变量的值应该是包含 JSON 配置文件的目录路径
// 例如：WEIBO_TOKEN_DIR=/path/to/tokens/
// 该目录下的所有 .json 文件都会被读取
func NewExternalConfigManager(envVar string) *ExternalConfigManager {
	return &ExternalConfigManager{
		v:       viper.New(),
		envVar:  envVar,
		loaded:  false,
		sources: make([]string, 0),
	}
}

// Load 加载外部配置文件
// 读取指定目录下的所有 .json 文件
// 如果环境变量未设置或目录不存在，返回错误
func (ecm *ExternalConfigManager) Load() error {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if ecm.loaded {
		return nil
	}

	dirPath := os.Getenv(ecm.envVar)
	if dirPath == "" {
		return fmt.Errorf("环境变量 %s 未设置", ecm.envVar)
	}

	// 检查目录是否存在
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("配置目录不存在: %s", dirPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("配置路径不是目录: %s", dirPath)
	}

	// 读取目录下的所有 .json 文件
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("读取配置目录失败: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}

		configPath := filepath.Join(dirPath, entry.Name())
		ecm.v.SetConfigFile(configPath)
		ecm.v.SetConfigType("json")

		// 读取配置（合并多个 JSON 文件）
		if err := ecm.v.MergeInConfig(); err != nil {
			return fmt.Errorf("读取配置文件 %s 失败: %w", entry.Name(), err)
		}

		ecm.sources = append(ecm.sources, configPath)
	}

	if len(ecm.sources) == 0 {
		return fmt.Errorf("配置目录 %s 下未找到 JSON 文件", dirPath)
	}

	ecm.loaded = true
	return nil
}

// Get 获取配置值
func (ecm *ExternalConfigManager) Get(key string) any {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return nil
	}
	return ecm.v.Get(key)
}

// GetString 获取字符串配置值
func (ecm *ExternalConfigManager) GetString(key string) string {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return ""
	}
	return ecm.v.GetString(key)
}

// GetStringSlice 获取字符串数组配置值
func (ecm *ExternalConfigManager) GetStringSlice(key string) []string {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return nil
	}
	return ecm.v.GetStringSlice(key)
}

// GetInt 获取整数配置值
func (ecm *ExternalConfigManager) GetInt(key string) int {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return 0
	}
	return ecm.v.GetInt(key)
}

// GetBool 获取布尔配置值
func (ecm *ExternalConfigManager) GetBool(key string) bool {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return false
	}
	return ecm.v.GetBool(key)
}

// GetAll 获取所有配置
func (ecm *ExternalConfigManager) GetAll() map[string]any {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.loaded {
		return nil
	}
	return ecm.v.AllSettings()
}

// IsLoaded 检查配置是否已加载
func (ecm *ExternalConfigManager) IsLoaded() bool {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()
	return ecm.loaded
}

// GetDirPath 获取配置目录路径（从环境变量读取）
func (ecm *ExternalConfigManager) GetDirPath() string {
	return os.Getenv(ecm.envVar)
}

// GetSources 获取已加载的配置文件路径列表
func (ecm *ExternalConfigManager) GetSources() []string {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	sources := make([]string, len(ecm.sources))
	copy(sources, ecm.sources)
	return sources
}
