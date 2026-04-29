package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestToken_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected bool
	}{
		{
			name: "未设置过期时间",
			token: Token{
				Token:     "test_token",
				ExpiresAt: time.Time{},
			},
			expected: false,
		},
		{
			name: "未来过期时间",
			token: Token{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "已过期",
			token: Token{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(-24 * time.Hour),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewTokenManager(t *testing.T) {
	tm := NewTokenManager("TEST_TOKEN_DIR")
	if tm == nil {
		t.Fatal("NewTokenManager() returned nil")
	}
	if tm.loaded {
		t.Error("新创建的 TokenManager 应该是未加载状态")
	}
	if tm.Count() != 0 {
		t.Error("新创建的 TokenManager 应该没有 Token")
	}
}

func TestTokenManager_Load(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	t.Run("空目录", func(t *testing.T) {
		tm := NewTokenManager("TEST_TOKEN_DIR")
		err := tm.Load()
		if err == nil {
			t.Error("空目录应该返回错误")
		}
	})

	t.Run("加载单个Token文件", func(t *testing.T) {
		// 创建单个 Token 文件
		singleTokenContent := `{"token": "token1", "name": "test_token"}`
		singleFile := filepath.Join(tempDir, "single.json")
		if err := os.WriteFile(singleFile, []byte(singleTokenContent), 0644); err != nil {
			t.Fatalf("写入文件失败: %v", err)
		}
		defer func() { _ = os.Remove(singleFile) }()

		tm := NewTokenManager("TEST_TOKEN_DIR")
		if err := tm.Load(); err != nil {
			t.Fatalf("Load() 失败: %v", err)
		}

		if !tm.IsLoaded() {
			t.Error("Load() 后应该是已加载状态")
		}
		if tm.Count() != 1 {
			t.Errorf("Token 数量错误: got %d, want 1", tm.Count())
		}
	})

	t.Run("加载Token数组文件", func(t *testing.T) {
		// 清理之前的文件
		_ = os.RemoveAll(tempDir)
		_ = os.MkdirAll(tempDir, 0755)

		// 创建 Token 数组文件
		arrayContent := `[
			{"token": "token1", "name": "token1"},
			{"token": "token2", "name": "token2"}
		]`
		arrayFile := filepath.Join(tempDir, "array.json")
		if err := os.WriteFile(arrayFile, []byte(arrayContent), 0644); err != nil {
			t.Fatalf("写入文件失败: %v", err)
		}

		tm := NewTokenManager("TEST_TOKEN_DIR")
		if err := tm.Load(); err != nil {
			t.Fatalf("Load() 失败: %v", err)
		}

		if tm.Count() != 2 {
			t.Errorf("Token 数量错误: got %d, want 2", tm.Count())
		}
	})

	t.Run("加载Map格式文件", func(t *testing.T) {
		// 清理之前的文件
		_ = os.RemoveAll(tempDir)
		_ = os.MkdirAll(tempDir, 0755)

		// 创建 Map 格式文件
		mapContent := `{"token1": "value1", "token2": "value2"}`
		mapFile := filepath.Join(tempDir, "map.json")
		if err := os.WriteFile(mapFile, []byte(mapContent), 0644); err != nil {
			t.Fatalf("写入文件失败: %v", err)
		}

		tm := NewTokenManager("TEST_TOKEN_DIR")
		if err := tm.Load(); err != nil {
			t.Fatalf("Load() 失败: %v", err)
		}

		if tm.Count() != 2 {
			t.Errorf("Token 数量错误: got %d, want 2", tm.Count())
		}
	})

	t.Run("环境变量未设置", func(t *testing.T) {
		_ = os.Unsetenv("TEST_TOKEN_DIR")
		tm := NewTokenManager("TEST_TOKEN_DIR")
		err := tm.Load()
		if err == nil {
			t.Error("环境变量未设置应该返回错误")
		}
	})
}

func TestTokenManager_GetRandom(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	// 创建 Token 文件
	content := `[
		{"token": "token1", "name": "token1", "weight": 1},
		{"token": "token2", "name": "token2", "weight": 2},
		{"token": "token3", "name": "token3", "weight": 3}
	]`
	tokenFile := filepath.Join(tempDir, "tokens.json")
	if err := os.WriteFile(tokenFile, []byte(content), 0644); err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	tm := NewTokenManager("TEST_TOKEN_DIR")
	if err := tm.Load(); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	t.Run("随机获取Token", func(t *testing.T) {
		// 多次获取，验证都能获取到
		results := make(map[string]bool)
		for i := 0; i < 100; i++ {
			token, err := tm.GetRandom()
			if err != nil {
				t.Fatalf("GetRandom() 失败: %v", err)
			}
			results[token.Token] = true
		}

		// 应该能获取到所有 Token
		if len(results) != 3 {
			t.Errorf("随机获取结果不完整: got %d unique tokens, want 3", len(results))
		}
	})

	t.Run("未加载时获取", func(t *testing.T) {
		newTm := NewTokenManager("TEST_TOKEN_DIR")
		_, err := newTm.GetRandom()
		if err == nil {
			t.Error("未加载时应该返回错误")
		}
	})
}

func TestTokenManager_GetRandom_WithExpired(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	// 创建包含过期 Token 的文件
	content := `[
		{"token": "valid_token", "name": "valid"},
		{"token": "expired_token", "name": "expired", "expires_at": "2020-01-01T00:00:00Z"}
	]`
	tokenFile := filepath.Join(tempDir, "tokens.json")
	if err := os.WriteFile(tokenFile, []byte(content), 0644); err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	tm := NewTokenManager("TEST_TOKEN_DIR")
	if err := tm.Load(); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	t.Run("只返回未过期Token", func(t *testing.T) {
		// 多次获取，验证只获取到未过期的
		for i := 0; i < 50; i++ {
			token, err := tm.GetRandom()
			if err != nil {
				t.Fatalf("GetRandom() 失败: %v", err)
			}
			if token.Token == "expired_token" {
				t.Error("不应该返回已过期的 Token")
			}
		}
	})

	t.Run("获取有效Token列表", func(t *testing.T) {
		validTokens := tm.GetValid()
		if len(validTokens) != 1 {
			t.Errorf("有效 Token 数量错误: got %d, want 1", len(validTokens))
		}
		if validTokens[0].Token != "valid_token" {
			t.Errorf("Token 值错误: got %s, want valid_token", validTokens[0].Token)
		}
	})

	t.Run("统计有效Token", func(t *testing.T) {
		if tm.Count() != 2 {
			t.Errorf("总数量错误: got %d, want 2", tm.Count())
		}
		if tm.CountValid() != 1 {
			t.Errorf("有效数量错误: got %d, want 1", tm.CountValid())
		}
	})
}

func TestTokenManager_AllExpired(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	// 创建全部过期的 Token 文件
	content := `[
		{"token": "token1", "name": "token1", "expires_at": "2020-01-01T00:00:00Z"},
		{"token": "token2", "name": "token2", "expires_at": "2020-06-01T00:00:00Z"}
	]`
	tokenFile := filepath.Join(tempDir, "tokens.json")
	if err := os.WriteFile(tokenFile, []byte(content), 0644); err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	tm := NewTokenManager("TEST_TOKEN_DIR")
	if err := tm.Load(); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	t.Run("全部过期时获取", func(t *testing.T) {
		_, err := tm.GetRandom()
		if err == nil {
			t.Error("全部过期时应该返回错误")
		}
	})

	t.Run("有效数量为零", func(t *testing.T) {
		if tm.CountValid() != 0 {
			t.Errorf("有效数量错误: got %d, want 0", tm.CountValid())
		}
	})
}

func TestTokenManager_RemoveExpired(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	// 创建包含过期 Token 的文件
	content := `[
		{"token": "valid_token", "name": "valid"},
		{"token": "expired_token1", "name": "expired1", "expires_at": "2020-01-01T00:00:00Z"},
		{"token": "expired_token2", "name": "expired2", "expires_at": "2020-06-01T00:00:00Z"}
	]`
	tokenFile := filepath.Join(tempDir, "tokens.json")
	if err := os.WriteFile(tokenFile, []byte(content), 0644); err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	tm := NewTokenManager("TEST_TOKEN_DIR")
	if err := tm.Load(); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	t.Run("清理过期Token", func(t *testing.T) {
		removed := tm.RemoveExpired()
		if removed != 2 {
			t.Errorf("清理数量错误: got %d, want 2", removed)
		}
		if tm.Count() != 1 {
			t.Errorf("剩余数量错误: got %d, want 1", tm.Count())
		}
	})
}

func TestTokenManager_AddToken(t *testing.T) {
	tm := NewTokenManager("TEST_TOKEN_DIR")

	token := Token{
		Token:  "dynamic_token",
		Name:   "dynamic",
		Weight: 5,
	}

	tm.AddToken(token)

	if tm.Count() != 1 {
		t.Errorf("添加后数量错误: got %d, want 1", tm.Count())
	}

	all := tm.GetAll()
	if len(all) != 1 || all[0].Token != "dynamic_token" {
		t.Error("添加的 Token 不正确")
	}
}

func TestTokenManager_GetSources(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "token_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 设置环境变量
	_ = os.Setenv("TEST_TOKEN_DIR", tempDir)
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	// 创建多个 Token 文件
	for i := 1; i <= 3; i++ {
		content := `{"token": "token", "name": "token"}`
		fileName := filepath.Join(tempDir, "token"+string(rune('0'+i))+".json")
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			t.Fatalf("写入文件失败: %v", err)
		}
	}

	tm := NewTokenManager("TEST_TOKEN_DIR")
	if err := tm.Load(); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	sources := tm.GetSources()
	if len(sources) != 3 {
		t.Errorf("来源数量错误: got %d, want 3", len(sources))
	}
}

func TestTokenManager_GetDirPath(t *testing.T) {
	_ = os.Setenv("TEST_TOKEN_DIR", "/path/to/tokens")
	defer func() { _ = os.Unsetenv("TEST_TOKEN_DIR") }()

	tm := NewTokenManager("TEST_TOKEN_DIR")
	path := tm.GetDirPath()
	if path != "/path/to/tokens" {
		t.Errorf("路径错误: got %s, want /path/to/tokens", path)
	}
}
