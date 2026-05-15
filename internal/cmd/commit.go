package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// commit flags
var (
	cmFlagPage    int
	cmFlagPerPage int
	cmFlagRefName string
	cmFlagSince   string
	cmFlagUntil   string
)

// commit comments flags
var (
	cmCommentsFlagPage    int
	cmCommentsFlagPerPage int
)

// commit refs flags
var (
	cmRefsFlagPage    int
	cmRefsFlagPerPage int
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "提交管理",
}

var commitListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取提交列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListCommitsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    cmFlagPage,
				PerPage: cmFlagPerPage,
			},
		}
		if cmFlagRefName != "" {
			opts.RefName = gongfeng.Ptr(cmFlagRefName)
		}
		if cmFlagSince != "" {
			opts.Since = gongfeng.Ptr(cmFlagSince)
		}
		if cmFlagUntil != "" {
			opts.Until = gongfeng.Ptr(cmFlagUntil)
		}

		commits, _, err := apiClient.Commits.ListCommits(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, commits, !flagPretty)
	},
}

var commitShowCmd = &cobra.Command{
	Use:   "show <sha>",
	Short: "获取提交详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		commit, _, err := apiClient.Commits.GetCommit(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(commit, "message")
	},
}

var commitDiffCmd = &cobra.Command{
	Use:   "diff <sha>",
	Short: "获取提交 Diff",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		diffs, _, err := apiClient.Commits.GetCommitDiff(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, diffs, !flagPretty)
	},
}

var commitCommentsCmd = &cobra.Command{
	Use:   "comments <sha>",
	Short: "获取提交评论",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListCommitCommentsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    cmCommentsFlagPage,
				PerPage: cmCommentsFlagPerPage,
			},
		}

		comments, _, err := apiClient.Commits.ListCommitComments(ctx, projectID(), args[0], opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, comments, !flagPretty)
	},
}

var commitRefsCmd = &cobra.Command{
	Use:   "refs <sha>",
	Short: "获取提交关联的分支和 Tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListCommitRefsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    cmRefsFlagPage,
				PerPage: cmRefsFlagPerPage,
			},
		}

		refs, _, err := apiClient.Commits.ListCommitRefs(ctx, projectID(), args[0], opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, refs, !flagPretty)
	},
}

func init() {
	// commit list flags
	commitListCmd.Flags().IntVar(&cmFlagPage, "page", 0, "页码")
	commitListCmd.Flags().IntVar(&cmFlagPerPage, "per-page", 0, "每页条数")
	commitListCmd.Flags().StringVar(&cmFlagRefName, "ref-name", "", "分支名或 Tag 名")
	commitListCmd.Flags().StringVar(&cmFlagSince, "since", "", "起始时间（ISO 8601 格式）")
	commitListCmd.Flags().StringVar(&cmFlagUntil, "until", "", "截止时间（ISO 8601 格式）")

	// commit comments flags
	commitCommentsCmd.Flags().IntVar(&cmCommentsFlagPage, "page", 0, "页码")
	commitCommentsCmd.Flags().IntVar(&cmCommentsFlagPerPage, "per-page", 0, "每页条数")

	// commit refs flags
	commitRefsCmd.Flags().IntVar(&cmRefsFlagPage, "page", 0, "页码")
	commitRefsCmd.Flags().IntVar(&cmRefsFlagPerPage, "per-page", 0, "每页条数")

	commitCmd.AddCommand(commitListCmd, commitShowCmd, commitDiffCmd, commitCommentsCmd, commitRefsCmd)
	rootCmd.AddCommand(commitCmd)
}
