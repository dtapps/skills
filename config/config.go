package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dtapp/skills/types"

	"github.com/spf13/viper"
)

// resolveWorkDir 解析工作目录（不依赖 viper 状态）
// 优先级：DTSKILLS_WORK_DIR > 用户配置目录 > 当前目录
func resolveWorkDir() string {
	// 检查环境变量
	workDir := os.Getenv(strings.ToUpper(types.ProjectName) + "_WORK_DIR")
	if workDir != "" {
		return workDir
	}

	// 使用用户配置目录
	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, types.ProjectName)
	}

	// 默认当前目录
	return "."
}

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 配置文件名称
	viper.SetConfigName("config")
	// 配置文件类型
	// viper.SetConfigType("yaml")

	// 获取工作目录
	workDir := resolveWorkDir()

	// 添加工作目录到配置搜索路径
	viper.AddConfigPath(workDir)

	// 自动读取环境变量
	viper.AutomaticEnv()
	// 环境变量前缀
	viper.SetEnvPrefix(strings.ToUpper(types.ProjectName))
	// 绑定环境变量
	_ = viper.BindEnv("api_keys.baidu", "BAIDU_API_KEY")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
		// 配置文件未找到，静默使用默认配置
	}
	// 配置文件加载成功，静默处理

	return nil
}

// SetConfigFile 设置配置文件路径
func SetConfigFile(path string) {
	viper.SetConfigFile(path)
}

// GetBaiduAPIKey 获取百度 API Key
func GetBaiduAPIKey() string {
	return viper.GetString("api_keys.baidu")
}

// GetSkillsPath 获取技能包路径
func GetSkillsPath() string {
	return viper.GetString("skills.path")
}

// SetBaiduAPIKey 设置百度 API Key
func SetBaiduAPIKey(apiKey string) {
	viper.Set("api_keys.baidu", apiKey)
}

// SetSkillsPath 设置技能包路径
func SetSkillsPath(path string) {
	viper.Set("skills.path", path)
}

// SaveConfig 保存配置到文件
func SaveConfig() error {
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	fmt.Printf("配置已保存到: %s\n", viper.ConfigFileUsed())
	return nil
}

// SaveConfigAs 保存配置到指定文件
func SaveConfigAs(path string) error {
	err := viper.WriteConfigAs(path)
	if err != nil {
		return err
	}
	fmt.Printf("配置已保存到: %s\n", path)
	return nil
}

// InitConfig 初始化配置文件
// 如果配置文件不存在，则创建默认配置文件
func InitConfig() (string, error) {
	// 检查是否已加载配置文件
	if viper.ConfigFileUsed() != "" {
		return viper.ConfigFileUsed(), nil
	}

	// 使用 resolveWorkDir 获取工作目录
	workDir := resolveWorkDir()

	// 确保工作目录存在
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", fmt.Errorf("创建工作目录失败: %w", err)
	}

	configPath := filepath.Join(workDir, "config.yaml")
	viper.SetConfigFile(configPath)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return "", fmt.Errorf("创建配置文件失败: %w", err)
	}

	return configPath, nil
}

// Get 获取任意配置项
func Get(key string) any {
	return viper.Get(key)
}

// Set 设置任意配置项
func Set(key string, value any) error {
	viper.Set(key, value)
	return nil
}

// All 获取所有配置
func All() map[string]any {
	return viper.AllSettings()
}

// GetConfigFile 获取当前使用的配置文件路径
func GetConfigFile() string {
	return viper.ConfigFileUsed()
}

// GetWorkDir 获取工作目录
// 优先级：DTSKILLS_WORK_DIR > 配置文件所在目录 > 用户配置目录 > 当前目录
func GetWorkDir() string {
	// 检查环境变量
	workDir := os.Getenv(strings.ToUpper(types.ProjectName) + "_WORK_DIR")
	if workDir != "" {
		return workDir
	}

	// 使用配置文件所在目录
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		return filepath.Dir(configFile)
	}

	// 使用用户配置目录
	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, types.ProjectName)
	}

	// 默认当前目录
	return "."
}
