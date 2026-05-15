package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// group list flags
var (
	grpListSearch  string
	grpListPage    int
	grpListPerPage int
)

// group create flags
var (
	grpCreateName        string
	grpCreatePath        string
	grpCreateDescription string
)

// group update flags
var (
	grpUpdateName        string
	grpUpdateDescription string
)

// group members flags
var (
	grpMembersPage    int
	grpMembersPerPage int
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "项目组管理",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取项目组列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListGroupsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    grpListPage,
				PerPage: grpListPerPage,
			},
		}
		if grpListSearch != "" {
			opts.Search = gongfeng.Ptr(grpListSearch)
		}

		groups, _, err := apiClient.Groups.ListGroups(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, groups, !flagPretty)
	},
}

var groupShowCmd = &cobra.Command{
	Use:   "show <group_id>",
	Short: "获取项目组详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		gid := atoi(args[0], "group_id")

		group, _, err := apiClient.Groups.GetGroup(ctx, gid)
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(group, "description")
	},
}

var groupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建项目组",
	RunE: func(cmd *cobra.Command, args []string) error {
		if grpCreateName == "" {
			exitWithParamError("--name is required", "Provide a group name with --name")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateGroupOptions{
			Name: gongfeng.Ptr(grpCreateName),
		}
		if grpCreatePath != "" {
			opts.Path = gongfeng.Ptr(grpCreatePath)
		}
		if grpCreateDescription != "" {
			opts.Description = gongfeng.Ptr(grpCreateDescription)
		}

		group, _, err := apiClient.Groups.CreateGroup(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return printSuccessResponse(
			fmt.Sprintf("%d", group.ID),
			group.WebURL,
			"",
		)
	},
}

var groupUpdateCmd = &cobra.Command{
	Use:   "update <group_id>",
	Short: "编辑项目组",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		gid := atoi(args[0], "group_id")

		opts := &gongfeng.EditGroupOptions{
			ID: gongfeng.Ptr(gid),
		}
		if grpUpdateName != "" {
			opts.Name = gongfeng.Ptr(grpUpdateName)
		}
		if grpUpdateDescription != "" {
			opts.Description = gongfeng.Ptr(grpUpdateDescription)
		}

		group, _, err := apiClient.Groups.EditGroup(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return printSuccessResponse(
			fmt.Sprintf("%d", group.ID),
			group.WebURL,
			"",
		)
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "delete <group_id>",
	Short: "删除项目组",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		gid := atoi(args[0], "group_id")

		_, err := apiClient.Groups.DeleteGroup(ctx, gid)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var groupMembersCmd = &cobra.Command{
	Use:   "members <group_id>",
	Short: "获取项目组成员",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		gid := atoi(args[0], "group_id")

		opts := &gongfeng.ListGroupMembersOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    grpMembersPage,
				PerPage: grpMembersPerPage,
			},
		}

		members, _, err := apiClient.Groups.ListGroupMembers(ctx, gid, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, members, !flagPretty)
	},
}

func init() {
	// group list flags
	groupListCmd.Flags().StringVar(&grpListSearch, "search", "", "按关键字搜索项目组")
	groupListCmd.Flags().IntVar(&grpListPage, "page", 0, "页码")
	groupListCmd.Flags().IntVar(&grpListPerPage, "per-page", 0, "每页条数")

	// group create flags
	groupCreateCmd.Flags().StringVar(&grpCreateName, "name", "", "项目组名称（必需）")
	groupCreateCmd.Flags().StringVar(&grpCreatePath, "path", "", "项目组路径")
	groupCreateCmd.Flags().StringVar(&grpCreateDescription, "description", "", "项目组描述")

	// group update flags
	groupUpdateCmd.Flags().StringVar(&grpUpdateName, "name", "", "项目组名称")
	groupUpdateCmd.Flags().StringVar(&grpUpdateDescription, "description", "", "项目组描述")

	// group members flags
	groupMembersCmd.Flags().IntVar(&grpMembersPage, "page", 0, "页码")
	groupMembersCmd.Flags().IntVar(&grpMembersPerPage, "per-page", 0, "每页条数")

	groupCmd.AddCommand(groupListCmd, groupShowCmd, groupCreateCmd, groupUpdateCmd, groupDeleteCmd, groupMembersCmd)
	rootCmd.AddCommand(groupCmd)
}
