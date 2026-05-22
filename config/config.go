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
	// HefengAPIKeyEnv 和风天气 API Key 环境变量
	HefengAPIKeyEnv = "HEFENG_API_KEY"
	// WeiboAPITokenEnv 微博API Token 环境变量
	WeiboAPITokenEnv = "WEIBO_API_TOKEN"
	// WeiboTokenDirEnv 微博 Token 目录环境变量
	WeiboTokenDirEnv = "WEIBO_TOKEN_DIR"
	// BaiduMapAPIKeyEnv 百度地图 API Key 环境变量
	BaiduMapAPIKeyEnv = "BAIDU_MAP_API_KEY"
	// AmapMapAPIKeyEnv 高德地图 API Key 环境变量
	AmapMapAPIKeyEnv = "AMAP_MAP_API_KEY"
	// TencentMapAPIKeyEnv 腾讯地图 API Key 环境变量
	TencentMapAPIKeyEnv = "TENCENT_MAP_API_KEY"
	// UAPIKeyEnv UAPI Key 环境变量
	UAPIKeyEnv = "UAPI_KEY"
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
	// KeyHefengAPIKey 和风天气 API Key 配置键
	KeyHefengAPIKey = "api_keys.hefeng"
	// KeyWeiboAPIToken 微博 API Token 配置键
	KeyWeiboAPIToken = "api_tokens.weibo"
	// KeyWeiboTokenDir 微博 Token 目录配置键
	KeyWeiboTokenDir = "api_tokens.weibo.token_dir"
	// KeyBaiduMapAPIKey 百度地图 API Key 配置键
	KeyBaiduMapAPIKey = "api_keys.baidu_map"
	// KeyAmapMapAPIKey 高德地图 API Key 配置键
	KeyAmapMapAPIKey = "api_keys.amap_map"
	// KeyTencentMapAPIKey 腾讯地图 API Key 配置键
	KeyTencentMapAPIKey = "api_keys.tencent_map"
	// KeyUAPIKey UAPI Key 配置键
	KeyUAPIKey = "api_keys.uapi"
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
	_ = viper.BindEnv(KeyHefengAPIKey, HefengAPIKeyEnv)
	_ = viper.BindEnv(KeyWeiboAPIToken, WeiboAPITokenEnv)
	_ = viper.BindEnv(KeyWeiboTokenDir, WeiboTokenDirEnv)
	_ = viper.BindEnv(KeyBaiduMapAPIKey, BaiduMapAPIKeyEnv)
	_ = viper.BindEnv(KeyAmapMapAPIKey, AmapMapAPIKeyEnv)
	_ = viper.BindEnv(KeyTencentMapAPIKey, TencentMapAPIKeyEnv)
	_ = viper.BindEnv(KeyUAPIKey, UAPIKeyEnv)

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

// GetHefengAPIKey 获取和风天气 API Key
func GetHefengAPIKey() string {
	return viper.GetString(KeyHefengAPIKey)
}

// GetWeiboAPIToken 获取微博 API Token
func GetWeiboAPIToken() string {
	return viper.GetString(KeyWeiboAPIToken)
}

// GetWeiboTokenDir 获取微博 Token 目录
func GetWeiboTokenDir() string {
	return viper.GetString(KeyWeiboTokenDir)
}

// GetBaiduMapAPIKey 获取百度地图 API Key
func GetBaiduMapAPIKey() string {
	return viper.GetString(KeyBaiduMapAPIKey)
}

// GetAmapMapAPIKey 获取高德地图 API Key
func GetAmapMapAPIKey() string {
	return viper.GetString(KeyAmapMapAPIKey)
}

// GetTencentMapAPIKey 获取腾讯地图 API Key
func GetTencentMapAPIKey() string {
	return viper.GetString(KeyTencentMapAPIKey)
}

// GetUAPIKey 获取 UAPI Key
func GetUAPIKey() string {
	return viper.GetString(KeyUAPIKey)
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

// SetHefengAPIKey 设置和风天气 API Key
func SetHefengAPIKey(apiKey string) {
	viper.Set(KeyHefengAPIKey, apiKey)
}

// SetWeiboAPIToken 设置微博 API Token
func SetWeiboAPIToken(apiToken string) {
	viper.Set(KeyWeiboAPIToken, apiToken)
}

// SetUAPIKey 设置 UAPI Key
func SetUAPIKey(apiKey string) {
	viper.Set(KeyUAPIKey, apiKey)
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

// hefengKeyManager 和风天气 API Key 管理器（懒加载）
var hefengKeyManager *APIKeyManager
var hefengKeyOnce sync.Once

// getHefengKeyManager 获取和风天气 Key 管理器
func getHefengKeyManager() *APIKeyManager {
	hefengKeyOnce.Do(func() {
		hefengKeyManager = NewAPIKeyManager(KeyHefengAPIKey, HefengAPIKeyEnv)
	})
	return hefengKeyManager
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
// 如果没有配置 token_dir，返回 nil
func getWeiboExternalTokenManager() *TokenManager {
	weiboExternalTokenOnce.Do(func() {
		// 优先从配置读取 token_dir，如果不存在则使用环境变量
		tokenDir := viper.GetString(KeyWeiboTokenDir)
		if tokenDir == "" {
			tokenDir = os.Getenv(WeiboTokenDirEnv)
		}
		if tokenDir == "" {
			return
		}
		weiboExternalTokenManager = NewTokenManager("")
		// 直接从配置的路径加载
		_ = weiboExternalTokenManager.LoadFromDir(tokenDir)
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

// GetHefengAPIKeyRandom 随机获取一个和风天气 API Key
func GetHefengAPIKeyRandom() string {
	return getHefengKeyManager().GetRandom()
}

// GetHefengAPIKeyNext 轮询获取下一个和风天气 API Key
func GetHefengAPIKeyNext() string {
	return getHefengKeyManager().GetNext()
}

// GetHefengAPIKeys 获取所有和风天气 API Key
func GetHefengAPIKeys() []string {
	return getHefengKeyManager().GetAll()
}

// baiduMapKeyManager 百度地图 API Key 管理器（懒加载）
var baiduMapKeyManager *APIKeyManager
var baiduMapKeyOnce sync.Once

// getBaiduMapKeyManager 获取百度地图 Key 管理器
func getBaiduMapKeyManager() *APIKeyManager {
	baiduMapKeyOnce.Do(func() {
		baiduMapKeyManager = NewAPIKeyManager(KeyBaiduMapAPIKey, BaiduMapAPIKeyEnv)
	})
	return baiduMapKeyManager
}

// amapMapKeyManager 高德地图 API Key 管理器（懒加载）
var amapMapKeyManager *APIKeyManager
var amapMapKeyOnce sync.Once

// getAmapMapKeyManager 获取高德地图 Key 管理器
func getAmapMapKeyManager() *APIKeyManager {
	amapMapKeyOnce.Do(func() {
		amapMapKeyManager = NewAPIKeyManager(KeyAmapMapAPIKey, AmapMapAPIKeyEnv)
	})
	return amapMapKeyManager
}

// uapiKeyManager UAPI Key 管理器（懒加载）
var uapiKeyManager *APIKeyManager
var uapiKeyOnce sync.Once

// getUapiKeyManager 获取 UAPI Key 管理器
func getUapiKeyManager() *APIKeyManager {
	uapiKeyOnce.Do(func() {
		uapiKeyManager = NewAPIKeyManager(KeyUAPIKey, UAPIKeyEnv)
	})
	return uapiKeyManager
}

// GetUAPIKeyRandom 随机获取一个 UAPI Key
func GetUAPIKeyRandom() string {
	return getUapiKeyManager().GetRandom()
}

// GetUAPIKeyNext 轮询获取下一个 UAPI Key
func GetUAPIKeyNext() string {
	return getUapiKeyManager().GetNext()
}

// GetUAPIKeys 获取所有 UAPI Key
func GetUAPIKeys() []string {
	return getUapiKeyManager().GetAll()
}

// tencentMapKeyManager 腾讯地图 API Key 管理器（懒加载）
var tencentMapKeyManager *APIKeyManager
var tencentMapKeyOnce sync.Once

// getTencentMapKeyManager 获取腾讯地图 Key 管理器
func getTencentMapKeyManager() *APIKeyManager {
	tencentMapKeyOnce.Do(func() {
		tencentMapKeyManager = NewAPIKeyManager(KeyTencentMapAPIKey, TencentMapAPIKeyEnv)
	})
	return tencentMapKeyManager
}

// GetBaiduMapAPIKeyRandom 随机获取一个百度地图 API Key
func GetBaiduMapAPIKeyRandom() string {
	return getBaiduMapKeyManager().GetRandom()
}

// GetBaiduMapAPIKeyNext 轮询获取下一个百度地图 API Key
func GetBaiduMapAPIKeyNext() string {
	return getBaiduMapKeyManager().GetNext()
}

// GetBaiduMapAPIKeys 获取所有百度地图 API Key
func GetBaiduMapAPIKeys() []string {
	return getBaiduMapKeyManager().GetAll()
}

// GetAmapMapAPIKeyRandom 随机获取一个高德地图 API Key
func GetAmapMapAPIKeyRandom() string {
	return getAmapMapKeyManager().GetRandom()
}

// GetAmapMapAPIKeyNext 轮询获取下一个高德地图 API Key
func GetAmapMapAPIKeyNext() string {
	return getAmapMapKeyManager().GetNext()
}

// GetAmapMapAPIKeys 获取所有高德地图 API Key
func GetAmapMapAPIKeys() []string {
	return getAmapMapKeyManager().GetAll()
}

// GetTencentMapAPIKeyRandom 随机获取一个腾讯地图 API Key
func GetTencentMapAPIKeyRandom() string {
	return getTencentMapKeyManager().GetRandom()
}

// GetTencentMapAPIKeyNext 轮询获取下一个腾讯地图 API Key
func GetTencentMapAPIKeyNext() string {
	return getTencentMapKeyManager().GetNext()
}

// GetTencentMapAPIKeys 获取所有腾讯地图 API Key
func GetTencentMapAPIKeys() []string {
	return getTencentMapKeyManager().GetAll()
}

// GetWeiboAPITokenRandom 随机获取一个微博 API Token
// 优先级：1. WEIBO_TOKEN_DIR 目录下的 JSON 配置文件 2. 配置文件/环境变量
func GetWeiboAPITokenRandom() string {
	// 首先尝试从外部目录加载
	tm := getWeiboExternalTokenManager()
	if tm != nil && tm.IsLoaded() {
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
