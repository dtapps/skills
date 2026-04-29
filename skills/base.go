package skills

import (
	"reflect"
	"sort"
	"strings"

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
			Name:            subCmd.Name(),
			Description:     subCmd.Short,
			LongDescription: subCmd.Long,
			Usage:           types.ProjectName + " " + cmd.Name() + " " + subCmd.Use,
			Flags:           extractFlagsFromCobra(subCmd),
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

// ExtractOutputFieldsFromStruct 从结构体反射获取输出字段信息
// 从 json tag 获取字段名，从 description tag 获取描述
// 返回的字段按名称排序，确保生成结果稳定
func ExtractOutputFieldsFromStruct(s any) []types.OutputField {
	var fields []types.OutputField
	t := reflect.TypeOf(s)
	for field := range t.Fields() {
		jsonTag := field.Tag.Get("json")
		// 提取 json tag 中的字段名（去掉 omitempty 等选项）
		fieldName := jsonTag
		if before, _, ok := strings.Cut(jsonTag, ","); ok {
			fieldName = before
		}
		if fieldName == "" {
			fieldName = field.Name
		}
		// 从 description tag 获取描述
		desc := field.Tag.Get("description")
		fields = append(fields, types.OutputField{
			Name:        fieldName,
			Description: desc,
		})
	}
	// 按字段名排序，确保生成结果稳定
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}
