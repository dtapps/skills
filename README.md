# dtSkill CLI 工具

一个基于 Go 语言的 Skills 管理工具，使用 Cobra 和 charmbracelet/huh 库实现。

## 功能特性

- **生成技能包**：通过交互式界面生成 SKILL.md 文件
- **执行命令**：执行技能包中的命令（百度百科搜索、百度搜索）
- **配置管理**：管理 API Key 等配置
- **自升级**：支持自动检查并升级到新版本

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://cnb.cool/dtapp/skills.git
cd skills

# 构建
make build

# 或者构建所有平台
make build-all
```

### 下载预编译二进制

从 [Releases](https://cnb.cool/dtapp/skills/-/releases) 页面下载对应平台的二进制文件。

### 通过包管理器安装

#### Debian/Ubuntu (deb 包)

```bash
# 下载 deb 包
wget https://cnb.cool/dtapp/skills/-/releases/download/latest/dtskills_amd64.deb

# 安装
sudo dpkg -i dtskills_amd64.deb

# 如果缺少依赖，修复依赖关系
sudo apt-get install -f
```

安装后可直接使用 `dtskills` 命令：
```bash
dtskills --version
dtskills generate
```

#### RHEL/CentOS/Fedora (rpm 包)

```bash
# 下载 rpm 包
wget https://cnb.cool/dtapp/skills/-/releases/download/latest/dtskills.x86_64.rpm

# 安装
sudo rpm -i dtskills.x86_64.rpm

# 或者使用 dnf
sudo dnf install dtskills.x86_64.rpm
```

安装后可直接使用 `dtskills` 命令：
```bash
dtskills --version
dtskills generate
```

**注意**：通过 deb/rpm 安装后，二进制文件位于 `/usr/bin/dtskills`，配置文件位于 `/etc/dtskills/config.yaml`。

## 使用方法

### 查看帮助

```bash
./dtskills help
```

### 生成技能包

```bash
./dtskills generate
```

执行后会显示多选界面，选择要生成的技能包。支持 PicoClaw 格式的 `.skill-origin.json` 文件生成。

### 执行技能命令

#### 百度百科搜索

```bash
./dtskills baidu-baike search -k <关键词>
```

参数：
- `-k, --key`：检索关键字（必填）
- `-t, --type`：检索类型，可选 `lemmaTitle`(按词条名) 或 `lemmaId`(按词条ID)，默认 `lemmaTitle`

例如：
```bash
./dtskills baidu-baike search -k 刘德华
./dtskills baidu-baike search -k 人工智能 -t lemmaTitle
```

#### 百度搜索

```bash
./dtskills baidu-search search -q <关键词>
```

参数：
- `-q, --query`：搜索关键词（必填）
- `-c, --count`：返回结果数量，默认 10
- `-f, --freshness`：时间筛选，如 `pd`(最近24小时)、`pw`(最近7天)、`pm`(最近31天)、`py`(最近365天) 或 `YYYY-MM-DDtoYYYY-MM-DD`

例如：
```bash
# 基础搜索
./dtskills baidu-search search -q 人工智能

# 设置返回数量
./dtskills baidu-search search -q 旅游景点 -c 20

# 时间范围筛选
./dtskills baidu-search search -q 最新新闻 -f pd
./dtskills baidu-search search -q 2025年新闻 -f 2025-01-01to2025-12-31
```

### 配置管理

#### 设置配置项

```bash
./dtskills config set <key> <value>
```

例如：
```bash
# 设置百度 API Key
./dtskills config set api_keys.baidu your-api-key

# 设置技能包路径
./dtskills config set skills.path ./my-skills
```

#### 查看配置

```bash
./dtskills config list
```

#### 查看配置文件路径

```bash
./dtskills config path
```

### 版本信息

```bash
./dtskills version
```

### 自升级

```bash
./dtskills upgrade
```

## 项目结构

```
skills/
├── config/              # 配置管理
│   └── config.go
├── generator/           # 技能包生成器
│   └── generator.go
├── skills/              # 技能实现
│   ├── baidu_baike/     # 百度百科技能
│   │   ├── baidu_baike.go
│   │   └── api.go
│   └── baidu_search/    # 百度搜索技能
│       ├── baidu_search.go
│       └── api.go
├── templates/           # 模板文件
│   └── skill.md.tpl
├── types/               # 公共类型定义
│   └── skill.go
├── ui/                  # UI 组件
│   ├── multiselect.go
│   └── output.go
├── updater/             # 自升级功能
│   └── updater.go
├── main.go              # 主入口
├── go.mod               # Go 模块文件
├── go.sum               # 依赖版本锁定
├── Makefile             # 构建脚本
└── README.md            # 本说明文件
```

## 配置说明

配置文件支持以下位置（按优先级排序）：

1. 环境变量 `DTSKILLS_CONFIG` 指定的路径
2. 当前目录下的 `config.yaml`
3. 系统配置目录（如 `~/.config/dtskills/config.yaml`）

### 配置示例

```yaml
api_keys:
  baidu: your-baidu-api-key

skills:
  path: ./skills
```

## 开发

### 构建

```bash
# 本地构建
make build

# 构建所有平台
make build-all

# 清理构建产物
make clean
```

### 测试

```bash
# 运行单元测试
make test

# 测试百度百科搜索
make test-baike

# 测试百度搜索
make test-search

# 运行所有测试
make test-all
```

### 代码检查

```bash
make check
```

### 安装依赖

```bash
make deps
```

## 环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `DTSKILLS_CONFIG` | 配置文件路径 | `./config.yaml` |
| `BAIDU_API_KEY` | 百度 API Key | `your-api-key` |
| `BINARY_NAME` | 构建时的二进制名称 | `dtskills` |
| `VERSION` | 构建时的版本号 | `1.0.0` |

## 技术栈

- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [charmbracelet/huh](https://github.com/charmbracelet/huh) - 交互式表单
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - 样式渲染
- [minio/selfupdate](https://github.com/minio/selfupdate) - 自升级功能

## 许可证

MIT License
