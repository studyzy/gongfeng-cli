package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// release 子命令独有的 flag 变量
var (
	rlFlagTagName     string
	rlFlagDescription string
	rlFlagPage        int
	rlFlagPerPage     int
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "版本发布管理",
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取 Release 列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListReleasesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    rlFlagPage,
				PerPage: rlFlagPerPage,
			},
		}
		releases, _, err := apiClient.Releases.ListReleases(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, releases, !flagPretty)
	},
}

var releaseShowCmd = &cobra.Command{
	Use:   "show <tag_name>",
	Short: "获取指定 Release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		tagName := args[0]
		release, _, err := apiClient.Releases.GetRelease(context.Background(), projectID(), tagName)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(release, "description")
	},
}

var releaseCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建 Release",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if rlFlagTagName == "" {
			exitWithParamError("--tag-name is required", "Specify the tag name for the release")
		}
		opts := &gongfeng.CreateReleaseOptions{
			TagName: gongfeng.Ptr(rlFlagTagName),
		}
		if rlFlagDescription != "" {
			opts.Description = gongfeng.Ptr(rlFlagDescription)
		}
		release, _, err := apiClient.Releases.CreateRelease(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, release, !flagPretty)
	},
}

var releaseUpdateCmd = &cobra.Command{
	Use:   "update <tag_name>",
	Short: "更新 Release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		tagName := args[0]
		opts := &gongfeng.UpdateReleaseOptions{}
		if rlFlagDescription != "" {
			opts.Description = gongfeng.Ptr(rlFlagDescription)
		}
		release, _, err := apiClient.Releases.UpdateRelease(context.Background(), projectID(), tagName, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, release, !flagPretty)
	},
}

var releaseDeleteCmd = &cobra.Command{
	Use:   "delete <tag_name>",
	Short: "删除 Release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		tagName := args[0]
		_, err := apiClient.Releases.DeleteRelease(context.Background(), projectID(), tagName)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// list flags
	releaseListCmd.Flags().IntVar(&rlFlagPage, "page", 0, "页码")
	releaseListCmd.Flags().IntVar(&rlFlagPerPage, "per-page", 0, "每页数量")

	// create flags
	releaseCreateCmd.Flags().StringVar(&rlFlagTagName, "tag-name", "", "Tag 名称（必需）")
	releaseCreateCmd.Flags().StringVar(&rlFlagDescription, "description", "", "Release 描述")

	// update flags
	releaseUpdateCmd.Flags().StringVar(&rlFlagDescription, "description", "", "Release 描述")

	releaseCmd.AddCommand(releaseListCmd, releaseShowCmd, releaseCreateCmd, releaseUpdateCmd, releaseDeleteCmd)
	rootCmd.AddCommand(releaseCmd)
}
