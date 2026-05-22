package cmd

import (
	"fmt"
	"strings"

	"github.com/dtapp/skills/config"
	"github.com/dtapp/skills/ui"

	"github.com/spf13/cobra"
)

// NewConfigCmd 创建配置命令
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "管理配置",
		Long:  "查看和管理应用程序配置",
	}

	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigAddCmd())
	cmd.AddCommand(newConfigRemoveCmd())
	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigPathCmd())

	return cmd
}

// newConfigGetCmd 获取配置项
func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "获取配置项",
		Long:  "获取指定配置项的值，key 支持点号路径，如: api_keys.baidu",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())
			key := args[0]
			value := config.Get(key)
			p.Printf("%s: %v\n", key, value)
		},
	}
}

// newConfigSetCmd 设置配置项
func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "设置配置项",
		Long:  "设置指定配置项的值，key 支持点号路径，如: api_keys.baidu",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())
			key := args[0]
			value := args[1]

			if err := config.Set(key, value); err != nil {
				p.Printf("设置失败: %v\n", err)
				return
			}

			if err := config.SaveConfig(); err != nil {
				p.Printf("保存失败: %v\n", err)
				return
			}

			p.Printf("✓ %s = %s\n", key, value)
		},
	}
}

// newConfigAddCmd 添加配置项（追加到数组）
func newConfigAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [key] [value]",
		Short: "添加配置项（支持多 Key）",
		Long:  "向指定配置项添加值，如果该配置项已存在则追加到数组，如: config add api_keys.baidu new_key",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())
			key := args[0]
			value := args[1]

			// 获取现有值
			existing := config.Get(key)
			var keys []string

			switch v := existing.(type) {
			case []string:
				keys = v
			case string:
				if v != "" {
					keys = []string{v}
				}
			}

			// 追加新值
			keys = append(keys, value)

			// 设置回配置
			if err := config.Set(key, keys); err != nil {
				p.Printf("添加失败: %v\n", err)
				return
			}

			if err := config.SaveConfig(); err != nil {
				p.Printf("保存失败: %v\n", err)
				return
			}

			p.Printf("✓ %s 已添加，当前共 %d 个 Key\n", key, len(keys))
		},
	}
}

// formatConfigValue 格式化配置值，递归展开嵌套结构
func formatConfigValue(value interface{}, indent string) string {
	switch v := value.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}"
		}
		var result string
		for k, val := range v {
			formatted := formatConfigValue(val, indent+"  ")
			if strings.Contains(formatted, "\n") {
				result += fmt.Sprintf("\n%s%s:", indent, k) + formatted
			} else {
				result += fmt.Sprintf("\n%s%s: %s", indent, k, formatted)
			}
		}
		return result
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}
		var result string
		for i, val := range v {
			formatted := formatConfigValue(val, indent+"  ")
			result += fmt.Sprintf("\n%s[%d]: %s", indent, i, formatted)
		}
		return result
	case string:
		if v == "" {
			return "(空)"
		}
		// 如果字符串太长，截断显示
		if len(v) > 50 {
			return v[:50] + "..."
		}
		return v
	case nil:
		return "(null)"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// flattenConfig 将嵌套配置展平为可选项列表
func flattenConfig(cfg map[string]interface{}, prefix string) []struct {
	Key     string
	Display string
} {
	var result []struct {
		Key     string
		Display string
	}

	for key, value := range cfg {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// 递归展平嵌套 map
			subItems := flattenConfig(v, fullKey)
			result = append(result, subItems...)
		default:
			display := formatConfigValue(value, "")
			result = append(result, struct {
				Key     string
				Display string
			}{
				Key:     fullKey,
				Display: fmt.Sprintf("%s: %s", fullKey, display),
			})
		}
	}

	return result
}

// newConfigRemoveCmd 交互式移除配置项
func newConfigRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "交互式移除配置项",
		Long:  "以交互方式选择并移除配置项",
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())

			for {
				// 获取所有配置
				all := config.All()
				if len(all) == 0 {
					p.Println("暂无配置可删除")
					return
				}

				// 展平配置列表
				items := flattenConfig(all, "")
				if len(items) == 0 {
					p.Println("暂无配置可删除")
					return
				}

				// 构建选项列表
				var options []ui.Option
				for _, item := range items {
					options = append(options, ui.Option{
						Display: item.Display,
						Value:   item.Key,
					})
				}

				// 显示选择列表
				p.Title("📋 选择要删除的配置项")
				p.Println()

				// 使用单选让用户选择
				selected, err := ui.RunSelect("请选择要删除的配置项（按 q 退出）:", options)
				if err != nil {
					return
				}

				// 用户选择退出
				if selected == "" {
					return
				}

				// 获取当前值用于显示
				currentValue := config.Get(selected)
				displayValue := formatConfigValue(currentValue, "")

				// 确认是否删除
				confirmed := ui.RunConfirm(fmt.Sprintf("确定要删除配置项 '%s'%s 吗？", selected, displayValue), false)
				if !confirmed {
					p.Println("已取消删除")
					p.Println()
					continue
				}

				// 删除配置项
				_ = config.Set(selected, "")
				if err := config.SaveConfig(); err != nil {
					p.Printf("保存失败: %v\n", err)
					return
				}

				p.Printf("✓ '%s' 已删除\n", selected)
				p.Println()
			}
		},
	}
}

// newConfigListCmd 列出所有配置
func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有配置",
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())
			p.Title("📋 当前配置")
			p.Println()

			all := config.All()
			if len(all) == 0 {
				p.Println("暂无配置")
				return
			}

			for key, value := range all {
				p.Printf("%s: %v\n", key, value)
			}
		},
	}
}

// newConfigPathCmd 显示配置文件路径
func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "显示配置文件路径",
		Run: func(cmd *cobra.Command, args []string) {
			p := ui.NewPrinter(cmd.OutOrStdout())
			configFile := config.GetConfigFile()
			if configFile == "" {
				p.Println("未找到配置文件")
			} else {
				p.Printf("配置文件路径: %s\n", configFile)
			}
		},
	}
}
