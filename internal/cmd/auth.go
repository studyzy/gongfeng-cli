package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/studyzy/gongfeng-cli/internal/config"
	"github.com/studyzy/gongfeng-cli/internal/output"
)

var flagLocal bool

// authCmd 是 auth 父命令
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "认证管理",
}

// loginCmd 是 auth login 子命令
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "保存 Private Token 到配置文件",
	Long:  "使用 --token 传入工蜂 Private Token。默认写入 ~/.gongfeng.json，使用 --local 写入当前目录。",
	RunE:  runLogin,
}

func init() {
	loginCmd.Flags().BoolVar(&flagLocal, "local", false, "将凭据写入当前目录的 .gongfeng.json")
	authCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(authCmd)
}

// runLogin 执行登录逻辑
func runLogin(cmd *cobra.Command, args []string) error {
	token := flagToken
	if token == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"Provide --token",
			"Usage: gongfeng auth login --token <private-token>")
		os.Exit(output.ExitParamError)
		return nil
	}

	cfg := &config.Config{Token: token}

	configPath, err := config.GetConfigPath(flagLocal)
	if err != nil {
		output.PrintError(os.Stderr, "config_error",
			"Failed to determine config path: "+err.Error(),
			"Check file system permissions.")
		os.Exit(output.ExitAPIError)
		return nil
	}

	if err := config.SaveConfig(cfg, configPath); err != nil {
		output.PrintError(os.Stderr, "config_error",
			"Failed to save config: "+err.Error(),
			"Check file system permissions for "+configPath)
		os.Exit(output.ExitAPIError)
		return nil
	}

	return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
}
