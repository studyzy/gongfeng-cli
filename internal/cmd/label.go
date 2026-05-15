package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// label 子命令独有的 flag 变量
var (
	lbFlagName        string
	lbFlagColor       string
	lbFlagDescription string
	lbFlagNewName     string
	lbFlagPage        int
	lbFlagPerPage     int
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "标签管理",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取标签列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListLabelsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    lbFlagPage,
				PerPage: lbFlagPerPage,
			},
		}
		labels, _, err := apiClient.Labels.ListLabels(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, labels, !flagPretty)
	},
}

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if lbFlagName == "" {
			exitWithParamError("--name is required", "Specify the label name")
		}
		if lbFlagColor == "" {
			exitWithParamError("--color is required", "Specify the label color, e.g. #428BCA")
		}
		opts := &gongfeng.CreateLabelOptions{
			Name:  gongfeng.Ptr(lbFlagName),
			Color: gongfeng.Ptr(lbFlagColor),
		}
		if lbFlagDescription != "" {
			opts.Description = gongfeng.Ptr(lbFlagDescription)
		}
		label, _, err := apiClient.Labels.CreateLabel(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, label, !flagPretty)
	},
}

var labelUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if lbFlagName == "" {
			exitWithParamError("--name is required", "Specify the current label name to update")
		}
		opts := &gongfeng.UpdateLabelOptions{
			Name: gongfeng.Ptr(lbFlagName),
		}
		if lbFlagNewName != "" {
			opts.NewName = gongfeng.Ptr(lbFlagNewName)
		}
		if lbFlagColor != "" {
			opts.Color = gongfeng.Ptr(lbFlagColor)
		}
		if lbFlagDescription != "" {
			opts.Description = gongfeng.Ptr(lbFlagDescription)
		}
		label, _, err := apiClient.Labels.UpdateLabel(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, label, !flagPretty)
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if lbFlagName == "" {
			exitWithParamError("--name is required", "Specify the label name to delete")
		}
		opts := &gongfeng.DeleteLabelOptions{
			Name: gongfeng.Ptr(lbFlagName),
		}
		_, err := apiClient.Labels.DeleteLabel(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// list flags
	labelListCmd.Flags().IntVar(&lbFlagPage, "page", 0, "页码")
	labelListCmd.Flags().IntVar(&lbFlagPerPage, "per-page", 0, "每页数量")

	// create flags
	labelCreateCmd.Flags().StringVar(&lbFlagName, "name", "", "标签名称（必需）")
	labelCreateCmd.Flags().StringVar(&lbFlagColor, "color", "", "标签颜色，如 #428BCA（必需）")
	labelCreateCmd.Flags().StringVar(&lbFlagDescription, "description", "", "标签描述")

	// update flags
	labelUpdateCmd.Flags().StringVar(&lbFlagName, "name", "", "当前标签名称（必需）")
	labelUpdateCmd.Flags().StringVar(&lbFlagNewName, "new-name", "", "新标签名称")
	labelUpdateCmd.Flags().StringVar(&lbFlagColor, "color", "", "标签颜色")
	labelUpdateCmd.Flags().StringVar(&lbFlagDescription, "description", "", "标签描述")

	// delete flags
	labelDeleteCmd.Flags().StringVar(&lbFlagName, "name", "", "标签名称（必需）")

	labelCmd.AddCommand(labelListCmd, labelCreateCmd, labelUpdateCmd, labelDeleteCmd)
	rootCmd.AddCommand(labelCmd)
}
