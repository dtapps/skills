# 自动补全

`dtskills` 支持为 Bash、Zsh 和 PowerShell 生成自动补全脚本，让你在输入命令时按 `Tab` 键自动补全命令和参数。

## 功能说明

安装自动补全后，你可以：

- 输入 `dtskills <Tab>` 显示所有可用命令
- 输入 `dtskills bai<Tab>` 自动补全为 `baidu-baike` 或 `baidu-search`
- 输入 `dtskills config <Tab>` 显示 config 的所有子命令
- 输入 `dtskills --<Tab>` 显示所有可用 flags

---

## Linux (Debian/Ubuntu)

Debian 和 Ubuntu 默认使用 Bash。

### 首次安装

```bash
# 生成补全脚本到系统目录
sudo dtskills completion bash > /etc/bash_completion.d/dtskills

# 重新加载配置
source /etc/bash_completion.d/dtskills
```

安装后重新打开终端，输入 `dtskills <Tab>` 即可看到自动补全效果。

### 更新补全脚本

当你升级 dtskills 或添加新命令后，必须重新生成补全脚本：

```bash
sudo dtskills completion bash > /etc/bash_completion.d/dtskills
```

更新后重新打开终端即可生效。

### 故障排除

**补全不生效？**

1. 检查补全脚本是否正确生成：
   ```bash
   cat /etc/bash_completion.d/dtskills
   ```

2. 确保当前 shell 是 Bash：
   ```bash
   echo $SHELL  # 应输出 /bin/bash
   ```

3. 重新加载配置：
   ```bash
   source ~/.bashrc
   ```

---

## macOS

macOS 新版默认使用 Zsh。

### 首次安装

**使用 oh-my-zsh**：
```bash
mkdir -p ~/.oh-my-zsh/completions
dtskills completion zsh > ~/.oh-my-zsh/completions/_dtskills
```

**使用 Homebrew**：
```bash
mkdir -p $(brew --prefix)/share/zsh/site-functions
dtskills completion zsh > $(brew --prefix)/share/zsh/site-functions/_dtskills
```

安装后重新打开终端即可生效。

### 更新补全脚本

当你升级 dtskills 或添加新命令后，必须重新生成补全脚本：

```bash
dtskills completion zsh > ~/.oh-my-zsh/completions/_dtskills
```

或（如果使用 Homebrew）：
```bash
dtskills completion zsh > $(brew --prefix)/share/zsh/site-functions/_dtskills
```

更新后重新打开终端即可生效。

### 故障排除

**补全不生效？**

1. 检查补全脚本是否正确生成：
   ```bash
   cat ~/.oh-my-zsh/completions/_dtskills
   ```

2. 确保当前 shell 是 Zsh：
   ```bash
   echo $SHELL  # 应输出 /bin/zsh
   ```

3. 重新加载配置：
   ```bash
   source ~/.zshrc
   ```

---

## Windows

Windows 使用 PowerShell。

### 首次安装

```powershell
# 创建配置文件（如果不存在）
New-Item -Path $PROFILE -Type File -Force

# 添加补全脚本
dtskills completion powershell >> $PROFILE
```

安装后重新打开 PowerShell，输入 `dtskills <Tab>` 即可看到自动补全效果。

### 更新补全脚本

当你升级 dtskills 或添加新命令后，必须重新生成补全脚本：

```powershell
dtskills completion powershell > $PROFILE
```

更新后重新打开 PowerShell 即可生效。

### 故障排除

**补全不生效？**

1. 检查配置文件路径：
   ```powershell
   echo $PROFILE
   ```

2. 检查补全脚本是否正确生成：
   ```powershell
   Get-Content $PROFILE
   ```

3. 确保执行策略允许运行脚本：
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

---

## 验证安装

安装完成后，重新打开终端，输入以下命令验证：

```bash
# 查看所有命令
dtskills <Tab><Tab>

# 补全命令
dtskills conf<Tab>  # 会自动补全为 config

# 查看子命令
dtskills config <Tab><Tab>
```

## 其他 Shell

如果你使用其他 shell（如 Fish），可以查看支持列表：

```bash
dtskills completion --help
```

然后使用对应的命令生成补全脚本。
