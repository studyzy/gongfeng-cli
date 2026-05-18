---
name: gongfeng
description: 腾讯工蜂 (git.code.tencent.com) 操作，涉及项目、合并请求、分支、提交、缺陷、代码评审等。
---

面向 AI Agent 的腾讯工蜂命令行工具。通过 `gongfeng` 命令与工蜂平台交互，所有输出针对最小 token 消耗优化。

## 安装

```bash
go install github.com/studyzy/gongfeng-cli/cmd/gongfeng@latest
```

## 认证

```bash
# 交互式登录持久化凭据
gongfeng auth login --token <your_token>

# 或通过命令行标志
gongfeng --token <your_token> ...
```

凭据优先级：CLI flags > `./.gongfeng.json` > `~/.gongfeng.json`

## 命令参考

{{COMMAND_REFERENCE}}
