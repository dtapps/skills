package types

// ProjectName 项目名称
// 在构建时通过 -ldflags 注入，例如: -X github.com/dtapp/skills/types.ProjectName=dtskills
var ProjectName = "dtskills"

// SkillInfo 技能配置信息
type SkillInfo struct {
	Name          string            // 技能名称 (必填, 最多64字符, 小写字母/数字/连字符)
	Version       string            // 版本号
	Author        string            // 作者
	Description   string            // 描述 (必填, 最多1024字符)
	Homepage      string            // 主页
	License       string            // 许可证 (可选)
	Compatibility string            // 兼容性 (可选, 最多500字符)
	Metadata      map[string]string // 元数据 (可选), 可包含额外信息
	AllowedTools  []string          // 允许使用的工具 (可选)
	Commands      []Command         // 命令列表
}

// Command 命令信息
type Command struct {
	Name        string // 命令名称
	Description string // 命令描述
	Usage       string // 使用示例
	Flags       []Flag // 命令参数列表
}

// Flag 命令参数信息
type Flag struct {
	Name        string // 参数名称
	Shorthand   string // 简写（如 -q）
	Description string // 参数描述
	Required    bool   // 是否必填
	Default     string // 默认值
}

// SkillOrigin skill 来源信息（用于 picoclaw）
type SkillOrigin struct {
	Version     int    `json:"version"`
	OriginKind  string `json:"origin_kind"`
	InstalledAt int64  `json:"installed_at"`
}
