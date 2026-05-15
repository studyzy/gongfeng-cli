package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// repo 子命令独有的 flag 变量
var (
	rpFlagPath          string
	rpFlagRef           string
	rpFlagFilePath      string
	rpFlagBranch        string
	rpFlagContent       string
	rpFlagCommitMessage string
	rpFlagFrom          string
	rpFlagTo            string
	rpFlagPage          int
	rpFlagPerPage       int
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "仓库文件管理",
}

var repoTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "获取文件树",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListTreeOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    rpFlagPage,
				PerPage: rpFlagPerPage,
			},
		}
		if rpFlagPath != "" {
			opts.Path = gongfeng.Ptr(rpFlagPath)
		}
		if rpFlagRef != "" {
			opts.RefName = gongfeng.Ptr(rpFlagRef)
		}
		tree, _, err := apiClient.Repositories.ListTree(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, tree, !flagPretty)
	},
}

var repoFileCmd = &cobra.Command{
	Use:   "file",
	Short: "获取文件内容",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rpFlagFilePath == "" {
			exitWithParamError("--file-path is required", "Specify the file path to retrieve")
		}
		opts := &gongfeng.GetFileOptions{
			FilePath: gongfeng.Ptr(rpFlagFilePath),
		}
		if rpFlagRef != "" {
			opts.Ref = gongfeng.Ptr(rpFlagRef)
		}
		file, _, err := apiClient.Repositories.GetFile(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(file, "content")
	},
}

var repoCreateFileCmd = &cobra.Command{
	Use:   "create-file",
	Short: "创建文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rpFlagFilePath == "" {
			exitWithParamError("--file-path is required", "Specify the file path to create")
		}
		if rpFlagBranch == "" {
			exitWithParamError("--branch-name is required", "Specify the branch name")
		}
		if rpFlagContent == "" {
			exitWithParamError("--content is required", "Specify the file content")
		}
		if rpFlagCommitMessage == "" {
			exitWithParamError("--commit-message is required", "Specify the commit message")
		}
		opts := &gongfeng.CreateFileOptions{
			FilePath:      gongfeng.Ptr(rpFlagFilePath),
			BranchName:    gongfeng.Ptr(rpFlagBranch),
			Content:       gongfeng.Ptr(rpFlagContent),
			CommitMessage: gongfeng.Ptr(rpFlagCommitMessage),
		}
		file, _, err := apiClient.Repositories.CreateFile(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(file, "content")
	},
}

var repoUpdateFileCmd = &cobra.Command{
	Use:   "update-file",
	Short: "更新文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rpFlagFilePath == "" {
			exitWithParamError("--file-path is required", "Specify the file path to update")
		}
		if rpFlagBranch == "" {
			exitWithParamError("--branch-name is required", "Specify the branch name")
		}
		if rpFlagContent == "" {
			exitWithParamError("--content is required", "Specify the file content")
		}
		if rpFlagCommitMessage == "" {
			exitWithParamError("--commit-message is required", "Specify the commit message")
		}
		opts := &gongfeng.UpdateFileOptions{
			FilePath:      gongfeng.Ptr(rpFlagFilePath),
			BranchName:    gongfeng.Ptr(rpFlagBranch),
			Content:       gongfeng.Ptr(rpFlagContent),
			CommitMessage: gongfeng.Ptr(rpFlagCommitMessage),
		}
		file, _, err := apiClient.Repositories.UpdateFile(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(file, "content")
	},
}

var repoDeleteFileCmd = &cobra.Command{
	Use:   "delete-file",
	Short: "删除文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rpFlagFilePath == "" {
			exitWithParamError("--file-path is required", "Specify the file path to delete")
		}
		if rpFlagBranch == "" {
			exitWithParamError("--branch-name is required", "Specify the branch name")
		}
		if rpFlagCommitMessage == "" {
			exitWithParamError("--commit-message is required", "Specify the commit message")
		}
		opts := &gongfeng.DeleteFileOptions{
			FilePath:      gongfeng.Ptr(rpFlagFilePath),
			BranchName:    gongfeng.Ptr(rpFlagBranch),
			CommitMessage: gongfeng.Ptr(rpFlagCommitMessage),
		}
		_, err := apiClient.Repositories.DeleteFile(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var repoCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "比较两个分支/Tag/SHA",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rpFlagFrom == "" {
			exitWithParamError("--from is required", "Specify the source branch/tag/SHA")
		}
		if rpFlagTo == "" {
			exitWithParamError("--to is required", "Specify the target branch/tag/SHA")
		}
		opts := &gongfeng.CompareOptions{
			From: gongfeng.Ptr(rpFlagFrom),
			To:   gongfeng.Ptr(rpFlagTo),
		}
		result, _, err := apiClient.Repositories.Compare(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, result, !flagPretty)
	},
}

func init() {
	// tree flags
	repoTreeCmd.Flags().StringVar(&rpFlagPath, "path", "", "子目录路径")
	repoTreeCmd.Flags().StringVar(&rpFlagRef, "ref-name", "", "分支/Tag/SHA")
	repoTreeCmd.Flags().IntVar(&rpFlagPage, "page", 0, "页码")
	repoTreeCmd.Flags().IntVar(&rpFlagPerPage, "per-page", 0, "每页数量")

	// file flags
	repoFileCmd.Flags().StringVar(&rpFlagFilePath, "file-path", "", "文件路径（必需）")
	repoFileCmd.Flags().StringVar(&rpFlagRef, "ref", "", "分支/Tag/SHA")

	// create-file flags
	repoCreateFileCmd.Flags().StringVar(&rpFlagFilePath, "file-path", "", "文件路径（必需）")
	repoCreateFileCmd.Flags().StringVar(&rpFlagBranch, "branch-name", "", "分支名（必需）")
	repoCreateFileCmd.Flags().StringVar(&rpFlagContent, "content", "", "文件内容（必需）")
	repoCreateFileCmd.Flags().StringVar(&rpFlagCommitMessage, "commit-message", "", "提交信息（必需）")

	// update-file flags
	repoUpdateFileCmd.Flags().StringVar(&rpFlagFilePath, "file-path", "", "文件路径（必需）")
	repoUpdateFileCmd.Flags().StringVar(&rpFlagBranch, "branch-name", "", "分支名（必需）")
	repoUpdateFileCmd.Flags().StringVar(&rpFlagContent, "content", "", "文件内容（必需）")
	repoUpdateFileCmd.Flags().StringVar(&rpFlagCommitMessage, "commit-message", "", "提交信息（必需）")

	// delete-file flags
	repoDeleteFileCmd.Flags().StringVar(&rpFlagFilePath, "file-path", "", "文件路径（必需）")
	repoDeleteFileCmd.Flags().StringVar(&rpFlagBranch, "branch-name", "", "分支名（必需）")
	repoDeleteFileCmd.Flags().StringVar(&rpFlagCommitMessage, "commit-message", "", "提交信息（必需）")

	// compare flags
	repoCompareCmd.Flags().StringVar(&rpFlagFrom, "from", "", "源分支/Tag/SHA（必需）")
	repoCompareCmd.Flags().StringVar(&rpFlagTo, "to", "", "目标分支/Tag/SHA（必需）")

	repoCmd.AddCommand(repoTreeCmd, repoFileCmd, repoCreateFileCmd, repoUpdateFileCmd, repoDeleteFileCmd, repoCompareCmd)
	rootCmd.AddCommand(repoCmd)
}
