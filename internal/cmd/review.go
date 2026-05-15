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
	rvFlagCommit              string
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

var reviewApproveCmd = &cobra.Command{
	Use:   "approve <mr_id>",
	Short: "通过 MR 评审",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.ApproveMROptions{}
		if rvFlagMessage != "" {
			opts.Message = gongfeng.Ptr(rvFlagMessage)
		}
		_, err := apiClient.Reviews.ApproveMR(context.Background(), projectID(), mrID, opts)
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
		opts := &gongfeng.RejectMROptions{}
		if rvFlagMessage != "" {
			opts.Message = gongfeng.Ptr(rvFlagMessage)
		}
		_, err := apiClient.Reviews.RejectMR(context.Background(), projectID(), mrID, opts)
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
		if rvFlagCommit == "" {
			exitWithParamError("--commit is required", "Specify the commit SHA for the review")
		}
		if rvFlagTitle == "" {
			exitWithParamError("--title is required", "Specify the title for the commit review")
		}
		opts := &gongfeng.CreateCommitReviewOptions{
			Commit: gongfeng.Ptr(rvFlagCommit),
			Title:  gongfeng.Ptr(rvFlagTitle),
		}
		if rvFlagDescription != "" {
			opts.Description = gongfeng.Ptr(rvFlagDescription)
		}
		review, _, err := apiClient.Reviews.CreateCommitReview(context.Background(), projectID(), opts)
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
		opts := &gongfeng.ApproveCommitReviewOptions{}
		if rvFlagMessage != "" {
			opts.Message = gongfeng.Ptr(rvFlagMessage)
		}
		_, err := apiClient.Reviews.ApproveCommitReview(context.Background(), projectID(), reviewID, opts)
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
		opts := &gongfeng.RejectCommitReviewOptions{}
		if rvFlagMessage != "" {
			opts.Message = gongfeng.Ptr(rvFlagMessage)
		}
		_, err := apiClient.Reviews.RejectCommitReview(context.Background(), projectID(), reviewID, opts)
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

	// commit-create flags
	commitReviewCreateCmd.Flags().StringVar(&rvFlagCommit, "commit", "", "提交 SHA（必需）")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagTitle, "title", "", "评审标题（必需）")
	commitReviewCreateCmd.Flags().StringVar(&rvFlagDescription, "description", "", "评审描述")

	// commit-approve flags
	commitReviewApproveCmd.Flags().StringVar(&rvFlagMessage, "message", "", "审批意见")

	// commit-reject flags
	commitReviewRejectCmd.Flags().StringVar(&rvFlagMessage, "message", "", "拒绝原因")

	reviewCmd.AddCommand(
		reviewInviteCmd, reviewRemoveCmd, reviewApproveCmd, reviewRejectCmd,
		commitReviewListCmd, commitReviewShowCmd, commitReviewCreateCmd,
		commitReviewApproveCmd, commitReviewRejectCmd,
	)
	rootCmd.AddCommand(reviewCmd)
}
