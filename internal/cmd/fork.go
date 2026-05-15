package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Fork 管理",
}

var forkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Fork 项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()

		project, _, err := apiClient.Forks.ForkProject(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, project, !flagPretty)
	},
}

var forkLinkCmd = &cobra.Command{
	Use:   "link <forked_from_id>",
	Short: "创建 Fork 关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		forkedFromID := atoi(args[0], "forked_from_id")
		ctx := context.Background()

		_, err := apiClient.Forks.CreateForkRelation(ctx, projectID(), forkedFromID)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var forkUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "取消 Fork 关系",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()

		_, err := apiClient.Forks.DeleteForkRelation(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	forkCmd.AddCommand(forkCreateCmd, forkLinkCmd, forkUnlinkCmd)
	rootCmd.AddCommand(forkCmd)
}
