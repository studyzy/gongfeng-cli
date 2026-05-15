package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// mr 子命令独有的 flag 变量
var (
	mrFlagState        string
	mrFlagOrderBy      string
	mrFlagSort         string
	mrFlagPage         int
	mrFlagPerPage      int
	mrFlagSourceBranch string
	mrFlagTargetBranch string
	mrFlagTitle        string
	mrFlagDescription  string
	mrFlagAssigneeID   int
	mrFlagReviewers    string
	mrFlagApproverRule string
	mrFlagStateEvent   string
	mrFlagMergeMsg     string
)

var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "合并请求管理",
}

var mrListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取 MR 列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListMergeRequestsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    mrFlagPage,
				PerPage: mrFlagPerPage,
			},
		}
		if mrFlagState != "" {
			opts.State = gongfeng.Ptr(mrFlagState)
		}
		if mrFlagOrderBy != "" {
			opts.OrderBy = gongfeng.Ptr(mrFlagOrderBy)
		}
		if mrFlagSort != "" {
			opts.Sort = gongfeng.Ptr(mrFlagSort)
		}
		mrs, _, err := apiClient.MergeRequests.ListMergeRequests(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, mrs, !flagPretty)
	},
}

var mrShowCmd = &cobra.Command{
	Use:   "show <mr_id>",
	Short: "获取 MR 详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		mr, _, err := apiClient.MergeRequests.GetMergeRequest(context.Background(), projectID(), mrID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(mr, "description")
	},
}

var mrCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建 MR",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if mrFlagSourceBranch == "" {
			exitWithParamError("--source-branch is required", "Specify the source branch for the merge request")
		}
		if mrFlagTargetBranch == "" {
			exitWithParamError("--target-branch is required", "Specify the target branch for the merge request")
		}
		if mrFlagTitle == "" {
			exitWithParamError("--title is required", "Specify the title for the merge request")
		}
		opts := &gongfeng.CreateMergeRequestOptions{
			SourceBranch: gongfeng.Ptr(mrFlagSourceBranch),
			TargetBranch: gongfeng.Ptr(mrFlagTargetBranch),
			Title:        gongfeng.Ptr(mrFlagTitle),
		}
		if mrFlagDescription != "" {
			opts.Description = gongfeng.Ptr(mrFlagDescription)
		}
		if mrFlagAssigneeID != 0 {
			opts.AssigneeID = gongfeng.Ptr(mrFlagAssigneeID)
		}
		if mrFlagReviewers != "" {
			opts.Reviewers = gongfeng.Ptr(mrFlagReviewers)
		}
		if mrFlagApproverRule != "" {
			opts.ApproverRule = gongfeng.Ptr(mrFlagApproverRule)
		}
		mr, _, err := apiClient.MergeRequests.CreateMergeRequest(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(mr, "description")
	},
}

var mrUpdateCmd = &cobra.Command{
	Use:   "update <mr_id>",
	Short: "更新 MR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.UpdateMergeRequestOptions{}
		if mrFlagTitle != "" {
			opts.Title = gongfeng.Ptr(mrFlagTitle)
		}
		if mrFlagDescription != "" {
			opts.Description = gongfeng.Ptr(mrFlagDescription)
		}
		if mrFlagTargetBranch != "" {
			opts.TargetBranch = gongfeng.Ptr(mrFlagTargetBranch)
		}
		if mrFlagAssigneeID != 0 {
			opts.AssigneeID = gongfeng.Ptr(mrFlagAssigneeID)
		}
		if mrFlagStateEvent != "" {
			opts.StateEvent = gongfeng.Ptr(mrFlagStateEvent)
		}
		mr, _, err := apiClient.MergeRequests.UpdateMergeRequest(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(mr, "description")
	},
}

var mrAcceptCmd = &cobra.Command{
	Use:   "accept <mr_id>",
	Short: "合并 MR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.AcceptMergeRequestOptions{}
		if mrFlagMergeMsg != "" {
			opts.MergeCommitMessage = gongfeng.Ptr(mrFlagMergeMsg)
		}
		mr, _, err := apiClient.MergeRequests.AcceptMergeRequest(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(mr, "description")
	},
}

var mrChangesCmd = &cobra.Command{
	Use:   "changes <mr_id>",
	Short: "获取 MR 代码变更",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		changes, _, err := apiClient.MergeRequests.GetMergeRequestChanges(context.Background(), projectID(), mrID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, changes, !flagPretty)
	},
}

var mrCommitsCmd = &cobra.Command{
	Use:   "commits <mr_id>",
	Short: "获取 MR 提交列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.ListMergeRequestCommitsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    mrFlagPage,
				PerPage: mrFlagPerPage,
			},
		}
		commits, _, err := apiClient.MergeRequests.ListMergeRequestCommits(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, commits, !flagPretty)
	},
}

var mrCommentsCmd = &cobra.Command{
	Use:   "comments <mr_id>",
	Short: "获取 MR 评论列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.ListMRCommentsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    mrFlagPage,
				PerPage: mrFlagPerPage,
			},
		}
		comments, _, err := apiClient.MergeRequests.ListMRComments(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, comments, !flagPretty)
	},
}

func init() {
	// list flags
	mrListCmd.Flags().StringVar(&mrFlagState, "state", "", "MR 状态过滤（opened/closed/merged/all）")
	mrListCmd.Flags().StringVar(&mrFlagOrderBy, "order-by", "", "排序字段")
	mrListCmd.Flags().StringVar(&mrFlagSort, "sort", "", "排序方向（asc/desc）")
	mrListCmd.Flags().IntVar(&mrFlagPage, "page", 0, "页码")
	mrListCmd.Flags().IntVar(&mrFlagPerPage, "per-page", 0, "每页数量")

	// create flags
	mrCreateCmd.Flags().StringVar(&mrFlagSourceBranch, "source-branch", "", "源分支（必需）")
	mrCreateCmd.Flags().StringVar(&mrFlagTargetBranch, "target-branch", "", "目标分支（必需）")
	mrCreateCmd.Flags().StringVar(&mrFlagTitle, "title", "", "MR 标题（必需）")
	mrCreateCmd.Flags().StringVar(&mrFlagDescription, "description", "", "MR 描述")
	mrCreateCmd.Flags().IntVar(&mrFlagAssigneeID, "assignee-id", 0, "指派人 ID")
	mrCreateCmd.Flags().StringVar(&mrFlagReviewers, "reviewers", "", "评审人")
	mrCreateCmd.Flags().StringVar(&mrFlagApproverRule, "approver-rule", "", "审批规则")

	// update flags
	mrUpdateCmd.Flags().StringVar(&mrFlagTitle, "title", "", "MR 标题")
	mrUpdateCmd.Flags().StringVar(&mrFlagDescription, "description", "", "MR 描述")
	mrUpdateCmd.Flags().StringVar(&mrFlagTargetBranch, "target-branch", "", "目标分支")
	mrUpdateCmd.Flags().IntVar(&mrFlagAssigneeID, "assignee-id", 0, "指派人 ID")
	mrUpdateCmd.Flags().StringVar(&mrFlagStateEvent, "state-event", "", "状态事件（close/reopen/merge）")

	// accept flags
	mrAcceptCmd.Flags().StringVar(&mrFlagMergeMsg, "merge-commit-message", "", "合并提交信息")

	// commits flags
	mrCommitsCmd.Flags().IntVar(&mrFlagPage, "page", 0, "页码")
	mrCommitsCmd.Flags().IntVar(&mrFlagPerPage, "per-page", 0, "每页数量")

	// comments flags
	mrCommentsCmd.Flags().IntVar(&mrFlagPage, "page", 0, "页码")
	mrCommentsCmd.Flags().IntVar(&mrFlagPerPage, "per-page", 0, "每页数量")

	mrCmd.AddCommand(mrListCmd, mrShowCmd, mrCreateCmd, mrUpdateCmd, mrAcceptCmd, mrChangesCmd, mrCommitsCmd, mrCommentsCmd)
	rootCmd.AddCommand(mrCmd)
}
