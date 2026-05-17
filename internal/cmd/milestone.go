package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// milestone 子命令独有的 flag 变量
var (
	msFlagPage        int
	msFlagPerPage     int
	msFlagTitle       string
	msFlagDescription string
	msFlagDueDate     string
	msFlagStateEvent  string
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "里程碑管理",
}

var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取里程碑列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListMilestonesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    msFlagPage,
				PerPage: msFlagPerPage,
			},
		}
		milestones, _, err := apiClient.Milestones.ListMilestones(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, milestones, !flagPretty)
	},
}

var milestoneShowCmd = &cobra.Command{
	Use:   "show <milestone_id>",
	Short: "获取里程碑详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		msID := atoi(args[0], "milestone_id")
		milestone, _, err := apiClient.Milestones.GetMilestone(context.Background(), projectID(), msID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(milestone, "description")
	},
}

var milestoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建里程碑",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if msFlagTitle == "" {
			exitWithParamError("--title is required", "Specify the title for the milestone")
		}
		opts := &gongfeng.CreateMilestoneOptions{
			Title: gongfeng.Ptr(msFlagTitle),
		}
		if msFlagDescription != "" {
			opts.Description = gongfeng.Ptr(msFlagDescription)
		}
		if msFlagDueDate != "" {
			opts.DueDate = gongfeng.Ptr(msFlagDueDate)
		}
		if msFlagStateEvent != "" {
			opts.StateEvent = gongfeng.Ptr(msFlagStateEvent)
		}
		milestone, _, err := apiClient.Milestones.CreateMilestone(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(milestone, "description")
	},
}

var milestoneUpdateCmd = &cobra.Command{
	Use:   "update <milestone_id>",
	Short: "更新里程碑",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		msID := atoi(args[0], "milestone_id")
		opts := &gongfeng.EditMilestoneOptions{}
		if msFlagTitle != "" {
			opts.Title = gongfeng.Ptr(msFlagTitle)
		}
		if msFlagDescription != "" {
			opts.Description = gongfeng.Ptr(msFlagDescription)
		}
		if msFlagDueDate != "" {
			opts.DueDate = gongfeng.Ptr(msFlagDueDate)
		}
		if msFlagStateEvent != "" {
			opts.StateEvent = gongfeng.Ptr(msFlagStateEvent)
		}
		milestone, _, err := apiClient.Milestones.EditMilestone(context.Background(), projectID(), msID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(milestone, "description")
	},
}

var milestoneDeleteCmd = &cobra.Command{
	Use:   "delete <milestone_id>",
	Short: "删除里程碑",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		msID := atoi(args[0], "milestone_id")
		_, err := apiClient.Milestones.DeleteMilestone(context.Background(), projectID(), msID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var milestoneIssuesCmd = &cobra.Command{
	Use:   "issues <milestone_id>",
	Short: "获取里程碑下的缺陷",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		msID := atoi(args[0], "milestone_id")
		opts := &gongfeng.ListMilestoneIssuesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    msFlagPage,
				PerPage: msFlagPerPage,
			},
		}
		issues, _, err := apiClient.Milestones.ListMilestoneIssues(context.Background(), projectID(), msID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, issues, !flagPretty)
	},
}

func init() {
	// list flags
	milestoneListCmd.Flags().IntVar(&msFlagPage, "page", 0, "页码")
	milestoneListCmd.Flags().IntVar(&msFlagPerPage, "per-page", 0, "每页数量")

	// create flags
	milestoneCreateCmd.Flags().StringVar(&msFlagTitle, "title", "", "里程碑标题（必需）")
	milestoneCreateCmd.Flags().StringVar(&msFlagDescription, "description", "", "里程碑描述")
	milestoneCreateCmd.Flags().StringVar(&msFlagDueDate, "due-date", "", "截止日期（格式：YYYY-MM-DD）")
	milestoneCreateCmd.Flags().StringVar(&msFlagStateEvent, "state-event", "", "状态事件（close/activate）")

	// update flags
	milestoneUpdateCmd.Flags().StringVar(&msFlagTitle, "title", "", "里程碑标题")
	milestoneUpdateCmd.Flags().StringVar(&msFlagDescription, "description", "", "里程碑描述")
	milestoneUpdateCmd.Flags().StringVar(&msFlagDueDate, "due-date", "", "截止日期（格式：YYYY-MM-DD）")
	milestoneUpdateCmd.Flags().StringVar(&msFlagStateEvent, "state-event", "", "状态事件（close/activate）")

	// issues flags
	milestoneIssuesCmd.Flags().IntVar(&msFlagPage, "page", 0, "页码")
	milestoneIssuesCmd.Flags().IntVar(&msFlagPerPage, "per-page", 0, "每页数量")

	milestoneCmd.AddCommand(milestoneListCmd, milestoneShowCmd, milestoneCreateCmd, milestoneUpdateCmd, milestoneDeleteCmd, milestoneIssuesCmd)
	rootCmd.AddCommand(milestoneCmd)
}
