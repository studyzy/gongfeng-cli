package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// tag flags
var (
	tgFlagPage    int
	tgFlagPerPage int
	tgFlagTagName string
	tgFlagRef     string
	tgFlagMessage string
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "标签管理",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取 Tag 列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListTagsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    tgFlagPage,
				PerPage: tgFlagPerPage,
			},
		}

		tags, _, err := apiClient.Tags.ListTags(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, tags, !flagPretty)
	},
}

var tagShowCmd = &cobra.Command{
	Use:   "show <tag_name>",
	Short: "获取指定 Tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		tag, _, err := apiClient.Tags.GetTag(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(tag, "message")
	},
}

var tagCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建 Tag",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		if tgFlagTagName == "" {
			exitWithParamError("--tag-name is required", "Provide a tag name with --tag-name")
		}
		if tgFlagRef == "" {
			exitWithParamError("--ref is required", "Provide a source ref with --ref")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateTagOptions{
			TagName: gongfeng.Ptr(tgFlagTagName),
			Ref:     gongfeng.Ptr(tgFlagRef),
		}
		if tgFlagMessage != "" {
			opts.Message = gongfeng.Ptr(tgFlagMessage)
		}

		tag, _, err := apiClient.Tags.CreateTag(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, tag, !flagPretty)
	},
}

var tagDeleteCmd = &cobra.Command{
	Use:   "delete <tag_name>",
	Short: "删除 Tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		_, err := apiClient.Tags.DeleteTag(ctx, projectID(), args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// tag list flags
	tagListCmd.Flags().IntVar(&tgFlagPage, "page", 0, "页码")
	tagListCmd.Flags().IntVar(&tgFlagPerPage, "per-page", 0, "每页条数")

	// tag create flags
	tagCreateCmd.Flags().StringVar(&tgFlagTagName, "tag-name", "", "Tag 名称（必需）")
	tagCreateCmd.Flags().StringVar(&tgFlagRef, "ref", "", "源分支或 SHA（必需）")
	tagCreateCmd.Flags().StringVar(&tgFlagMessage, "message", "", "Tag 描述信息")

	tagCmd.AddCommand(tagListCmd, tagShowCmd, tagCreateCmd, tagDeleteCmd)
	rootCmd.AddCommand(tagCmd)
}
