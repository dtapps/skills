package stats

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dtapp/skills/types"
)

const (
	statsFileName  = "stats.json"
	maxRecentUsage = 100 // 最多保留100条最近使用记录
)

// StatsData 统计数据根结构
type StatsData struct {
	Version     string                `json:"version"`
	LastUpdated time.Time             `json:"last_updated"`
	Skills      map[string]*SkillStat `json:"skills"`
	RecentUsage []UsageRecord         `json:"recent_usage"`
}

// SkillStat 单个技能统计
type SkillStat struct {
	TotalCalls   int64                   `json:"total_calls"`
	SuccessCalls int64                   `json:"success_calls"`
	FailedCalls  int64                   `json:"failed_calls"`
	AvgDuration  float64                 `json:"avg_duration_ms"`
	LastUsedAt   time.Time               `json:"last_used_at"`
	FirstUsedAt  time.Time               `json:"first_used_at"`
	Commands     map[string]*CommandStat `json:"commands"`
}

// CommandStat 命令统计
type CommandStat struct {
	Calls    int64     `json:"calls"`
	LastUsed time.Time `json:"last_used"`
}

// UsageRecord 单次使用记录
type UsageRecord struct {
	Skill     string        `json:"skill"`
	Command   string        `json:"command"`
	Args      string        `json:"args"`
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration_ms"`
	Success   bool          `json:"success"`
	ErrorMsg  string        `json:"error_msg,omitempty"`
}

// Manager 统计管理器
type Manager struct {
	data     *StatsData
	filePath string
	mu       sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取统计管理器单例
func GetManager() (*Manager, error) {
	var err error
	once.Do(func() {
		instance, err = newManager()
	})
	return instance, err
}

// newManager 创建新的统计管理器
func newManager() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("获取配置目录失败: %w", err)
	}

	projectConfigDir := filepath.Join(configDir, types.ProjectName)
	if err := os.MkdirAll(projectConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	filePath := filepath.Join(projectConfigDir, statsFileName)

	m := &Manager{
		data: &StatsData{
			Version:     "1.0",
			Skills:      make(map[string]*SkillStat),
			RecentUsage: make([]UsageRecord, 0),
		},
		filePath: filePath,
	}

	// 加载已有数据
	if err := m.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("加载统计数据失败: %w", err)
	}

	return m, nil
}

// load 从文件加载统计数据
func (m *Manager) load() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.data)
}

// save 保存统计数据到文件
func (m *Manager) save() error {
	m.data.LastUpdated = time.Now()

	data, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化统计数据失败: %w", err)
	}

	return os.WriteFile(m.filePath, data, 0644)
}

// RecordUsage 记录技能使用
func (m *Manager) RecordUsage(skill, command, args string, duration time.Duration, success bool, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// 获取或创建技能统计
	stat, exists := m.data.Skills[skill]
	if !exists {
		stat = &SkillStat{
			Commands: make(map[string]*CommandStat),
		}
		m.data.Skills[skill] = stat
		stat.FirstUsedAt = now
	}

	// 更新技能统计
	stat.TotalCalls++
	if success {
		stat.SuccessCalls++
	} else {
		stat.FailedCalls++
	}
	stat.LastUsedAt = now

	// 更新平均执行时间
	if stat.TotalCalls == 1 {
		stat.AvgDuration = float64(duration.Milliseconds())
	} else {
		stat.AvgDuration = (stat.AvgDuration*float64(stat.TotalCalls-1) + float64(duration.Milliseconds())) / float64(stat.TotalCalls)
	}

	// 更新命令统计
	cmdStat, exists := stat.Commands[command]
	if !exists {
		cmdStat = &CommandStat{}
		stat.Commands[command] = cmdStat
	}
	cmdStat.Calls++
	cmdStat.LastUsed = now

	// 添加使用记录
	record := UsageRecord{
		Skill:     skill,
		Command:   command,
		Args:      args,
		StartTime: now.Add(-duration),
		Duration:  duration,
		Success:   success,
		ErrorMsg:  errorMsg,
	}
	m.data.RecentUsage = append([]UsageRecord{record}, m.data.RecentUsage...)

	// 限制最近使用记录数量
	if len(m.data.RecentUsage) > maxRecentUsage {
		m.data.RecentUsage = m.data.RecentUsage[:maxRecentUsage]
	}

	return m.save()
}

// GetSkillStats 获取技能统计
func (m *Manager) GetSkillStats(skill string) (*SkillStat, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stat, exists := m.data.Skills[skill]
	return stat, exists
}

// GetAllStats 获取所有统计
func (m *Manager) GetAllStats() map[string]*SkillStat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 复制数据避免外部修改
	result := make(map[string]*SkillStat, len(m.data.Skills))
	maps.Copy(result, m.data.Skills)
	return result
}

// GetRecentUsage 获取最近使用记录
func (m *Manager) GetRecentUsage(limit int) []UsageRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.data.RecentUsage) {
		limit = len(m.data.RecentUsage)
	}

	result := make([]UsageRecord, limit)
	copy(result, m.data.RecentUsage[:limit])
	return result
}

// GetTotalCalls 获取总调用次数
func (m *Manager) GetTotalCalls() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int64
	for _, stat := range m.data.Skills {
		total += stat.TotalCalls
	}
	return total
}

// ResetStats 重置统计数据
func (m *Manager) ResetStats() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data.Skills = make(map[string]*SkillStat)
	m.data.RecentUsage = make([]UsageRecord, 0)
	m.data.LastUpdated = time.Now()

	return m.save()
}

// GetStatsFilePath 获取统计文件路径
func (m *Manager) GetStatsFilePath() string {
	return m.filePath
}
