package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dtapp/skills/types"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 设置默认值
	viper.SetDefault("skills.path", "../")

	// 配置文件名称
	viper.SetConfigName("config")
	// 配置文件类型
	viper.SetConfigType("yaml")

	// 检查是否通过环境变量指定了配置文件路径
	configPath := os.Getenv(strings.ToUpper(types.ProjectName) + "_CONFIG")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 配置文件搜索路径
		// 当前目录
		viper.AddConfigPath(".")

		// 跨平台标准配置目录
		configDir, err := os.UserConfigDir()
		if err == nil {
			projectConfigDir := filepath.Join(configDir, types.ProjectName)
			viper.AddConfigPath(projectConfigDir)
		}
	}

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
