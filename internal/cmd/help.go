package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// groupPriority 定义命令分组的显示优先级，数字越小越靠前
var groupPriority = map[string]int{
	"project": 1,
	"mr":      2,
	"branch":  3,
	"commit":  4,
	"issue":   5,
	"review":  6,
}

// specLine 表示一条命令参考行
type specLine struct {
	group string
	text  string
}

// buildSpecLines 遍历命令树，为每个叶子命令生成参考行，并按优先级排序
func buildSpecLines(root *cobra.Command) []specLine {
	var lines []specLine
	walkSpecCommands(root, "", "", &lines)
	sortSpecLines(lines)
	return lines
}

// sortSpecLines 按 groupPriority 对参考行排序
func sortSpecLines(lines []specLine) {
	sort.SliceStable(lines, func(i, j int) bool {
		pi := getGroupPriority(lines[i].group)
		pj := getGroupPriority(lines[j].group)
		if pi != pj {
			return pi < pj
		}
		return lines[i].group < lines[j].group
	})
}

// getGroupPriority 返回分组的显示优先级
func getGroupPriority(group string) int {
	if p, ok := groupPriority[group]; ok {
		return p
	}
	return 100
}

// walkSpecCommands 递归遍历命令树，收集叶子命令的参考行
func walkSpecCommands(cmd *cobra.Command, prefix string, group string, lines *[]specLine) {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}

		fullPath := child.Name()
		if prefix != "" {
			fullPath = prefix + " " + child.Name()
		}

		currentGroup := group
		if currentGroup == "" {
			currentGroup = child.Name()
		}

		if child.HasSubCommands() {
			walkSpecCommands(child, fullPath, currentGroup, lines)
		} else {
			line := commandToLine(child, fullPath)
			*lines = append(*lines, specLine{group: currentGroup, text: line})
		}
	}
}

// commandToLine 将 Cobra 命令转换为一行紧凑参考文本
func commandToLine(cmd *cobra.Command, path string) string {
	var b strings.Builder
	b.WriteString("gongfeng ")
	b.WriteString(path)

	if argName := extractArgName(cmd.Use); argName != "" {
		b.WriteString(" <")
		b.WriteString(argName)
		b.WriteString(">")
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		if isGlobalFlag(f.Name) {
			return
		}
		b.WriteString(" ")
		b.WriteString(formatFlag(f))
	})

	if cmd.Short != "" {
		b.WriteString("  # ")
		b.WriteString(cmd.Short)
	}

	return b.String()
}

// formatFlag 将一个 flag 格式化为紧凑文本
func formatFlag(f *pflag.Flag) string {
	hint := extractEnumHint(f.Usage)
	if isFlagRequired(f) {
		if hint != "" {
			return "--" + f.Name + "=<" + hint + ">"
		}
		return "--" + f.Name + "=<" + f.Name + ">"
	}
	if hint != "" {
		return "[--" + f.Name + "=<" + hint + ">]"
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
		return "[--" + f.Name + "=" + f.DefValue + "]"
	}
	return "[--" + f.Name + "]"
}

// extractEnumHint 从 Usage 文本的全角括号（）中提取枚举提示
func extractEnumHint(usage string) string {
	start := strings.Index(usage, "（")
	end := strings.LastIndex(usage, "）")
	if start < 0 || end <= start {
		return ""
	}
	content := usage[start+len("（") : end]
	content = strings.TrimSuffix(content, "，必需")
	content = strings.TrimSuffix(content, "，必填")
	return content
}

// isFlagRequired 判断标志是否为必填
func isFlagRequired(f *pflag.Flag) bool {
	return strings.Contains(f.Usage, "必需") || strings.Contains(f.Usage, "必填")
}

// isGlobalFlag 判断是否为全局标志
func isGlobalFlag(name string) bool {
	switch name {
	case "token", "base-url", "project-id", "pretty", "json":
		return true
	default:
		return false
	}
}

// printSpecOutput 输出完整的参考卡文本
func printSpecOutput(w *os.File, root *cobra.Command, lines []specLine) {
	fmt.Fprintf(w, "gongfeng - %s\n", root.Short)
	fmt.Fprintln(w, "Global: [--project-id=<id>] [--json] [--pretty]")

	lastGroup := ""
	for _, l := range lines {
		if l.group != lastGroup {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "# %s\n", l.group)
			lastGroup = l.group
		}
		fmt.Fprintln(w, l.text)
	}
}

// extractArgName 从 Use 字段提取位置参数名
func extractArgName(use string) string {
	start := -1
	for i, c := range use {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start >= 0 {
			return use[start:i]
		}
	}
	return ""
}
