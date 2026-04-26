package skills

import (
	"github.com/dtapp/skills/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Skill 技能接口，统一命令定义和 SkillInfo 生成
type Skill interface {
	// GetCommand 获取技能的 cobra.Command
	// 所有子命令和 flags 都在这里定义
	GetCommand() *cobra.Command

	// GetInfo 获取技能元数据（不包含 Commands，Commands 从 cobra.Command 自动提取）
	GetInfo() types.SkillInfo
}

// ExtractSkillInfo 从 Skill 接口提取完整的 SkillInfo（包含自动提取的 Commands）
func ExtractSkillInfo(skill Skill) types.SkillInfo {
	info := skill.GetInfo()
	cmd := skill.GetCommand()
	info.Commands = extractCommandsFromCobra(cmd)
	return info
}

// extractCommandsFromCobra 从 cobra.Command 提取 Commands 信息
func extractCommandsFromCobra(cmd *cobra.Command) []types.Command {
	var commands []types.Command

	for _, subCmd := range cmd.Commands() {
		command := types.Command{
			Name:        subCmd.Name(),
			Description: subCmd.Short,
			Usage:       types.ProjectName + " " + cmd.Name() + " " + subCmd.Use,
			Flags:       extractFlagsFromCobra(subCmd),
		}
		commands = append(commands, command)
	}

	return commands
}

// extractFlagsFromCobra 从 cobra.Command 提取 Flags 信息
func extractFlagsFromCobra(cmd *cobra.Command) []types.Flag {
	var flags []types.Flag

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flag := types.Flag{
			Name:        f.Name,
			Shorthand:   f.Shorthand,
			Description: f.Usage,
			Default:     f.DefValue,
		}
		// 检查是否是必填参数 (通过 Annotations 判断)
		if f.Annotations != nil {
			if _, ok := f.Annotations["cobra_annotation_bash_completion_one_required_flag"]; ok {
				flag.Required = true
			}
		}
		flags = append(flags, flag)
	})

	return flags
}
