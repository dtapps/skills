package config

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Token 定义标准的 Token 结构
type Token struct {
	Token     string    `json:"token"`            // Token 值
	Name      string    `json:"name"`             // Token 名称/标识
	ExpiresAt time.Time `json:"expires_at"`       // 过期时间（可选）
	Weight    int       `json:"weight,omitempty"` // 权重，用于随机选择（默认 1）
}

// IsExpired 检查 Token 是否已过期
// 如果 ExpiresAt 为零值，表示永不过期
func (t *Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.ExpiresAt)
}

// TokenManager Token 管理器
// 用于管理多个 Token，支持随机选择和过期检查
type TokenManager struct {
	tokens  []Token
	envVar  string
	loaded  bool
	mu      sync.RWMutex
	random  *rand.Rand
	sources []string // 记录加载的配置文件路径
}

// NewTokenManager 创建 Token 管理器
// envVar: 环境变量名，该环境变量的值应该是包含 JSON 配置文件的目录路径
// 例如：WEIBO_TOKEN_DIR=/path/to/tokens/
func NewTokenManager(envVar string) *TokenManager {
	return &TokenManager{
		tokens: make([]Token, 0),
		envVar: envVar,
		//nolint:gosec // 用于 Token 随机选择，不需要加密安全级别
		random:  rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		sources: make([]string, 0),
	}
}

// Load 加载 Token 配置文件
// 读取指定目录下的所有 .json 文件
func (tm *TokenManager) Load() error {
	dirPath := os.Getenv(tm.envVar)
	if dirPath == "" {
		return fmt.Errorf("环境变量 %s 未设置", tm.envVar)
	}
	return tm.LoadFromDir(dirPath)
}

// LoadFromDir 从指定目录加载 Token 配置文件
// 读取目录下的所有 .json 文件
func (tm *TokenManager) LoadFromDir(dirPath string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.loaded {
		return nil
	}

	// 检查目录是否存在
	dirPath = filepath.Clean(dirPath)
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
		if err := tm.loadTokenFile(configPath); err != nil {
			return fmt.Errorf("读取配置文件 %s 失败: %w", entry.Name(), err)
		}

		tm.sources = append(tm.sources, configPath)
	}

	if len(tm.tokens) == 0 {
		return fmt.Errorf("配置目录 %s 下未找到有效的 Token", dirPath)
	}

	tm.loaded = true
	return nil
}

// loadTokenFile 加载单个 Token 文件
func (tm *TokenManager) loadTokenFile(path string) error {
	path = filepath.Clean(path)
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return err
	}

	// 尝试解析为 Token 数组
	var tokenList []Token
	if err := json.Unmarshal(data, &tokenList); err == nil {
		tm.tokens = append(tm.tokens, tokenList...)
		return nil
	}

	// 尝试解析为单个 Token
	var singleToken Token
	if err := json.Unmarshal(data, &singleToken); err == nil {
		tm.tokens = append(tm.tokens, singleToken)
		return nil
	}

	// 尝试解析为 map 格式（兼容旧格式）
	var tokenMap map[string]string
	if err := json.Unmarshal(data, &tokenMap); err == nil {
		for name, tokenValue := range tokenMap {
			tm.tokens = append(tm.tokens, Token{
				Token:  tokenValue,
				Name:   name,
				Weight: 1,
			})
		}
		return nil
	}

	return fmt.Errorf("无法解析 Token 文件格式")
}

// GetRandom 随机获取一个未过期的 Token
// 支持按权重随机选择
func (tm *TokenManager) GetRandom() (*Token, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if !tm.loaded {
		return nil, fmt.Errorf("Token 未加载")
	}

	// 过滤出未过期的 Token
	validTokens := make([]Token, 0)
	totalWeight := 0
	for _, token := range tm.tokens {
		if !token.IsExpired() {
			validTokens = append(validTokens, token)
			weight := token.Weight
			if weight <= 0 {
				weight = 1
			}
			totalWeight += weight
		}
	}

	if len(validTokens) == 0 {
		return nil, fmt.Errorf("没有可用的 Token（可能都已过期）")
	}

	// 按权重随机选择
	r := tm.random.IntN(totalWeight)
	for _, token := range validTokens {
		weight := token.Weight
		if weight <= 0 {
			weight = 1
		}
		r -= weight
		if r < 0 {
			return &token, nil
		}
	}

	// 兜底返回最后一个
	return &validTokens[len(validTokens)-1], nil
}

// GetAll 获取所有 Token（包括已过期的）
func (tm *TokenManager) GetAll() []Token {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tokens := make([]Token, len(tm.tokens))
	copy(tokens, tm.tokens)
	return tokens
}

// GetValid 获取所有未过期的 Token
func (tm *TokenManager) GetValid() []Token {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	validTokens := make([]Token, 0)
	for _, token := range tm.tokens {
		if !token.IsExpired() {
			validTokens = append(validTokens, token)
		}
	}
	return validTokens
}

// Count 获取 Token 总数
func (tm *TokenManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.tokens)
}

// CountValid 获取未过期的 Token 数量
func (tm *TokenManager) CountValid() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	count := 0
	for _, token := range tm.tokens {
		if !token.IsExpired() {
			count++
		}
	}
	return count
}

// IsLoaded 检查是否已加载
func (tm *TokenManager) IsLoaded() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.loaded
}

// GetDirPath 获取配置目录路径
func (tm *TokenManager) GetDirPath() string {
	return os.Getenv(tm.envVar)
}

// GetSources 获取已加载的配置文件路径列表
func (tm *TokenManager) GetSources() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	sources := make([]string, len(tm.sources))
	copy(sources, tm.sources)
	return sources
}

// AddToken 动态添加 Token（用于运行时更新）
func (tm *TokenManager) AddToken(token Token) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tokens = append(tm.tokens, token)
}

// RemoveExpired 移除所有已过期的 Token
func (tm *TokenManager) RemoveExpired() int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	originalCount := len(tm.tokens)
	validTokens := make([]Token, 0)
	for _, token := range tm.tokens {
		if !token.IsExpired() {
			validTokens = append(validTokens, token)
		}
	}
	tm.tokens = validTokens
	return originalCount - len(validTokens)
}
