package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/dtapp/skills/types"
)

//go:embed templates/skill.md.tpl
var templatesFS embed.FS

// SkillTemplateData 技能模板数据
type SkillTemplateData struct {
	types.SkillInfo
	CliName         string
	CommandsStr     string
	AllowedToolsStr string
	Date            string
}

// GenerateSkill 生成 skill 目录和文件
func GenerateSkill(skillInfo types.SkillInfo, outputDir string, isPicoclaw bool, cliName string) error {
	// 创建 skill 目录
	skillDir := filepath.Join(outputDir, skillInfo.Name)
	if err := os.MkdirAll(skillDir, 0750); err != nil {
		return fmt.Errorf("创建 skill 目录失败: %v", err)
	}

	// 生成 SKILL.md 文件（带变化检测）
	if err := generateSkillMd(skillInfo, skillDir, cliName); err != nil {
		return fmt.Errorf("生成 SKILL.md 失败: %v", err)
	}

	// 如果是 picoclaw，生成 .skill-origin.json 文件（如果文件不存在）
	if isPicoclaw {
		if err := generateSkillOrigin(skillDir); err != nil {
			return fmt.Errorf("生成 .skill-origin.json 失败: %v", err)
		}
	}

	return nil
}

// generateSkillMd 生成 SKILL.md 文件（带变化检测）
func generateSkillMd(skillInfo types.SkillInfo, skillDir string, cliName string) error {
	// 将 Version 和 Author 和 Date 加入到 Metadata 中
	if skillInfo.Metadata == nil {
		skillInfo.Metadata = make(map[string]string)
	}
	skillInfo.Metadata["version"] = skillInfo.Version
	skillInfo.Metadata["author"] = skillInfo.Author
	skillInfo.Metadata["date"] = time.Now().Format("2006-01-02 15:04:05")

	data := SkillTemplateData{
		SkillInfo:       skillInfo,
		CliName:         cliName,
		CommandsStr:     generateCommandsFromSkillInfo(skillInfo, cliName),
		AllowedToolsStr: strings.Join(skillInfo.AllowedTools, " "),
		Date:            time.Now().Format("2006-01-02"),
	}

	// 读取嵌入的模板文件
	templateContent, err := templatesFS.ReadFile("templates/skill.md.tpl")
	if err != nil {
		return fmt.Errorf("读取嵌入模板失败: %v", err)
	}
	tmpl, err := template.New("skill.md.tpl").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 生成新的内容到内存
	var newContent bytes.Buffer
	if err := tmpl.Execute(&newContent, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	filePath := filepath.Join(skillDir, "SKILL.md")

	// 检查现有文件是否存在
	// #nosec G304 - filePath 是程序内部构造的路径
	existingContent, err := os.ReadFile(filePath)
	if err == nil {
		// 文件存在，比较内容（忽略 date 元数据）
		if contentsEqualIgnoreDate(string(existingContent), newContent.String()) {
			// 内容相同，跳过写入
			return nil
		}
	}
	// 文件不存在或内容不同，写入新内容
	return os.WriteFile(filePath, newContent.Bytes(), 0600)
}

// contentsEqualIgnoreDate 比较两个 SKILL.md 内容是否相同（忽略 date 元数据）
func contentsEqualIgnoreDate(existing, new string) bool {
	// 使用正则表达式移除 date 元数据行
	datePattern := regexp.MustCompile(`(?m)^\s*date:\s*"[^"]*"\s*$`)

	existingNormalized := datePattern.ReplaceAllString(existing, "")
	newNormalized := datePattern.ReplaceAllString(new, "")

	return strings.TrimSpace(existingNormalized) == strings.TrimSpace(newNormalized)
}

// generateSkillOrigin 生成 .skill-origin.json 文件（仅在文件不存在时创建）
func generateSkillOrigin(skillDir string) error {
	filePath := filepath.Join(skillDir, ".skill-origin.json")

	// 如果文件已存在，跳过生成，保留原有的 installed_at
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	origin := types.SkillOrigin{
		Version:     1,
		OriginKind:  "manual",
		InstalledAt: time.Now().UnixMilli(),
	}

	data, err := json.MarshalIndent(origin, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}

// generateCommandsFromSkillInfo 从 SkillInfo 生成命令说明
func generateCommandsFromSkillInfo(skillInfo types.SkillInfo, cliName string) string {
	if len(skillInfo.Commands) == 0 {
		return fmt.Sprintf("- %s %s version - 显示版本信息\n- %s %s [command] - 执行具体命令", cliName, skillInfo.Name, cliName, skillInfo.Name)
	}

	var commands []string
	for _, cmd := range skillInfo.Commands {
		usage := strings.Replace(cmd.Usage, "dtskills", cliName, 1)
		commands = append(commands, fmt.Sprintf("- %s - %s", usage, cmd.Description))
	}
	return strings.Join(commands, "\n")
}
