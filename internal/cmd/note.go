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
	ntFlagBody    string
	ntFlagPage    int
	ntFlagPerPage int
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

func init() {
	// mr-list flags
	mrNoteListCmd.Flags().IntVar(&ntFlagPage, "page", 0, "页码")
	mrNoteListCmd.Flags().IntVar(&ntFlagPerPage, "per-page", 0, "每页数量")

	// mr-create flags
	mrNoteCreateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	// mr-update flags
	mrNoteUpdateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	// issue-list flags
	issueNoteListCmd.Flags().IntVar(&ntFlagPage, "page", 0, "页码")
	issueNoteListCmd.Flags().IntVar(&ntFlagPerPage, "per-page", 0, "每页数量")

	// issue-create flags
	issueNoteCreateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	// issue-update flags
	issueNoteUpdateCmd.Flags().StringVar(&ntFlagBody, "body", "", "评论内容（必需）")

	noteCmd.AddCommand(
		mrNoteListCmd, mrNoteShowCmd, mrNoteCreateCmd, mrNoteUpdateCmd,
		issueNoteListCmd, issueNoteShowCmd, issueNoteCreateCmd, issueNoteUpdateCmd,
	)
	rootCmd.AddCommand(noteCmd)
}
