# dtSkill CLI 工具

一个基于 Go 语言的 Skills 管理工具，使用 Cobra 和 charmbracelet/huh 库实现。

<div align="center">

[![CNB 版本](https://cnb.cool/dtapp/skills/-/badge/release)](https://cnb.cool/dtapp/skills/-/badge/release.link) [![CNB Star 数](https://cnb.cool/dtapp/skills/-/badge/star)](https://cnb.cool/dtapp/skills) [![CNB Fork 数](https://cnb.cool/dtapp/skills/-/badge/fork)](https://cnb.cool/dtapp/skills) [![CNB 构建](https://cnb.cool/dtapp/skills/-/badge/git/latest/ci/status/push)](https://cnb.cool/dtapp/skills) [![CNB 覆盖率](https://cnb.cool/dtapp/skills/-/badge/git/latest/testing/unit/coverage)](https://cnb.cool/dtapp/skills/-/insights/unit-test) [![CNB PR 覆盖率](https://cnb.cool/dtapp/skills/-/badge/git/latest/testing/unit/coverage-pr)](https://cnb.cool/dtapp/skills/-/insights/unit-test) ![CNB 下载请求数](https://img.shields.io/endpoint?label=CNB%20下载请求数&url=https%3A%2F%2Fdl-stats.dtapp.net%2Fshields.json%3Fplatform%3Dcnb.cool%26repository%3D%2Fdtapp%2Fskills)

[![CNB 漏洞](https://cnb.cool/dtapp/skills/-/badge/security/vulnerability/full/latest)](https://cnb.cool/dtapp/skills/-/badge/security/vulnerability/full/latest.link) [![CNB 漏洞差异](https://cnb.cool/dtapp/skills/-/badge/security/vulnerability/diff/latest)](https://cnb.cool/dtapp/skills/-/badge/security/vulnerability/diff/latest.link) [![CNB 许可证](https://cnb.cool/dtapp/skills/-/badge/security/license/full/latest)](https://cnb.cool/dtapp/skills/-/badge/security/license/full/latest.link) [![CNB 许可证差异](https://cnb.cool/dtapp/skills/-/badge/security/license/diff/latest)](https://cnb.cool/dtapp/skills/-/badge/security/license/diff/latest.link) [![CNB SBOM](https://cnb.cool/dtapp/skills/-/badge/security/sbom/full/latest)](https://cnb.cool/dtapp/skills/-/badge/security/sbom/full/latest.link) [![CNB SBOM 差异](https://cnb.cool/dtapp/skills/-/badge/security/sbom/diff/latest)](https://cnb.cool/dtapp/skills/-/badge/security/sbom/diff/latest.link)

[![CNB 访问](https://commit.cool/badge/visit/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 发布频率](https://commit.cool/badge/release/frequency/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 距离上次发布](https://commit.cool/badge/release/last/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 发布总数](https://commit.cool/badge/releases/dtapp/skills)](https://cnb.cool/dtapp/skills/-/releases) [![CNB 仓库年龄](https://commit.cool/badge/age/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 主要语言](https://commit.cool/badge/language/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 仓库许可证](https://commit.cool/badge/license/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB 存储用量](https://commit.cool/badge/usage/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB Git 存储](https://commit.cool/badge/usage/git/dtapp/skills)](https://cnb.cool/dtapp/skills) [![CNB LFS 存储](https://commit.cool/badge/usage/lfs/dtapp/skills)](https://cnb.cool/dtapp/skills)

[![GitHub 版本](https://img.shields.io/github/v/release/dtapps/skills?style=flat-square)](https://github.com/dtapps/skills/releases) [![GitHub 下载量](https://img.shields.io/github/downloads/dtapps/skills/total?style=flat-square)](https://github.com/dtapps/skills/releases) [![GitHub Star 数](https://img.shields.io/github/stars/dtapps/skills?style=flat-square)](https://github.com/dtapps/skills/stargazers) [![GitHub Fork 数](https://img.shields.io/github/forks/dtapps/skills?style=flat-square)](https://github.com/dtapps/skills/network/members) [![GitHub 许可证](https://img.shields.io/github/license/dtapps/skills?style=flat-square)](LICENSE) [![GitHub 支持平台](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-007EC6?style=flat-square)](https://github.com/dtapps/skills) [![GitHub 构建](https://img.shields.io/github/actions/workflow/status/dtapps/skills/push.yml?style=flat-square)](https://github.com/dtapps/skills/actions/workflows/push.yml) [![Codecov](https://img.shields.io/codecov/c/github/dtapps/skills?style=flat-square&label=Codecov&logo=codecov)](https://codecov.io/gh/dtapps/skills) [![Coveralls](https://img.shields.io/coveralls/github/dtapps/skills?style=flat-square&label=Coveralls&logo=coveralls)](https://coveralls.io/github/dtapps/skills) ![GitHub 下载请求数](https://img.shields.io/endpoint?label=GitHub%20下载请求数&url=https%3A%2F%2Fdl-stats.dtapp.top%2Fshields.json%3Fplatform%3Dgithub.com%26repository%3D%2Fdtapps%2Fskills)

[![GitHub 访问](https://visitor-badge.laobi.icu/badge?page_id=github/dtapps/skills)](https://github.com/dtapps/skills) [![GitHub 距离上次发布](https://img.shields.io/github/release-date/dtapps/skills)](https://github.com/dtapps/skills) [![GitHub 仓库年龄](https://img.shields.io/github/created-at/dtapps/skills)](https://github.com/dtapps/skills) [![GitHub 主要语言](https://img.shields.io/github/languages/top/dtapps/skills)](https://github.com/dtapps/skills) [![GitHub 存储用量](https://img.shields.io/github/repo-size/dtapps/skills)](https://github.com/dtapps/skills)

</div>

## 功能特性

- **生成技能包**：通过交互式界面生成 SKILL.md 文件
- **执行命令**：执行技能包中的命令（百度百科搜索、百度搜索、腾讯新闻、腾讯天气）
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
dtskills version
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
dtskills version
dtskills generate
dtskills upgrade
dtskills help
```

**注意**：通过 deb/rpm 安装后，二进制文件位于 `/usr/bin/dtskills`；配置文件位于 `~/.config/dtskills/config.yaml`，可用 `dtskills config init` 初始化。

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

#### 腾讯新闻

```bash
./dtskills tencent-news [command] [flags]
```

**子命令：**

- `search` - 搜索新闻
- `hot` - 获取热点新闻
- `morning` - 获取新闻早报
- `evening` - 获取新闻晚报

**示例：**

```bash
# 搜索新闻
./dtskills tencent-news search -q 人工智能

# 获取热点新闻
./dtskills tencent-news hot

# 获取新闻早报
./dtskills tencent-news morning

# 获取新闻晚报
./dtskills tencent-news evening
```

#### 腾讯天气

```bash
./dtskills tencent-weather [command] [flags]
```

**子命令：**

- `current` - 获取实时天气
- `forecast` - 获取天气预报
- `airquality` - 获取空气质量
- `ip2city` - 根据 IP 获取城市信息

**示例：**

```bash
# 获取北京实时天气（adcode: 110000）
./dtskills tencent-weather current -a 110000

# 获取北京天气预报
./dtskills tencent-weather forecast -a 110000

# 获取北京空气质量
./dtskills tencent-weather airquality -a 110000

# 根据 IP 获取城市信息
./dtskills tencent-weather ip2city
```

### 自动补全

支持为 Bash、Zsh 和 PowerShell 生成自动补全脚本，让你在输入命令时按 `Tab` 键自动补全命令和参数。

**快速安装（Debian/Ubuntu）**：

```bash
sudo dtskills completion bash > /etc/bash_completion.d/dtskills
source /etc/bash_completion.d/dtskills
```

更多平台的安装方法（macOS、Windows）和详细说明，请参阅 [自动补全文档](docs/completion.md)。

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

配置文件 `config.yaml` 位于工作目录下（由 `DTSKILLS_WORK_DIR` 环境变量指定，默认为当前目录或用户配置目录 `~/.config/dtskills/`）。

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

| 变量名                 | 说明                             | 示例           |
|------------------------|----------------------------------|----------------|
| `DTSKILLS_WORK_DIR`    | 工作目录（存储技能包和统计数据） | `./`           |
| `BAIDU_API_KEY`        | 百度 API Key                     | `your-api-key` |
| `TENCENT_NEWS_API_KEY` | 腾讯新闻/天气 API Key            | `your-api-key` |

## 技术栈

- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [charmbracelet/huh](https://github.com/charmbracelet/huh) - 交互式表单
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - 样式渲染
- [minio/selfupdate](https://github.com/minio/selfupdate) - 自升级功能

## 仓库

| 平台   | 地址                             |
|--------|----------------------------------|
| CNB    | https://cnb.cool/dtapp/skills    |
| GitHub | https://github.com/dtapps/skills |

## 许可证

见 [LICENSE](LICENSE) 文件。

## Star 走势

<details>
<summary><b>点击展开查看图表</b></summary>

**CNB**

[![CNB Star 走势](https://commit.cool/stars/dtapp/skills)](https://cnb.cool/dtapp/skills)

**GitHub**

[![Star 走势](https://api.star-history.com/chart?repos=dtapps/skills&type=date&legend=top-left)](https://www.star-history.com/?repos=dtapps%2Fskills&type=date&legend=top-left)

</details>

## 下载量统计图

<details>
<summary><b>点击展开查看图表</b></summary>

[![下载量统计图](https://release-monitor.com/chart/dtapps/skills.svg?stable=5)](https://release-monitor.com/#/dtapps/skills)

</details>
