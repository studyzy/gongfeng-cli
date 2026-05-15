package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// commit-status list flags
var (
	csFlagPage    int
	csFlagPerPage int
	csFlagRef     string
	csFlagStage   string
	csFlagName    string
	csFlagAll     bool
)

// commit-status create flags
var (
	csFlagState       string
	csFlagCreateRef   string
	csFlagCreateName  string
	csFlagTargetURL   string
	csFlagDescription string
	csFlagContext      string
)

var commitStatusCmd = &cobra.Command{
	Use:   "commit-status",
	Short: "提交检测状态管理",
}

var commitStatusListCmd = &cobra.Command{
	Use:   "list <sha>",
	Short: "查询提交检测结果列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListCommitStatusesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    csFlagPage,
				PerPage: csFlagPerPage,
			},
		}
		if csFlagRef != "" {
			opts.Ref = gongfeng.Ptr(csFlagRef)
		}
		if csFlagStage != "" {
			opts.Stage = gongfeng.Ptr(csFlagStage)
		}
		if csFlagName != "" {
			opts.Name = gongfeng.Ptr(csFlagName)
		}
		if cmd.Flags().Changed("all") {
			opts.All = gongfeng.Ptr(csFlagAll)
		}

		statuses, _, err := apiClient.CommitStatuses.ListCommitStatuses(ctx, projectID(), args[0], opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, statuses, !flagPretty)
	},
}

var commitStatusCreateCmd = &cobra.Command{
	Use:   "create <sha>",
	Short: "创建提交检测结果",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		if csFlagState == "" {
			exitWithParamError("--state is required", "Provide a state with --state (pending/running/success/failed/canceled)")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateCommitStatusOptions{
			State: gongfeng.Ptr(csFlagState),
		}
		if csFlagCreateRef != "" {
			opts.Ref = gongfeng.Ptr(csFlagCreateRef)
		}
		if csFlagCreateName != "" {
			opts.Name = gongfeng.Ptr(csFlagCreateName)
		}
		if csFlagTargetURL != "" {
			opts.TargetURL = gongfeng.Ptr(csFlagTargetURL)
		}
		if csFlagDescription != "" {
			opts.Description = gongfeng.Ptr(csFlagDescription)
		}
		if csFlagContext != "" {
			opts.Context = gongfeng.Ptr(csFlagContext)
		}

		status, _, err := apiClient.CommitStatuses.CreateCommitStatus(ctx, projectID(), args[0], opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, status, !flagPretty)
	},
}

var commitStatusResultCmd = &cobra.Command{
	Use:   "result <ref>",
	Short: "查询组合检测结果",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		result, _, err := apiClient.CommitStatuses.GetCommitStatusResult(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, result, !flagPretty)
	},
}

func init() {
	// commit-status list flags
	commitStatusListCmd.Flags().IntVar(&csFlagPage, "page", 0, "页码")
	commitStatusListCmd.Flags().IntVar(&csFlagPerPage, "per-page", 0, "每页条数")
	commitStatusListCmd.Flags().StringVar(&csFlagRef, "ref", "", "分支名或 Tag 名")
	commitStatusListCmd.Flags().StringVar(&csFlagStage, "stage", "", "构建阶段")
	commitStatusListCmd.Flags().StringVar(&csFlagName, "name", "", "检测名称")
	commitStatusListCmd.Flags().BoolVar(&csFlagAll, "all", false, "返回所有状态（包括重试）")

	// commit-status create flags
	commitStatusCreateCmd.Flags().StringVar(&csFlagState, "state", "", "状态（必需：pending/running/success/failed/canceled）")
	commitStatusCreateCmd.Flags().StringVar(&csFlagCreateRef, "ref", "", "分支名或 Tag 名")
	commitStatusCreateCmd.Flags().StringVar(&csFlagCreateName, "name", "", "检测名称")
	commitStatusCreateCmd.Flags().StringVar(&csFlagTargetURL, "target-url", "", "目标 URL")
	commitStatusCreateCmd.Flags().StringVar(&csFlagDescription, "description", "", "描述信息")
	commitStatusCreateCmd.Flags().StringVar(&csFlagContext, "context", "", "上下文标识")

	commitStatusCmd.AddCommand(commitStatusListCmd, commitStatusCreateCmd, commitStatusResultCmd)
	rootCmd.AddCommand(commitStatusCmd)
}
