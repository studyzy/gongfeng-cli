package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// review 子命令独有的 flag 变量
var (
	rvFlagReviewerID          int
	rvFlagNecessaryReviewerID int
	rvFlagMessage             string
	rvFlagState               string
	rvFlagOrderBy             string
	rvFlagSort                string
	rvFlagPage                int
	rvFlagPerPage             int
	rvFlagAuthorID            int
	rvFlagSourceCommit        string
	rvFlagSourceBranch        string
	rvFlagTargetBranch        string
	rvFlagTargetCommit        string
	rvFlagTargetProjectID     int
	rvFlagReviewerIDs         string
	rvFlagNecessaryReviewerIDs string
	rvFlagApproverRule        int
	rvFlagNecessaryApproverRule int
	rvFlagTitle               string
	rvFlagDescription         string
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "代码评审管理",
}

// --- MR 评审 ---

var reviewInviteCmd = &cobra.Command{
	Use:   "invite <mr_id>",
	Short: "邀请 MR 评审人",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		if rvFlagReviewerID == 0 {
			exitWithParamError("--reviewer-id is required", "Specify the reviewer user ID")
		}
		opts := &gongfeng.InviteMRReviewerOptions{
			ReviewerID: gongfeng.Ptr(rvFlagReviewerID),
		}
		if rvFlagNecessaryReviewerID != 0 {
			opts.NecessaryReviewerID = gongfeng.Ptr(rvFlagNecessaryReviewerID)
		}
		_, err := apiClient.Reviews.InviteMRReviewer(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var reviewRemoveCmd = &cobra.Command{
	Use:   "remove <mr_id>",
	Short: "移除 MR 评审人",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		if rvFlagReviewerID == 0 {
			exitWithParamError("--reviewer-id is required", "Specify the reviewer user ID to remove")
		}
		opts := &gongfeng.RemoveMRReviewerOptions{
			ReviewerID: gongfeng.Ptr(rvFlagReviewerID),
		}
		_, err := apiClient.Reviews.RemoveMRReviewer(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var reviewMRShowCmd = &cobra.Command{
	Use:   "mr-show <mr_id>",
	Short: "获取 MR 评审信息",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		review, _, err := apiClient.Reviews.GetMRReview(context.Background(), projectID(), mrID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(review, "description")
	},
}

var reviewApproveCmd = &cobra.Command{
	Use:   "approve <mr_id>",
	Short: "通过 MR 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.SubmitMRReviewSummaryOptions{
			ReviewerEvent: gongfeng.Ptr("approve"),
		}
		if rvFlagMessage != "" {
			opts.Summary = gongfeng.Ptr(rvFlagMessage)
		}
		_, _, err := apiClient.Reviews.SubmitMRReviewSummary(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var reviewRejectCmd = &cobra.Command{
	Use:   "reject <mr_id>",
	Short: "拒绝 MR 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.SubmitMRReviewSummaryOptions{
			ReviewerEvent: gongfeng.Ptr("reject"),
		}
		if rvFlagMessage != "" {
			opts.Summary = gongfeng.Ptr(rvFlagMessage)
		}
		_, _, err := apiClient.Reviews.SubmitMRReviewSummary(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var reviewMRReopenCmd = &cobra.Command{
	Use:   "mr-reopen <mr_id>",
	Short: "重置 MR 评审状态",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		_, _, err := apiClient.Reviews.ReopenMRReview(context.Background(), projectID(), mrID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var reviewMRCancelCmd = &cobra.Command{
	Use:   "mr-cancel <mr_id>",
	Short: "取消 MR 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		_, err := apiClient.Reviews.CancelMRReview(context.Background(), projectID(), mrID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

// --- Commit 评审 ---

var commitReviewListCmd = &cobra.Command{
	Use:   "commit-list",
	Short: "获取 Commit 评审列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListCommitReviewsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    rvFlagPage,
				PerPage: rvFlagPerPage,
			},
		}
		if rvFlagState != "" {
			opts.State = gongfeng.Ptr(rvFlagState)
		}
		if rvFlagOrderBy != "" {
			opts.OrderBy = gongfeng.Ptr(rvFlagOrderBy)
		}
		if rvFlagSort != "" {
			opts.Sort = gongfeng.Ptr(rvFlagSort)
		}
		if rvFlagAuthorID != 0 {
			opts.AuthorID = gongfeng.Ptr(rvFlagAuthorID)
		}
		reviews, _, err := apiClient.Reviews.ListCommitReviews(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, reviews, !flagPretty)
	},
}

var commitReviewShowCmd = &cobra.Command{
	Use:   "commit-show <review_id>",
	Short: "获取 Commit 评审详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		review, _, err := apiClient.Reviews.GetCommitReview(context.Background(), projectID(), reviewID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(review, "description")
	},
}

var commitReviewCreateCmd = &cobra.Command{
	Use:   "commit-create",
	Short: "创建 Commit 评审",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rvFlagTitle == "" {
			exitWithParamError("--title is required", "Specify the title for the commit review")
		}
		opts := &gongfeng.CreateCommitReviewOptions{
			Title: gongfeng.Ptr(rvFlagTitle),
		}
		if rvFlagSourceCommit != "" {
			opts.SourceCommit = gongfeng.Ptr(rvFlagSourceCommit)
		}
		if rvFlagSourceBranch != "" {
			opts.SourceBranch = gongfeng.Ptr(rvFlagSourceBranch)
		}
		if rvFlagTargetBranch != "" {
			opts.TargetBranch = gongfeng.Ptr(rvFlagTargetBranch)
		}
		if rvFlagTargetCommit != "" {
			opts.TargetCommit = gongfeng.Ptr(rvFlagTargetCommit)
		}
		if rvFlagTargetProjectID != 0 {
			opts.TargetProjectID = gongfeng.Ptr(rvFlagTargetProjectID)
		}
		if rvFlagDescription != "" {
			opts.Description = gongfeng.Ptr(rvFlagDescription)
		}
		if rvFlagReviewerIDs != "" {
			opts.ReviewerIDs = gongfeng.Ptr(rvFlagReviewerIDs)
		}
		if rvFlagNecessaryReviewerIDs != "" {
			opts.NecessaryReviewerIDs = gongfeng.Ptr(rvFlagNecessaryReviewerIDs)
		}
		if rvFlagApproverRule != 0 {
			opts.ApproverRule = gongfeng.Ptr(rvFlagApproverRule)
		}
		if rvFlagNecessaryApproverRule != 0 {
			opts.NecessaryApproverRule = gongfeng.Ptr(rvFlagNecessaryApproverRule)
		}
		review, _, err := apiClient.Reviews.CreateCommitReview(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(review, "description")
	},
}

var commitReviewUpdateCmd = &cobra.Command{
	Use:   "commit-update <review_id>",
	Short: "更新 Commit 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		opts := &gongfeng.UpdateCommitReviewOptions{}
		if rvFlagTitle != "" {
			opts.Title = gongfeng.Ptr(rvFlagTitle)
		}
		if rvFlagDescription != "" {
			opts.Description = gongfeng.Ptr(rvFlagDescription)
		}
		review, _, err := apiClient.Reviews.UpdateCommitReview(context.Background(), projectID(), reviewID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(review, "description")
	},
}

var commitReviewApproveCmd = &cobra.Command{
	Use:   "commit-approve <review_id>",
	Short: "通过 Commit 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		opts := &gongfeng.SubmitCommitReviewSummaryOptions{
			ReviewerEvent: gongfeng.Ptr("approve"),
		}
		if rvFlagMessage != "" {
			opts.Summary = gongfeng.Ptr(rvFlagMessage)
		}
		_, _, err := apiClient.Reviews.SubmitCommitReviewSummary(context.Background(), projectID(), reviewID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var commitReviewRejectCmd = &cobra.Command{
	Use:   "commit-reject <review_id>",
	Short: "拒绝 Commit 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		opts := &gongfeng.SubmitCommitReviewSummaryOptions{
			ReviewerEvent: gongfeng.Ptr("reject"),
		}
		if rvFlagMessage != "" {
			opts.Summary = gongfeng.Ptr(rvFlagMessage)
		}
		_, _, err := apiClient.Reviews.SubmitCommitReviewSummary(context.Background(), projectID(), reviewID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var commitReviewReopenCmd = &cobra.Command{
	Use:   "commit-reopen <review_id>",
	Short: "重置 Commit 评审状态",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		_, _, err := apiClient.Reviews.ReopenCommitReview(context.Background(), projectID(), reviewID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// invite flags
	reviewInviteCmd.Flags().IntVar(&rvFlagReviewerID, "reviewer-id", 0, "评审人用户 ID（必需）")
	reviewInviteCmd.Flags().IntVar(&rvFlagNecessaryReviewerID, "necessary-reviewer-id", 0, "必要评审人用户 ID")

	// remove flags
	reviewRemoveCmd.Flags().IntVar(&rvFlagReviewerID, "reviewer-id", 0, "评审人用户 ID（必需）")

	// approve flags
	reviewApproveCmd.Flags().StringVar(&rvFlagMessage, "message", "", "审批意见")

	// reject flags
	reviewRejectCmd.Flags().StringVar(&rvFlagMessage, "message", "", "拒绝原因")

	// commit-list flags
	commitReviewListCmd.Flags().StringVar(&rvFlagState, "state", "", "评审状态过滤")
	commitReviewListCmd.Flags().StringVar(&rvFlagOrderBy, "order-by", "", "排序字段")
	commitReviewListCmd.Flags().StringVar(&rvFlagSort, "sort", "", "排序方向（asc/desc）")
	commitReviewListCmd.Flags().IntVar(&rvFlagPage, "page", 0, "页码")
	commitReviewListCmd.Flags().IntVar(&rvFlagPerPage, "per-page", 0, "每页数量")
	commitReviewListCmd.Flags().IntVar(&rvFlagAuthorID, "author-id", 0, "作者 ID 过滤")

	// commit-create flags
	commitReviewCreateCmd.Flags().StringVar(&rvFlagTitle, "title", "", "评审标题（必需）")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagSourceCommit, "source-commit", "", "源提交 SHA")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagSourceBranch, "source-branch", "", "源分支")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagTargetBranch, "target-branch", "", "目标分支")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagTargetCommit, "target-commit", "", "目标提交 SHA")
	commitReviewCreateCmd.Flags().IntVar(&rvFlagTargetProjectID, "target-project-id", 0, "目标项目 ID")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagDescription, "description", "", "评审描述")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagReviewerIDs, "reviewer-ids", "", "评审人 ID 列表（逗号分隔）")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagNecessaryReviewerIDs, "necessary-reviewer-ids", "", "必要评审人 ID 列表（逗号分隔）")
	commitReviewCreateCmd.Flags().IntVar(&rvFlagApproverRule, "approver-rule", 0, "审批规则")
	commitReviewCreateCmd.Flags().IntVar(&rvFlagNecessaryApproverRule, "necessary-approver-rule", 0, "必要审批规则")

	// commit-update flags
	commitReviewUpdateCmd.Flags().StringVar(&rvFlagTitle, "title", "", "评审标题")
	commitReviewUpdateCmd.Flags().StringVar(&rvFlagDescription, "description", "", "评审描述")

	// commit-approve flags
	commitReviewApproveCmd.Flags().StringVar(&rvFlagMessage, "message", "", "审批意见")

	// commit-reject flags
	commitReviewRejectCmd.Flags().StringVar(&rvFlagMessage, "message", "", "拒绝原因")

	reviewCmd.AddCommand(
		reviewInviteCmd, reviewRemoveCmd, reviewMRShowCmd, reviewApproveCmd, reviewRejectCmd,
		reviewMRReopenCmd, reviewMRCancelCmd,
		commitReviewListCmd, commitReviewShowCmd, commitReviewCreateCmd, commitReviewUpdateCmd,
		commitReviewApproveCmd, commitReviewRejectCmd, commitReviewReopenCmd,
	)
	rootCmd.AddCommand(reviewCmd)
}
