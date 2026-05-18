// Package cmd 定义了 gongfeng-cli 的所有 Cobra 命令
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/config"
	"github.com/studyzy/gongfeng-cli/internal/output"
)

var (
	// 全局标志
	flagProjectID string
	flagPretty    bool
	flagJSON      bool
	flagToken     string
	flagBaseURL   string

	// 全局共享的客户端
	apiClient *gongfeng.Client
)

// rootCmd 是 CLI 的根命令
var rootCmd = &cobra.Command{
	Use:   "gongfeng",
	Short: "面向 AI Agent 的腾讯工蜂命令行工具",
	Long:  "gongfeng-cli 是一个面向 AI Agent 的腾讯工蜂命令行工具，通过工蜂 API 实现代码托管核心操作。输出针对最小 token 消耗优化。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// auth login 命令不需要预加载配置和客户端
		if cmd.Name() == "login" {
			return nil
		}
		// skill 子命令（如 skill init）不需要认证
		if parent := cmd.Parent(); parent != nil && parent.Name() == "skill" {
			return nil
		}
		// --version 不需要认证
		if v, _ := cmd.Flags().GetBool("version"); v {
			return nil
		}
		return initClientAndConfig(cmd)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// 根命令自定义 help：输出紧凑参考卡
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd != rootCmd {
			defaultHelp(cmd, args)
			return
		}
		lines := buildSpecLines(rootCmd)
		printSpecOutput(os.Stdout, rootCmd, lines)
	})

	rootCmd.PersistentFlags().StringVar(&flagProjectID, "project-id", "", "项目 ID 或 namespace/project 路径（覆盖本地配置）")
	rootCmd.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "输出带缩进的 JSON，仅供人类阅读，AI Agent 不应使用（浪费 token）")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "强制 JSON 输出（列表默认已是 JSON，详情默认 Markdown 更省 token）")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "工蜂 Private Token")
	rootCmd.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "工蜂实例 URL（默认 https://git.code.tencent.com/）")
}

// initClientAndConfig 初始化配置和 API 客户端
func initClientAndConfig(cmd *cobra.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// 命令行标志覆盖配置
	token := flagToken
	if token == "" {
		token = cfg.Token
	}

	if token == "" {
		output.PrintError(os.Stderr, "authentication_required",
			"No valid credentials found",
			"Run 'gongfeng auth login --token <token>'. "+
				"You can also set GONGFENG_TOKEN environment variable.")
		os.Exit(output.ExitAuthError)
	}

	// 构建客户端选项
	var opts []gongfeng.ClientOptionFunc
	baseURL := flagBaseURL
	if baseURL == "" {
		baseURL = cfg.BaseURL
	}
	if baseURL != "" {
		opts = append(opts, gongfeng.WithBaseURL(baseURL))
	}

	client, err := gongfeng.NewClient(token, opts...)
	if err != nil {
		output.PrintError(os.Stderr, "client_error",
			"Failed to create API client: "+err.Error(),
			"Check your token and base URL configuration.")
		os.Exit(output.ExitAPIError)
	}
	apiClient = client

	// project-id 标志覆盖配置
	if flagProjectID == "" {
		flagProjectID = cfg.ProjectID
	}

	return nil
}

// useJSONOutput 判断是否应使用 JSON 格式输出，--pretty 隐含 --json
func useJSONOutput() bool {
	return flagJSON || flagPretty
}

// printDetail 输出单条详情，默认 Markdown 格式，--json/--pretty 时输出 JSON
func printDetail(data interface{}, bodyField string) error {
	if useJSONOutput() {
		return output.PrintJSON(os.Stdout, data, !flagPretty)
	}
	return output.PrintMarkdown(os.Stdout, data, bodyField)
}

// printSuccessResponse 输出操作成功的精简响应
func printSuccessResponse(id, url, projectID string) error {
	resp := &output.SuccessResponse{
		Success:   true,
		ID:        id,
		URL:       url,
		ProjectID: projectID,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

// requireProjectID 检查 project-id 是否已配置，未配置时输出错误并退出
func requireProjectID() {
	if flagProjectID == "" {
		output.PrintError(os.Stderr, "project_required",
			"No project ID configured",
			"Use --project-id flag or run 'gongfeng auth login' and set project_id in config.")
		os.Exit(output.ExitParamError)
	}
}

// projectID 返回用于 SDK 调用的项目 ID（优先尝试整数，否则作为字符串路径）
func projectID() interface{} {
	if id, err := strconv.Atoi(flagProjectID); err == nil {
		return id
	}
	return flagProjectID
}

// exitWithAPIError 输出 API 错误并退出
func exitWithAPIError(err error) {
	output.PrintError(os.Stderr, "api_error", err.Error(), "")
	os.Exit(output.ExitAPIError)
}

// exitWithParamError 输出参数错误并退出
func exitWithParamError(message, hint string) {
	output.PrintError(os.Stderr, "missing_parameter", message, hint)
	os.Exit(output.ExitParamError)
}

// atoi 将字符串转换为整数，失败时输出错误并退出
func atoi(s string, name string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		exitWithParamError(fmt.Sprintf("invalid %s: %s", name, s), fmt.Sprintf("%s must be an integer", name))
	}
	return v
}
