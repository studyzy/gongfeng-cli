package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// branch flags
var (
	brFlagPage       int
	brFlagPerPage    int
	brFlagBranchName string
	brFlagRef        string
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "分支管理",
}

var branchListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取分支列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListBranchesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    brFlagPage,
				PerPage: brFlagPerPage,
			},
		}

		branches, _, err := apiClient.Branches.ListBranches(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, branches, !flagPretty)
	},
}

var branchShowCmd = &cobra.Command{
	Use:   "show <branch_name>",
	Short: "获取分支详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		branch, _, err := apiClient.Branches.GetBranch(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(branch, "")
	},
}

var branchCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建分支",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		if brFlagBranchName == "" {
			exitWithParamError("--branch-name is required", "Provide a branch name with --branch-name")
		}
		if brFlagRef == "" {
			exitWithParamError("--ref is required", "Provide a source ref with --ref")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateBranchOptions{
			BranchName: gongfeng.Ptr(brFlagBranchName),
			Ref:        gongfeng.Ptr(brFlagRef),
		}

		branch, _, err := apiClient.Branches.CreateBranch(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, branch, !flagPretty)
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete <branch_name>",
	Short: "删除分支",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		_, err := apiClient.Branches.DeleteBranch(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// branch list flags
	branchListCmd.Flags().IntVar(&brFlagPage, "page", 0, "页码")
	branchListCmd.Flags().IntVar(&brFlagPerPage, "per-page", 0, "每页条数")

	// branch create flags
	branchCreateCmd.Flags().StringVar(&brFlagBranchName, "branch-name", "", "分支名称（必需）")
	branchCreateCmd.Flags().StringVar(&brFlagRef, "ref", "", "源分支或 SHA（必需）")

	branchCmd.AddCommand(branchListCmd, branchShowCmd, branchCreateCmd, branchDeleteCmd)
	rootCmd.AddCommand(branchCmd)
}
