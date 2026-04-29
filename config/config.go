package config

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dtapp/skills/types"

	"github.com/spf13/viper"
)

// 环境变量名
const (
	// BaiduAPIKeyEnv 百度 API Key 环境变量
	BaiduAPIKeyEnv = "BAIDU_API_KEY"
	// TencentNewsAPIKeyEnv 腾讯新闻 API Key 环境变量
	TencentNewsAPIKeyEnv = "TENCENT_NEWS_API_KEY"
)

// getWorkDirEnv 获取工作目录环境变量名 (动态生成: {PROJECT_NAME}_WORK_DIR)
func getWorkDirEnv() string {
	return strings.ToUpper(types.ProjectName) + "_WORK_DIR"
}

// 配置键常量
const (
	// KeyBaiduAPIKey 百度 API Key 配置键
	KeyBaiduAPIKey = "api_keys.baidu"
	// KeyTencentNewsAPIKey 腾讯新闻 API Key 配置键
	KeyTencentNewsAPIKey = "api_keys.tencent_news"
	// KeyRestyDebug Resty Debug 配置键
	KeyRestyDebug = "debug.resty"
	// KeySkillsPath 技能包路径配置键
	KeySkillsPath = "skills.path"
)

// resolveWorkDir 解析工作目录（不依赖 viper 状态）
// 优先级：DTSKILLS_WORK_DIR > 用户配置目录 > 当前目录
func resolveWorkDir() string {
	// 检查环境变量
	workDir := os.Getenv(getWorkDirEnv())
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

// SetSkillsPath 设置技能包路径
func SetSkillsPath(path string) {
	viper.Set(KeySkillsPath, path)
}

// apiKeyManager 多 Key 管理器
type apiKeyManager struct {
	keys   []string
	index  int
	mu     sync.Mutex
	random *rand.Rand
}

// newAPIKeyManager 创建 API Key 管理器
// 支持四种配置格式：
// 1. 配置文件 - 字符串（单 Key）：api_keys.baidu: "key1"
// 2. 配置文件 - 数组（多 Key）：api_keys.baidu: ["key1", "key2", "key3"]
// 3. 环境变量（单 Key）：BAIDU_API_KEY="key1"
// 4. 环境变量（多 Key）：BAIDU_API_KEY="key1", BAIDU_API_KEY_1="key2", BAIDU_API_KEY_2="key3"
func newAPIKeyManager(configKey string, envVar string) *apiKeyManager {
	var keys []string

	// 优先检查环境变量（支持 _数字 后缀格式）
	keys = readEnvKeys(envVar)

	// 环境变量未设置，从配置文件读取
	if len(keys) == 0 {
		// 尝试获取数组
		if arr := viper.GetStringSlice(configKey); len(arr) > 0 {
			keys = arr
		} else if str := viper.GetString(configKey); str != "" {
			// 兼容旧格式：单字符串
			keys = []string{str}
		}
	}

	return &apiKeyManager{
		keys:   keys,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// readEnvKeys 从环境变量读取多个 Key
// 格式：ENV_VAR, ENV_VAR_1, ENV_VAR_2, ... (最多 20 个)
// 先读取基础变量，然后读取 _1 到 _19，中间有空不停止
func readEnvKeys(baseEnvVar string) []string {
	var keys []string
	maxKeys := 20

	// 先读取基础变量
	if baseValue := os.Getenv(baseEnvVar); baseValue != "" {
		keys = append(keys, baseValue)
	}

	// 并发读取 _1 到 _19
	var wg sync.WaitGroup
	keyChan := make(chan string, maxKeys-1)

	for i := 1; i < maxKeys; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			envKey := fmt.Sprintf("%s_%d", baseEnvVar, idx)
			if value := os.Getenv(envKey); value != "" {
				keyChan <- value
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(keyChan)
	}()

	// 收集结果
	for value := range keyChan {
		keys = append(keys, value)
	}

	return keys
}

// GetRandom 随机获取一个 Key
func (m *apiKeyManager) GetRandom() string {
	if len(m.keys) == 0 {
		return ""
	}
	return m.keys[m.random.Intn(len(m.keys))]
}

// GetNext 轮询获取下一个 Key
func (m *apiKeyManager) GetNext() string {
	if len(m.keys) == 0 {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	key := m.keys[m.index]
	m.index = (m.index + 1) % len(m.keys)
	return key
}

// GetAll 获取所有 Key
func (m *apiKeyManager) GetAll() []string {
	return m.keys
}

// baiduKeyManager 百度 API Key 管理器（懒加载）
var baiduKeyManager *apiKeyManager
var baiduKeyOnce sync.Once

// getBaiduKeyManager 获取百度 Key 管理器
func getBaiduKeyManager() *apiKeyManager {
	baiduKeyOnce.Do(func() {
		baiduKeyManager = newAPIKeyManager(KeyBaiduAPIKey, BaiduAPIKeyEnv)
	})
	return baiduKeyManager
}

// tencentNewsKeyManager 腾讯新闻 API Key 管理器（懒加载）
var tencentNewsKeyManager *apiKeyManager
var tencentNewsKeyOnce sync.Once

// getTencentNewsKeyManager 获取腾讯新闻 Key 管理器
func getTencentNewsKeyManager() *apiKeyManager {
	tencentNewsKeyOnce.Do(func() {
		tencentNewsKeyManager = newAPIKeyManager(KeyTencentNewsAPIKey, TencentNewsAPIKeyEnv)
	})
	return tencentNewsKeyManager
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
	workDir := os.Getenv(getWorkDirEnv())
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
