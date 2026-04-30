package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dtapp/skills/types"

	"github.com/spf13/viper"
)

//nolint:gosec // 环境变量名不是凭证
const (
	// BaiduAPIKeyEnv 百度 API Key 环境变量
	BaiduAPIKeyEnv = "BAIDU_API_KEY"
	// TencentNewsAPIKeyEnv 腾讯新闻 API Key 环境变量
	TencentNewsAPIKeyEnv = "TENCENT_NEWS_API_KEY"
	// WeiboAPITokenEnv 微博API Token 环境变量
	WeiboAPITokenEnv = "WEIBO_API_TOKEN"
	// WeiboTokenDirEnv 微博 Token 目录环境变量
	WeiboTokenDirEnv = "WEIBO_TOKEN_DIR"
)

// GetWorkDirEnv 获取工作目录环境变量名 (动态生成: {PROJECT_NAME}_WORK_DIR)
func GetWorkDirEnv() string {
	return strings.ToUpper(types.ProjectName) + "_WORK_DIR"
}

//nolint:gosec // 配置键名不是凭证
const (
	// KeyBaiduAPIKey 百度 API Key 配置键
	KeyBaiduAPIKey = "api_keys.baidu"
	// KeyTencentNewsAPIKey 腾讯新闻 API Key 配置键
	KeyTencentNewsAPIKey = "api_keys.tencent_news"
	// KeyWeiboAPIToken 微博 API Token 配置键
	KeyWeiboAPIToken = "api_tokens.weibo"
	// KeyRestyDebug Resty Debug 配置键
	KeyRestyDebug = "debug.resty"
	// KeySkillsPath 技能包路径配置键
	KeySkillsPath = "skills.path"
)

// resolveWorkDir 解析工作目录（不依赖 viper 状态）
// 优先级：DTSKILLS_WORK_DIR > 用户配置目录 > 当前目录
func resolveWorkDir() string {
	// 检查环境变量
	workDir := os.Getenv(GetWorkDirEnv())
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
	// viper 的优先级：环境变量 > 配置文件 > 默认值
	// 所以如果同时设置了环境变量和配置文件，环境变量会覆盖配置文件的值
	viper.AutomaticEnv()
	// 环境变量前缀
	viper.SetEnvPrefix(strings.ToUpper(types.ProjectName))
	// 绑定环境变量
	_ = viper.BindEnv(KeyBaiduAPIKey, BaiduAPIKeyEnv)
	_ = viper.BindEnv(KeyTencentNewsAPIKey, TencentNewsAPIKeyEnv)
	_ = viper.BindEnv(KeyWeiboAPIToken, WeiboAPITokenEnv)

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
	return viper.GetString(KeyBaiduAPIKey)
}

// GetTencentNewsAPIKey 获取腾讯新闻 API Key
func GetTencentNewsAPIKey() string {
	return viper.GetString(KeyTencentNewsAPIKey)
}

// GetWeiboAPIToken 获取微博 API Token
func GetWeiboAPIToken() string {
	return viper.GetString(KeyWeiboAPIToken)
}

// IsRestyDebugEnabled 是否开启 Resty Debug 模式
// 注意：此参数只能从配置文件读取，不能通过环境变量设置
func IsRestyDebugEnabled() bool {
	return viper.GetBool(KeyRestyDebug)
}

// GetSkillsPath 获取技能包路径
func GetSkillsPath() string {
	return viper.GetString(KeySkillsPath)
}

// SetBaiduAPIKey 设置百度 API Key
func SetBaiduAPIKey(apiKey string) {
	viper.Set(KeyBaiduAPIKey, apiKey)
}

// SetTencentNewsAPIKey 设置腾讯新闻 API Key
func SetTencentNewsAPIKey(apiKey string) {
	viper.Set(KeyTencentNewsAPIKey, apiKey)
}

// SetWeiboAPIToken 设置微博 API Token
func SetWeiboAPIToken(apiToken string) {
	viper.Set(KeyWeiboAPIToken, apiToken)
}

// SetSkillsPath 设置技能包路径
func SetSkillsPath(path string) {
	viper.Set(KeySkillsPath, path)
}

// baiduKeyManager 百度 API Key 管理器（懒加载）
var baiduKeyManager *APIKeyManager
var baiduKeyOnce sync.Once

// getBaiduKeyManager 获取百度 Key 管理器
func getBaiduKeyManager() *APIKeyManager {
	baiduKeyOnce.Do(func() {
		baiduKeyManager = NewAPIKeyManager(KeyBaiduAPIKey, BaiduAPIKeyEnv)
	})
	return baiduKeyManager
}

// tencentNewsKeyManager 腾讯新闻 API Key 管理器（懒加载）
var tencentNewsKeyManager *APIKeyManager
var tencentNewsKeyOnce sync.Once

// getTencentNewsKeyManager 获取腾讯新闻 Key 管理器
func getTencentNewsKeyManager() *APIKeyManager {
	tencentNewsKeyOnce.Do(func() {
		tencentNewsKeyManager = NewAPIKeyManager(KeyTencentNewsAPIKey, TencentNewsAPIKeyEnv)
	})
	return tencentNewsKeyManager
}

// weiboTokenManager 微博 API Token 管理器（懒加载）
var weiboTokenManager *APIKeyManager
var weiboTokenOnce sync.Once

// getWeiboTokenManager 获取微博 Token 管理器
func getWeiboTokenManager() *APIKeyManager {
	weiboTokenOnce.Do(func() {
		weiboTokenManager = NewAPIKeyManager(KeyWeiboAPIToken, WeiboAPITokenEnv)
	})
	return weiboTokenManager
}

// weiboExternalTokenManager 微博外部 Token 管理器（通过目录加载）
var weiboExternalTokenManager *TokenManager
var weiboExternalTokenOnce sync.Once

// getWeiboExternalTokenManager 获取微博外部 Token 管理器
func getWeiboExternalTokenManager() *TokenManager {
	weiboExternalTokenOnce.Do(func() {
		weiboExternalTokenManager = NewTokenManager(WeiboTokenDirEnv)
		// 尝试加载，如果失败则不标记为加载成功
		_ = weiboExternalTokenManager.Load()
	})
	return weiboExternalTokenManager
}

// GetBaiduAPIKeyRandom 随机获取一个百度 API Key
func GetBaiduAPIKeyRandom() string {
	return getBaiduKeyManager().GetRandom()
}

// GetBaiduAPIKeyNext 轮询获取下一个百度 API Key
func GetBaiduAPIKeyNext() string {
	return getBaiduKeyManager().GetNext()
}

// GetBaiduAPIKeys 获取所有百度 API Key
func GetBaiduAPIKeys() []string {
	return getBaiduKeyManager().GetAll()
}

// GetTencentNewsAPIKeyRandom 随机获取一个腾讯新闻 API Key
func GetTencentNewsAPIKeyRandom() string {
	return getTencentNewsKeyManager().GetRandom()
}

// GetTencentNewsAPIKeyNext 轮询获取下一个腾讯新闻 API Key
func GetTencentNewsAPIKeyNext() string {
	return getTencentNewsKeyManager().GetNext()
}

// GetTencentNewsAPIKeys 获取所有腾讯新闻 API Key
func GetTencentNewsAPIKeys() []string {
	return getTencentNewsKeyManager().GetAll()
}

// GetWeiboAPITokenRandom 随机获取一个微博 API Token
// 优先级：1. WEIBO_TOKEN_DIR 目录下的 JSON 配置文件 2. 配置文件/环境变量
func GetWeiboAPITokenRandom() string {
	// 首先尝试从外部目录加载
	tm := getWeiboExternalTokenManager()
	if tm.IsLoaded() {
		if token, err := tm.GetRandom(); err == nil && token != nil {
			return token.Token
		}
	}
	// 回退到从配置文件或环境变量获取
	return getWeiboTokenManager().GetRandom()
}

// GetWeiboAPITokenNext 轮询获取下一个微博 API Token
func GetWeiboAPITokenNext() string {
	return getWeiboTokenManager().GetNext()
}

// GetWeiboAPITokens 获取所有微博 API Token
func GetWeiboAPITokens() []string {
	return getWeiboTokenManager().GetAll()
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
	if err := os.MkdirAll(workDir, 0750); err != nil {
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
// 优先级：{PROJECT_NAME}_WORK_DIR > 配置文件所在目录 > 用户配置目录 > 当前目录
func GetWorkDir() string {
	// 检查环境变量
	workDir := os.Getenv(GetWorkDirEnv())
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
