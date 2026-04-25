package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dtapp/skills/types"
)

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
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("创建 skill 目录失败: %v", err)
	}

	// 生成 SKILL.md 文件
	if err := generateSkillMd(skillInfo, skillDir, cliName); err != nil {
		return fmt.Errorf("生成 SKILL.md 失败: %v", err)
	}

	// 如果是 picoclaw，生成 .skill-origin.json 文件
	if isPicoclaw {
		if err := generateSkillOrigin(skillDir); err != nil {
			return fmt.Errorf("生成 .skill-origin.json 失败: %v", err)
		}
	}

	return nil
}

// generateSkillMd 生成 SKILL.md 文件
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

	// 读取模板文件
	templatePath := filepath.Join("templates", "skill.md.tpl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("模板文件不存在或解析失败: %v", err)
	}

	filePath := filepath.Join(skillDir, "SKILL.md")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return tmpl.Execute(file, data)
}

// generateSkillOrigin 生成 .skill-origin.json 文件
func generateSkillOrigin(skillDir string) error {
	origin := types.SkillOrigin{
		Version:     1,
		OriginKind:  "manual",
		InstalledAt: time.Now().UnixMilli(),
	}

	data, err := json.MarshalIndent(origin, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(skillDir, ".skill-origin.json")
	return os.WriteFile(filePath, data, 0644)
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
