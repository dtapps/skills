package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// WrapCommand 包装命令以记录使用统计
func WrapCommand(skillName string, cmd *cobra.Command) *cobra.Command {
	// 保存原始 Run 和 RunE
	originalRun := cmd.Run
	originalRunE := cmd.RunE

	// 包装 Run
	if originalRun != nil {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			_ = recordUsage(skillName, cmd.Name(), cmd, args, func() error {
				originalRun(cmd, args)
				return nil
			})
		}
	}

	// 包装 RunE
	if originalRunE != nil {
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return recordUsage(skillName, cmd.Name(), cmd, args, func() error {
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
func recordUsage(skillName, commandName string, cmd *cobra.Command, args []string, fn func() error) error {
	manager, err := GetManager()
	if err != nil {
		// 统计失败不影响主功能
		return fn()
	}

	startTime := time.Now()
	err = fn()
	duration := time.Since(startTime)

	// 构建参数字符串（包括 positional args 和 flags）
	argsStr := buildArgsString(cmd, args)

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

// buildArgsString 构建参数字符串，包括 positional args 和 flags
func buildArgsString(cmd *cobra.Command, args []string) string {
	var parts []string

	// 添加 positional args
	if len(args) > 0 {
		parts = append(parts, args...)
	}

	// 添加 flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			parts = append(parts, fmt.Sprintf("--%s=%s", flag.Name, flag.Value.String()))
		}
	})

	// 合并并截断
	result := strings.Join(parts, " ")
	if len(result) > 200 {
		result = result[:200] + "..."
	}

	return result
}

// WrapSkillCommands 包装技能命令组
// skillName 从 cmd.Name() 自动获取
func WrapSkillCommands(cmd *cobra.Command) *cobra.Command {
	// 为每个子命令添加统计
	for _, subCmd := range cmd.Commands() {
		WrapCommand(cmd.Name(), subCmd)
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
	manager := MustGetManager()
	if manager == nil {
		return
	}

	stat, exists := manager.GetSkillStats(skillName)
	if !exists {
		return
	}

	fmt.Printf("\n📊 %s 使用统计: 总调用 %d 次，成功 %d 次\n",
		skillName, stat.TotalCalls, stat.SuccessCalls)
}
