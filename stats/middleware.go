package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// WrapCommand 包装命令以记录使用统计
func WrapCommand(skillName string, cmd *cobra.Command) *cobra.Command {
	// 保存原始 Run 和 RunE
	originalRun := cmd.Run
	originalRunE := cmd.RunE

	// 包装 Run
	if originalRun != nil {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			_ = recordUsage(skillName, cmd.Name(), args, func() error {
				originalRun(cmd, args)
				return nil
			})
		}
	}

	// 包装 RunE
	if originalRunE != nil {
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return recordUsage(skillName, cmd.Name(), args, func() error {
				return originalRunE(cmd, args)
			})
		}
	}

	// 递归包装子命令
	for _, subCmd := range cmd.Commands() {
		WrapCommand(skillName, subCmd)
	}

	return cmd
}

// recordUsage 记录使用情况
func recordUsage(skillName, commandName string, args []string, fn func() error) error {
	manager, err := GetManager()
	if err != nil {
		// 统计失败不影响主功能
		return fn()
	}

	startTime := time.Now()
	err = fn()
	duration := time.Since(startTime)

	// 构建参数字符串
	argsStr := strings.Join(args, " ")
	if len(argsStr) > 200 {
		argsStr = argsStr[:200] + "..."
	}

	// 记录统计（忽略错误）
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		if len(errorMsg) > 100 {
			errorMsg = errorMsg[:100] + "..."
		}
	}

	_ = manager.RecordUsage(
		skillName,
		commandName,
		argsStr,
		duration,
		err == nil,
		errorMsg,
	)

	return err
}

// WrapSkillCommands 包装技能命令组
func WrapSkillCommands(skillName string, cmd *cobra.Command) *cobra.Command {
	// 为每个子命令添加统计
	for _, subCmd := range cmd.Commands() {
		WrapCommand(skillName, subCmd)
	}
	return cmd
}

// MustGetManager 获取管理器，失败时返回 nil
func MustGetManager() *Manager {
	m, _ := GetManager()
	return m
}

// PrintUsageHint 打印使用统计提示
func PrintUsageHint(skillName string) {
	manager, err := GetManager()
	if err != nil {
		return
	}

	stat, exists := manager.GetSkillStats(skillName)
	if !exists {
		fmt.Printf("\n💡 提示: 使用 '%s stats show %s' 查看使用统计\n", skillName, skillName)
		return
	}

	fmt.Printf("\n📊 %s 已使用 %d 次，平均耗时 %.0fms\n",
		skillName, stat.TotalCalls, stat.AvgDuration)
	fmt.Printf("💡 使用 '%s stats show %s' 查看详细统计\n", skillName, skillName)
}
