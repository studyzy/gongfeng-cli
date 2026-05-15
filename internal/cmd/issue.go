package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// issue 子命令独有的 flag 变量
var (
	isFlagState       string
	isFlagLabels      string
	isFlagMilestone   string
	isFlagOrderBy     string
	isFlagSort        string
	isFlagPage        int
	isFlagPerPage     int
	isFlagTitle       string
	isFlagDescription string
	isFlagAssigneeID  int
	isFlagMilestoneID int
	isFlagStateEvent  string
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "缺陷管理",
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取缺陷列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		opts := &gongfeng.ListIssuesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    isFlagPage,
				PerPage: isFlagPerPage,
			},
		}
		if isFlagState != "" {
			opts.State = gongfeng.Ptr(isFlagState)
		}
		if isFlagLabels != "" {
			opts.Labels = gongfeng.Ptr(isFlagLabels)
		}
		if isFlagMilestone != "" {
			opts.Milestone = gongfeng.Ptr(isFlagMilestone)
		}
		if isFlagOrderBy != "" {
			opts.OrderBy = gongfeng.Ptr(isFlagOrderBy)
		}
		if isFlagSort != "" {
			opts.Sort = gongfeng.Ptr(isFlagSort)
		}
		issues, _, err := apiClient.Issues.ListIssues(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, issues, !flagPretty)
	},
}

var issueShowCmd = &cobra.Command{
	Use:   "show <issue_id>",
	Short: "获取缺陷详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		issue, _, err := apiClient.Issues.GetIssue(context.Background(), projectID(), issueID)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(issue, "description")
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建缺陷",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if isFlagTitle == "" {
			exitWithParamError("--title is required", "Specify the issue title")
		}
		opts := &gongfeng.CreateIssueOptions{
			Title: gongfeng.Ptr(isFlagTitle),
		}
		if isFlagDescription != "" {
			opts.Description = gongfeng.Ptr(isFlagDescription)
		}
		if isFlagAssigneeID != 0 {
			opts.AssigneeID = gongfeng.Ptr(isFlagAssigneeID)
		}
		if isFlagMilestoneID != 0 {
			opts.MilestoneID = gongfeng.Ptr(isFlagMilestoneID)
		}
		if isFlagLabels != "" {
			opts.Labels = gongfeng.Ptr(isFlagLabels)
		}
		issue, _, err := apiClient.Issues.CreateIssue(context.Background(), projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(issue, "description")
	},
}

var issueUpdateCmd = &cobra.Command{
	Use:   "update <issue_id>",
	Short: "更新缺陷",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		opts := &gongfeng.UpdateIssueOptions{}
		if isFlagTitle != "" {
			opts.Title = gongfeng.Ptr(isFlagTitle)
		}
		if isFlagDescription != "" {
			opts.Description = gongfeng.Ptr(isFlagDescription)
		}
		if isFlagAssigneeID != 0 {
			opts.AssigneeID = gongfeng.Ptr(isFlagAssigneeID)
		}
		if isFlagMilestoneID != 0 {
			opts.MilestoneID = gongfeng.Ptr(isFlagMilestoneID)
		}
		if isFlagLabels != "" {
			opts.Labels = gongfeng.Ptr(isFlagLabels)
		}
		if isFlagStateEvent != "" {
			opts.StateEvent = gongfeng.Ptr(isFlagStateEvent)
		}
		issue, _, err := apiClient.Issues.UpdateIssue(context.Background(), projectID(), issueID, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return printDetail(issue, "description")
	},
}

var issueDeleteCmd = &cobra.Command{
	Use:   "delete <issue_id>",
	Short: "删除缺陷",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		issueID := atoi(args[0], "issue_id")
		_, err := apiClient.Issues.DeleteIssue(context.Background(), projectID(), issueID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// list flags
	issueListCmd.Flags().StringVar(&isFlagState, "state", "", "缺陷状态过滤（opened/closed）")
	issueListCmd.Flags().StringVar(&isFlagLabels, "labels", "", "标签过滤")
	issueListCmd.Flags().StringVar(&isFlagMilestone, "milestone", "", "里程碑过滤")
	issueListCmd.Flags().StringVar(&isFlagOrderBy, "order-by", "", "排序字段")
	issueListCmd.Flags().StringVar(&isFlagSort, "sort", "", "排序方向（asc/desc）")
	issueListCmd.Flags().IntVar(&isFlagPage, "page", 0, "页码")
	issueListCmd.Flags().IntVar(&isFlagPerPage, "per-page", 0, "每页数量")

	// create flags
	issueCreateCmd.Flags().StringVar(&isFlagTitle, "title", "", "缺陷标题（必需）")
	issueCreateCmd.Flags().StringVar(&isFlagDescription, "description", "", "缺陷描述")
	issueCreateCmd.Flags().IntVar(&isFlagAssigneeID, "assignee-id", 0, "指派人 ID")
	issueCreateCmd.Flags().IntVar(&isFlagMilestoneID, "milestone-id", 0, "里程碑 ID")
	issueCreateCmd.Flags().StringVar(&isFlagLabels, "labels", "", "标签")

	// update flags
	issueUpdateCmd.Flags().StringVar(&isFlagTitle, "title", "", "缺陷标题")
	issueUpdateCmd.Flags().StringVar(&isFlagDescription, "description", "", "缺陷描述")
	issueUpdateCmd.Flags().IntVar(&isFlagAssigneeID, "assignee-id", 0, "指派人 ID")
	issueUpdateCmd.Flags().IntVar(&isFlagMilestoneID, "milestone-id", 0, "里程碑 ID")
	issueUpdateCmd.Flags().StringVar(&isFlagLabels, "labels", "", "标签")
	issueUpdateCmd.Flags().StringVar(&isFlagStateEvent, "state-event", "", "状态事件（close/reopen）")

	issueCmd.AddCommand(issueListCmd, issueShowCmd, issueCreateCmd, issueUpdateCmd, issueDeleteCmd)
	rootCmd.AddCommand(issueCmd)
}
