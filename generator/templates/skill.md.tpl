---
name: "{{.Name}}"
title: "{{.Title}}"
description: "{{.Description}}"
{{- if .License }}
license: "{{.License}}"
{{- end }}
{{- if .Compatibility }}
compatibility: "{{.Compatibility}}"
{{- end }}
{{- if .Metadata }}
metadata:
{{- range $key, $value := .Metadata }}
  {{$key}}: "{{$value}}"
{{- end }}
{{- end }}
{{- if .AllowedTools }}
allowed-tools: "{{.AllowedToolsStr}}"
{{- end }}
---

# {{.Title}}（{{.Name}}）

## 描述

{{.Description}}

## 使用方法

```bash
{{.CliName}} {{.Name}} [command] [flags]
```

## 命令

{{- range .Commands }}
### {{.Name}}

{{.Description}}

{{- if ne .LongDescription ""}}
{{if ne .Description .LongDescription }}
###### {{.LongDescription}}
{{end}}
{{- end}}

**用法：**
```bash
{{.Usage}}{{- range .Flags }}{{if .Required}} -{{.Shorthand}} <{{.Name}}>{{end}}{{end}}
```
{{- if .Flags }}
**参数：**

| 参数 | 简写 | 描述 | 必填 | 默认值 |
|------|------|------|------|--------|
{{- range .Flags }}
| `--{{.Name}}` | {{if .Shorthand}}`-{{.Shorthand}}`{{else}}-{{end}} | {{.Description}} | {{if .Required}}是{{else}}否{{end}} | {{if .Default}}{{.Default}}{{else}}-{{end}} |
{{- end }}
{{- end }}
{{- if .OutputFields }}

**输出字段：**

| 字段名 | 描述 |
|--------|------|
{{- range .OutputFields }}
| `{{.Name}}` | {{.Description}} |
{{- end }}
{{- end }}
{{ end }}
