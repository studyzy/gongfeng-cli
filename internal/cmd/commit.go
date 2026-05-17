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
	cmFlagPath    string
	cmFlagSince   string
	cmFlagUntil   string
)

// commit diff flags
var (
	cmDiffFlagPath             string
	cmDiffFlagIgnoreWhiteSpace bool
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
	cmRefsFlagType    string
)

// commit create-comment flags
var (
	cmCreateCommentFlagNote     string
	cmCreateCommentFlagPath     string
	cmCreateCommentFlagLine     int
	cmCreateCommentFlagLineType string
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
		if cmFlagPath != "" {
			opts.Path = gongfeng.Ptr(cmFlagPath)
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
		var opts *gongfeng.GetCommitDiffOptions
		if cmDiffFlagPath != "" || cmDiffFlagIgnoreWhiteSpace {
			opts = &gongfeng.GetCommitDiffOptions{}
			if cmDiffFlagPath != "" {
				opts.Path = gongfeng.Ptr(cmDiffFlagPath)
			}
			if cmDiffFlagIgnoreWhiteSpace {
				opts.IgnoreWhiteSpace = gongfeng.Ptr(true)
			}
		}
		diffs, _, err := apiClient.Commits.GetCommitDiff(ctx, projectID(), args[0], opts)
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

var commitCreateCommentCmd = &cobra.Command{
	Use:   "create-comment <sha>",
	Short: "创建提交评论",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if cmCreateCommentFlagNote == "" {
			exitWithParamError("--note is required", "Specify the comment content")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateCommitCommentOptions{
			Note: gongfeng.Ptr(cmCreateCommentFlagNote),
		}
		if cmCreateCommentFlagPath != "" {
			opts.Path = gongfeng.Ptr(cmCreateCommentFlagPath)
		}
		if cmCreateCommentFlagLine != 0 {
			opts.Line = gongfeng.Ptr(cmCreateCommentFlagLine)
		}
		if cmCreateCommentFlagLineType != "" {
			opts.LineType = gongfeng.Ptr(cmCreateCommentFlagLineType)
		}

		comment, _, err := apiClient.Commits.CreateCommitComment(ctx, projectID(), args[0], opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, comment, !flagPretty)
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
		if cmRefsFlagType != "" {
			opts.Type = gongfeng.Ptr(cmRefsFlagType)
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
	commitListCmd.Flags().StringVar(&cmFlagPath, "path", "", "文件路径过滤")
	commitListCmd.Flags().StringVar(&cmFlagSince, "since", "", "起始时间（ISO 8601 格式）")
	commitListCmd.Flags().StringVar(&cmFlagUntil, "until", "", "截止时间（ISO 8601 格式）")

	// commit diff flags
	commitDiffCmd.Flags().StringVar(&cmDiffFlagPath, "path", "", "按文件路径过滤")
	commitDiffCmd.Flags().BoolVar(&cmDiffFlagIgnoreWhiteSpace, "ignore-white-space", false, "忽略空白字符差异")

	// commit comments flags
	commitCommentsCmd.Flags().IntVar(&cmCommentsFlagPage, "page", 0, "页码")
	commitCommentsCmd.Flags().IntVar(&cmCommentsFlagPerPage, "per-page", 0, "每页条数")

	// commit create-comment flags
	commitCreateCommentCmd.Flags().StringVar(&cmCreateCommentFlagNote, "note", "", "评论内容（必需）")
	commitCreateCommentCmd.Flags().StringVar(&cmCreateCommentFlagPath, "path", "", "文件路径")
	commitCreateCommentCmd.Flags().IntVar(&cmCreateCommentFlagLine, "line", 0, "行号")
	commitCreateCommentCmd.Flags().StringVar(&cmCreateCommentFlagLineType, "line-type", "", "行类型（new/old）")

	// commit refs flags
	commitRefsCmd.Flags().IntVar(&cmRefsFlagPage, "page", 0, "页码")
	commitRefsCmd.Flags().IntVar(&cmRefsFlagPerPage, "per-page", 0, "每页条数")
	commitRefsCmd.Flags().StringVar(&cmRefsFlagType, "type", "", "引用类型（branch/tag）")

	commitCmd.AddCommand(commitListCmd, commitShowCmd, commitDiffCmd, commitCommentsCmd, commitCreateCommentCmd, commitRefsCmd)
	rootCmd.AddCommand(commitCmd)
}
