package ui

import (
	"charm.land/huh/v2"
)

// Option 选项结构
type Option struct {
	Display string // 显示文本
	Value   string // 实际值
}

// RunMultiSelect 运行多选并返回选中的值
// 支持全选功能（第一个选项为"[全选]"）
func RunMultiSelect(message string, options []Option) ([]string, error) {
	var selectedValues []string

	// 添加全选选项
	allOption := Option{
		Display: "[✓ 全选所有技能]",
		Value:   "__ALL__",
	}
	allOptions := append([]Option{allOption}, options...)

	// 构建 huh 选项
	huhOptions := make([]huh.Option[string], len(allOptions))
	for i, opt := range allOptions {
		huhOptions[i] = huh.NewOption(opt.Display, opt.Value)
	}

	// 创建多选组件
	field := huh.NewMultiSelect[string]().
		Title(message).
		Options(huhOptions...).
		Value(&selectedValues).
		Description("按空格选择/取消，方向键移动，回车确认").
		Filterable(false).
		Height(len(allOptions) + 2)

	// 运行表单
	form := huh.NewForm(huh.NewGroup(field))
	err := form.Run()
	if err != nil {
		return nil, err
	}

	// 检查是否选择了全选
	for _, v := range selectedValues {
		if v == "__ALL__" {
			// 返回所有选项（排除全选标记）
			result := make([]string, len(options))
			for i, opt := range options {
				result[i] = opt.Value
			}
			return result, nil
		}
	}

	return selectedValues, nil
}

// RunSelect 运行单选并返回选中的值
// 如果用户取消（按 q 或 Ctrl+C），返回空字符串
func RunSelect(message string, options []Option) (string, error) {
	var selectedValue string

	// 构建 huh 选项
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Display, opt.Value)
	}

	// 创建单选组件
	field := huh.NewSelect[string]().
		Title(message).
		Options(huhOptions...).
		Value(&selectedValue).
		Description("方向键移动，回车确认，按 q 退出")

	// 运行表单
	form := huh.NewForm(huh.NewGroup(field))
	err := form.Run()
	if err != nil {
		// 用户取消（按 q 或 Ctrl+C）
		return "", err
	}

	return selectedValue, nil
}

// RunConfirm 运行确认并返回结果
func RunConfirm(message string, defaultValue bool) bool {
	var value bool

	// 创建确认组件
	field := huh.NewConfirm().
		Title(message).
		Value(&value).
		Affirmative("是").
		Negative("否")

	// 运行表单
	form := huh.NewForm(huh.NewGroup(field))
	err := form.Run()
	if err != nil {
		// 出错时（如用户中断 Ctrl+C）返回 false，取消操作
		return false
	}

	return value
}
