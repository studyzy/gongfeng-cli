# gongfeng-cli

面向 AI Agent 的腾讯工蜂命令行工具，通过工蜂 API 实现代码托管核心操作。

## 安装

### 方式一：go install（推荐）

```bash
go install github.com/studyzy/gongfeng-cli/cmd/gongfeng@latest
```

### 方式二：从源码构建并安装

```bash
git clone git@github.com:studyzy/gongfeng-cli.git
cd gongfeng-cli
make install   # 编译并安装到 $GOPATH/bin
```

### 方式三：仅构建二进制

```bash
git clone git@github.com:studyzy/gongfeng-cli.git
cd gongfeng-cli
make build     # 在当前目录生成 ./gongfeng
```

## 认证

使用工蜂 Private Token 进行认证。

```bash
# 命令行登录（写入 ~/.gongfeng.json）
gongfeng auth login --token <your_private_token>

# 写入当前目录 .gongfeng.json
gongfeng auth login --token <your_private_token> --local

# 或设置环境变量
export GONGFENG_TOKEN=<your_private_token>
```

凭据也可以直接写入配置文件 `~/.gongfeng.json` 或当前目录的 `.gongfeng.json`。

**凭据优先级**：CLI flags > 环境变量 > `./.gongfeng.json` > `~/.gongfeng.json`

### 自定义工蜂站点地址

如需连接私有部署的工蜂实例，可通过环境变量或配置文件指定：

```bash
# 环境变量
export GONGFENG_BASE_URL=https://your-gongfeng.example.com/
```

或写入配置文件：

```json
{
  "token": "your-token",
  "base_url": "https://your-gongfeng.example.com/"
}
```

| 配置项 | 环境变量 | JSON 字段 | 默认值 |
|--------|----------|-----------|--------|
| API 地址 | `GONGFENG_BASE_URL` | `base_url` | `https://code.tencent.com/` |

## 基本用法

```bash
# 查询项目列表
gongfeng project list

# 查看项目详情（支持 namespace/project 路径或数字 ID）
gongfeng project show --project-id namespace/myproject

# 查询 MR 列表
gongfeng mr list

# 创建 MR
gongfeng mr create --source-branch feature/xxx --target-branch main --title "feat: 新功能"

# 合并 MR
gongfeng mr accept 42

# 查询分支列表
gongfeng branch list

# 查询提交列表
gongfeng commit list

# 查询 Issue 列表
gongfeng issue list

# 创建 Issue
gongfeng issue create --title "Bug: 登录失败" --description "详情..."

# 查看所有命令参考（AI 自发现）
gongfeng --help
```

## 命令一览

```
gongfeng
├── auth           login --token <token> [--local]
├── project        list | show | create | search | members
├── mr             list | show | create | update | accept | changes | commits | comments
├── commit         list | show | comments | refs | diff
├── branch         list | show | create | delete | protect
├── tag            list | show | create | delete
├── issue          list | show | create | update | close | reopen
├── release        list | create
├── review         list | show | create | update | comment
├── commit-status  list | create | update
├── repository     info
├── group          list | show | create | delete
├── namespace      list
├── label          list | create | delete | subscribe
├── milestone      list | show | create | update | close
├── note           list | create | update | delete
├── fork           create | list
├── user           info
├── watcher        list | add | remove
├── webhook        list | create | delete
└── version
```

## 全局标志

| 标志 | 说明 |
|------|------|
| `--project-id <id>` | 项目 ID 或 `namespace/project` 路径（覆盖本地配置），支持整数 ID 或字符串路径两种格式 |
| `--token <token>` | 工蜂 Private Token（覆盖配置文件） |
| `--base-url <url>` | 工蜂实例 URL（覆盖配置文件） |
| `--pretty` | 输出带缩进的 JSON，仅供人类阅读；AI Agent 不应使用（浪费 token） |
| `--json` | 强制 JSON 输出（列表默认已是 JSON，详情默认 Markdown 更省 token） |

## 配置文件

配置文件为 `~/.gongfeng.json`（全局）或 `./.gongfeng.json`（项目级），两者可共存，项目级优先。

```json
{
  "token": "your-private-token",
  "project_id": "namespace/project",
  "base_url": "https://code.tencent.com/"
}
```

## 输出格式

- **列表命令**：默认输出紧凑 JSON（无缩进），最小化 token 消耗
- **详情命令**：默认输出 Markdown 格式，加 `--json` 可切换为 JSON
- **人类阅读**：使用 `--pretty` 输出带缩进的 JSON

## SDK

工蜂 Go SDK 已独立为单独的模块，可直接引入使用：

```bash
go get github.com/studyzy/gongfeng-sdk-go@latest
```

详见 [gongfeng-sdk-go](https://github.com/studyzy/gongfeng-sdk-go)。

## 开发

```bash
make build      # 构建
make install    # 安装到 $GOPATH/bin
make test       # 运行测试
make coverage   # 测试覆盖率报告
make lint       # 代码检查
make fmt        # 代码格式化
make clean      # 清理构建产物
```

## 许可证

Apache License 2.0
