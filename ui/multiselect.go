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
func RunMultiSelect(message string, options []Option) ([]string, error) {
	var selectedValues []string

	// 构建 huh 选项
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Display, opt.Value)
	}

	// 创建多选组件
	field := huh.NewMultiSelect[string]().
		Title(message).
		Options(huhOptions...).
		Value(&selectedValues).
		Description("按空格选择/取消，方向键移动，回车确认").
		Filterable(false).
		Height(len(options) + 2)

	// 运行表单
	form := huh.NewForm(huh.NewGroup(field))
	err := form.Run()
	if err != nil {
		return nil, err
	}

	return selectedValues, nil
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
