package config

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIKeyManager(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		configKey string
		envKey    string
		wantKeys  []string
	}{
		{
			name:      "环境变量单Key",
			envVars:   map[string]string{"TEST_API_KEY": "key1"},
			configKey: "test.api_key",
			envKey:    "TEST_API_KEY",
			wantKeys:  []string{"key1"},
		},
		{
			name: "环境变量多Key",
			envVars: map[string]string{
				"TEST_API_KEY":   "key1",
				"TEST_API_KEY_1": "key2",
				"TEST_API_KEY_2": "key3",
			},
			configKey: "test.api_key",
			envKey:    "TEST_API_KEY",
			wantKeys:  []string{"key1", "key2", "key3"},
		},
		{
			name: "环境变量中间有空",
			envVars: map[string]string{
				"TEST_API_KEY":   "key1",
				"TEST_API_KEY_2": "key3",
			},
			configKey: "test.api_key",
			envKey:    "TEST_API_KEY",
			wantKeys:  []string{"key1", "key3"},
		},
		{
			name:      "环境变量为空，使用配置文件",
			envVars:   map[string]string{},
			configKey: "test.api_key",
			envKey:    "TEST_API_KEY_EMPTY",
			wantKeys:  []string{"config_key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// 设置配置文件
			if len(tt.envVars) == 0 {
				viper.Set(tt.configKey, "config_key")
				defer viper.Reset()
			}

			manager := NewAPIKeyManager(tt.configKey, tt.envKey)
			gotKeys := manager.GetAll()

			// 由于并发读取，顺序可能不同，所以比较长度和内容
			assert.Equal(t, len(tt.wantKeys), len(gotKeys))
			for _, want := range tt.wantKeys {
				assert.Contains(t, gotKeys, want)
			}
		})
	}
}

func TestAPIKeyManagerGetRandom(t *testing.T) {
	manager := &APIKeyManager{
		keys: []string{"key1", "key2", "key3"},
	}

	// 测试随机获取
	for i := 0; i < 10; i++ {
		key := manager.GetRandom()
		assert.Contains(t, []string{"key1", "key2", "key3"}, key)
	}

	// 测试空管理器
	emptyManager := &APIKeyManager{keys: []string{}}
	assert.Equal(t, "", emptyManager.GetRandom())
}

func TestAPIKeyManagerGetNext(t *testing.T) {
	manager := &APIKeyManager{
		keys: []string{"key1", "key2", "key3"},
	}

	// 测试轮询
	assert.Equal(t, "key1", manager.GetNext())
	assert.Equal(t, "key2", manager.GetNext())
	assert.Equal(t, "key3", manager.GetNext())
	assert.Equal(t, "key1", manager.GetNext()) // 循环

	// 测试空管理器
	emptyManager := &APIKeyManager{keys: []string{}}
	assert.Equal(t, "", emptyManager.GetNext())
}

func TestAPIKeyManagerGetAll(t *testing.T) {
	manager := &APIKeyManager{
		keys: []string{"key1", "key2", "key3"},
	}

	keys := manager.GetAll()
	assert.Equal(t, []string{"key1", "key2", "key3"}, keys)

	// 测试空管理器
	emptyManager := &APIKeyManager{keys: []string{}}
	assert.Equal(t, []string{}, emptyManager.GetAll())
}

func TestReadEnvKeys(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		baseEnvVar string
		wantKeys   []string
	}{
		{
			name:       "基础变量",
			envVars:    map[string]string{"BASE_KEY": "key1"},
			baseEnvVar: "BASE_KEY",
			wantKeys:   []string{"key1"},
		},
		{
			name: "基础变量加后缀",
			envVars: map[string]string{
				"BASE_KEY":   "key1",
				"BASE_KEY_1": "key2",
				"BASE_KEY_2": "key3",
			},
			baseEnvVar: "BASE_KEY",
			wantKeys:   []string{"key1", "key2", "key3"},
		},
		{
			name:       "无变量",
			envVars:    map[string]string{},
			baseEnvVar: "NOT_EXIST",
			wantKeys:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			gotKeys := ReadEnvKeys(tt.baseEnvVar)

			if tt.wantKeys == nil {
				assert.Nil(t, gotKeys)
			} else {
				assert.Equal(t, len(tt.wantKeys), len(gotKeys))
				for _, want := range tt.wantKeys {
					assert.Contains(t, gotKeys, want)
				}
			}
		})
	}
}

func TestGetBaiduAPIKeyFunctions(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("BAIDU_API_KEY", "baidu_key1")
	os.Setenv("BAIDU_API_KEY_1", "baidu_key2")
	defer os.Unsetenv("BAIDU_API_KEY")
	defer os.Unsetenv("BAIDU_API_KEY_1")

	// 重置单例，确保使用新的环境变量
	baiduKeyManager = nil
	baiduKeyOnce = sync.Once{}

	// 测试获取所有 Key
	keys := GetBaiduAPIKeys()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "baidu_key1")
	assert.Contains(t, keys, "baidu_key2")

	// 测试随机获取
	randomKey := GetBaiduAPIKeyRandom()
	assert.Contains(t, []string{"baidu_key1", "baidu_key2"}, randomKey)

	// 测试轮询获取
	nextKey1 := GetBaiduAPIKeyNext()
	nextKey2 := GetBaiduAPIKeyNext()
	assert.Contains(t, []string{"baidu_key1", "baidu_key2"}, nextKey1)
	assert.Contains(t, []string{"baidu_key1", "baidu_key2"}, nextKey2)
}

func TestGetTencentNewsAPIKeyFunctions(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TENCENT_NEWS_API_KEY", "tencent_key1")
	defer os.Unsetenv("TENCENT_NEWS_API_KEY")

	// 重置单例
	tencentNewsKeyManager = nil
	tencentNewsKeyOnce = sync.Once{}

	// 测试获取所有 Key
	keys := GetTencentNewsAPIKeys()
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "tencent_key1", keys[0])
}
