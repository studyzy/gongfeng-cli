package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// note 子命令独有的 flag 变量
var (
	ntFlagBody          string
	ntFlagPage          int
	ntFlagPerPage       int
	ntFlagPath          string
	ntFlagLine          string
	ntFlagLineType      string
	ntFlagReviewerState string
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "评论管理",
}

// MR 评论
var mrNoteListCmd = &cobra.Command{
	Use:   "mr-list <mr_id>",
	Short: "获取 MR 评论列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		opts := &gongfeng.ListMergeRequestNotesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    ntFlagPage,
				PerPage: ntFlagPerPage,
			},
		}
		notes, _, err := apiClient.Notes.ListMergeRequestNotes(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, notes, !flagPretty)
	},
}

var mrNoteShowCmd = &cobra.Command{
	Use:   "mr-show <mr_id> <note_id>",
	Short: "获取指定 MR 评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		noteID := atoi(args[1], "note_id")
		note, _, err := apiClient.Notes.GetMergeRequestNote(context.Background(), projectID(), mrID, noteID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var mrNoteCreateCmd = &cobra.Command{
	Use:   "mr-create <mr_id>",
	Short: "创建 MR 评论",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.CreateMergeRequestNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		if ntFlagPath != "" {
			opts.Path = gongfeng.Ptr(ntFlagPath)
		}
		if ntFlagLine != "" {
			opts.Line = gongfeng.Ptr(ntFlagLine)
		}
		if ntFlagLineType != "" {
			opts.LineType = gongfeng.Ptr(ntFlagLineType)
		}
		if ntFlagReviewerState != "" {
			opts.ReviewerState = gongfeng.Ptr(ntFlagReviewerState)
		}
		note, _, err := apiClient.Notes.CreateMergeRequestNote(context.Background(), projectID(), mrID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var mrNoteUpdateCmd = &cobra.Command{
	Use:   "mr-update <mr_id> <note_id>",
	Short: "更新 MR 评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		mrID := atoi(args[0], "mr_id")
		noteID := atoi(args[1], "note_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.UpdateMergeRequestNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		if ntFlagReviewerState != "" {
			opts.ReviewerState = gongfeng.Ptr(ntFlagReviewerState)
		}
		note, _, err := apiClient.Notes.UpdateMergeRequestNote(context.Background(), projectID(), mrID, noteID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

// Issue 评论
var issueNoteListCmd = &cobra.Command{
	Use:   "issue-list <issue_id>",
	Short: "获取 Issue 评论列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		opts := &gongfeng.ListIssueNotesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    ntFlagPage,
				PerPage: ntFlagPerPage,
			},
		}
		notes, _, err := apiClient.Notes.ListIssueNotes(context.Background(), projectID(), issueID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, notes, !flagPretty)
	},
}

var issueNoteShowCmd = &cobra.Command{
	Use:   "issue-show <issue_id> <note_id>",
	Short: "获取指定 Issue 评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		noteID := atoi(args[1], "note_id")
		note, _, err := apiClient.Notes.GetIssueNote(context.Background(), projectID(), issueID, noteID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var issueNoteCreateCmd = &cobra.Command{
	Use:   "issue-create <issue_id>",
	Short: "创建 Issue 评论",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.CreateIssueNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		note, _, err := apiClient.Notes.CreateIssueNote(context.Background(), projectID(), issueID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var issueNoteUpdateCmd = &cobra.Command{
	Use:   "issue-update <issue_id> <note_id>",
	Short: "更新 Issue 评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		noteID := atoi(args[1], "note_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.UpdateIssueNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		note, _, err := apiClient.Notes.UpdateIssueNote(context.Background(), projectID(), issueID, noteID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

// Review 评论
var reviewNoteListCmd = &cobra.Command{
	Use:   "review-list <review_id>",
	Short: "获取代码评审评论列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		opts := &gongfeng.ListReviewNotesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    ntFlagPage,
				PerPage: ntFlagPerPage,
			},
		}
		notes, _, err := apiClient.Notes.ListReviewNotes(context.Background(), projectID(), reviewID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, notes, !flagPretty)
	},
}

var reviewNoteShowCmd = &cobra.Command{
	Use:   "review-show <review_id> <note_id>",
	Short: "获取指定代码评审评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		noteID := atoi(args[1], "note_id")
		note, _, err := apiClient.Notes.GetReviewNote(context.Background(), projectID(), reviewID, noteID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var reviewNoteCreateCmd = &cobra.Command{
	Use:   "review-create <review_id>",
	Short: "创建代码评审评论",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.CreateReviewNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		if ntFlagPath != "" {
			opts.Path = gongfeng.Ptr(ntFlagPath)
		}
		if ntFlagLine != "" {
			opts.Line = gongfeng.Ptr(ntFlagLine)
		}
		if ntFlagLineType != "" {
			opts.LineType = gongfeng.Ptr(ntFlagLineType)
		}
		if ntFlagReviewerState != "" {
			opts.ReviewerState = gongfeng.Ptr(ntFlagReviewerState)
		}
		note, _, err := apiClient.Notes.CreateReviewNote(context.Background(), projectID(), reviewID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

var reviewNoteUpdateCmd = &cobra.Command{
	Use:   "review-update <review_id> <note_id>",
	Short: "更新代码评审评论",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		reviewID := atoi(args[0], "review_id")
		noteID := atoi(args[1], "note_id")
		if ntFlagBody == "" {
			exitWithParamError("--body is required", "Specify the note body")
		}
		opts := &gongfeng.UpdateReviewNoteOptions{
			Body: gongfeng.Ptr(ntFlagBody),
		}
		if ntFlagReviewerState != "" {
			opts.ReviewerState = gongfeng.Ptr(ntFlagReviewerState)
		}
		note, _, err := apiClient.Notes.UpdateReviewNote(context.Background(), projectID(), reviewID, noteID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(note, "body")
	},
}

func init() {
	// mr-list flags
	mrNoteListCmd.Flags().IntVar(&ntFlagPage, "page", 0, "页码")
	mrNoteListCmd.Flags().IntVar(&ntFlagPerPage, "per-page", 0, "每页数量")

	// mr-create flags
	mrNoteCreateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")
	mrNoteCreateCmd.Flags().StringVar(&ntFlagPath, "path", "", "文件路径")
	mrNoteCreateCmd.Flags().StringVar(&ntFlagLine, "line", "", "行号")
	mrNoteCreateCmd.Flags().StringVar(&ntFlagLineType, "line-type", "", "行类型（new/old）")
	mrNoteCreateCmd.Flags().StringVar(&ntFlagReviewerState, "reviewer-state", "", "评审状态")

	// mr-update flags
	mrNoteUpdateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")
	mrNoteUpdateCmd.Flags().StringVar(&ntFlagReviewerState, "reviewer-state", "", "评审状态")

	// issue-list flags
	issueNoteListCmd.Flags().IntVar(&ntFlagPage, "page", 0, "页码")
	issueNoteListCmd.Flags().IntVar(&ntFlagPerPage, "per-page", 0, "每页数量")

	// issue-create flags
	issueNoteCreateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	// issue-update flags
	issueNoteUpdateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	// review-list flags
	reviewNoteListCmd.Flags().IntVar(&ntFlagPage, "page", 0, "页码")
	reviewNoteListCmd.Flags().IntVar(&ntFlagPerPage, "per-page", 0, "每页数量")

	// review-create flags
	reviewNoteCreateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")
	reviewNoteCreateCmd.Flags().StringVar(&ntFlagPath, "path", "", "文件路径")
	reviewNoteCreateCmd.Flags().StringVar(&ntFlagLine, "line", "", "行号")
	reviewNoteCreateCmd.Flags().StringVar(&ntFlagLineType, "line-type", "", "行类型（new/old）")
	reviewNoteCreateCmd.Flags().StringVar(&ntFlagReviewerState, "reviewer-state", "", "评审状态")

	// review-update flags
	reviewNoteUpdateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")
	reviewNoteUpdateCmd.Flags().StringVar(&ntFlagReviewerState, "reviewer-state", "", "评审状态")

	noteCmd.AddCommand(
		mrNoteListCmd, mrNoteShowCmd, mrNoteCreateCmd, mrNoteUpdateCmd,
		issueNoteListCmd, issueNoteShowCmd, issueNoteCreateCmd, issueNoteUpdateCmd,
		reviewNoteListCmd, reviewNoteShowCmd, reviewNoteCreateCmd, reviewNoteUpdateCmd,
	)
	rootCmd.AddCommand(noteCmd)
}
