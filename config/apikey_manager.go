package config

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sync"

	"github.com/spf13/viper"
)

// APIKeyManager 多 Key 管理器
type APIKeyManager struct {
	keys   []string
	index  int
	mu     sync.Mutex
	random *rand.Rand
}

// NewAPIKeyManager 创建 API Key 管理器
// 支持四种配置格式：
// 1. 配置文件 - 字符串（单 Key）：api_keys.baidu: "key1"
// 2. 配置文件 - 数组（多 Key）：api_keys.baidu: ["key1", "key2", "key3"]
// 3. 环境变量（单 Key）：BAIDU_API_KEY="key1"
// 4. 环境变量（多 Key）：BAIDU_API_KEY="key1", BAIDU_API_KEY_1="key2", BAIDU_API_KEY_2="key3"
func NewAPIKeyManager(configKey string, envVar string) *APIKeyManager {
	var keys []string

	// 优先检查环境变量（支持 _数字 后缀格式）
	keys = ReadEnvKeys(envVar)

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

	return &APIKeyManager{
		keys: keys,
		//nolint:gosec // 用于 API Key 随机选择，不需要加密安全级别
		random: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}
}

// ReadEnvKeys 从环境变量读取多个 Key
// 格式：ENV_VAR, ENV_VAR_1, ENV_VAR_2, ... (最多 20 个)
// 先读取基础变量，然后读取 _1 到 _19，中间有空不停止
func ReadEnvKeys(baseEnvVar string) []string {
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
func (m *APIKeyManager) GetRandom() string {
	if len(m.keys) == 0 {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// 支持直接用结构体字面量构造（测试与部分调用方），此时 random 未初始化。
	if m.random == nil {
		//nolint:gosec // 用于 API Key 随机选择，不需要加密安全级别
		m.random = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	return m.keys[m.random.IntN(len(m.keys))]
}

// GetNext 轮询获取下一个 Key
func (m *APIKeyManager) GetNext() string {
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
func (m *APIKeyManager) GetAll() []string {
	return m.keys
}
