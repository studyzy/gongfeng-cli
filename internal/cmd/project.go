package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// project list flags
var (
	projListSearch  string
	projListOrderBy string
	projListSort    string
	projListPage    int
	projListPerPage int
)

// project create flags
var (
	projCreateName        string
	projCreatePath        string
	projCreateDescription string
	projCreateNamespaceID int
)

// project members flags
var (
	projMembersPage    int
	projMembersPerPage int
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目管理",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取项目列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListProjectsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    projListPage,
				PerPage: projListPerPage,
			},
		}
		if projListSearch != "" {
			opts.Search = gongfeng.Ptr(projListSearch)
		}
		if projListOrderBy != "" {
			opts.OrderBy = gongfeng.Ptr(projListOrderBy)
		}
		if projListSort != "" {
			opts.Sort = gongfeng.Ptr(projListSort)
		}

		projects, _, err := apiClient.Projects.ListProjects(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, projects, !flagPretty)
	},
}

var projectShowCmd = &cobra.Command{
	Use:   "show <project_id>",
	Short: "获取项目详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// 项目 ID 可以是数字或路径
		var pid interface{}
		if id, err := strconv.Atoi(args[0]); err == nil {
			pid = id
		} else {
			pid = args[0]
		}

		project, _, err := apiClient.Projects.GetProject(ctx, pid)
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(project, "description")
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		if projCreateName == "" {
			exitWithParamError("--name is required", "Provide a project name with --name")
		}

		ctx := context.Background()
		opts := &gongfeng.CreateProjectOptions{
			Name: gongfeng.Ptr(projCreateName),
		}
		if projCreatePath != "" {
			opts.Path = gongfeng.Ptr(projCreatePath)
		}
		if projCreateDescription != "" {
			opts.Description = gongfeng.Ptr(projCreateDescription)
		}
		if projCreateNamespaceID != 0 {
			opts.NamespaceID = gongfeng.Ptr(projCreateNamespaceID)
		}

		project, _, err := apiClient.Projects.CreateProject(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return printSuccessResponse(
			fmt.Sprintf("%d", project.ID),
			project.WebURL,
			project.PathWithNamespace,
		)
	},
}

var projectSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索项目",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		projects, _, err := apiClient.Projects.SearchProjects(ctx, args[0])
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, projects, !flagPretty)
	},
}

var projectMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "获取项目成员列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListProjectMembersOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    projMembersPage,
				PerPage: projMembersPerPage,
			},
		}

		members, _, err := apiClient.Projects.ListProjectMembers(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, members, !flagPretty)
	},
}

func init() {
	// project list flags
	projectListCmd.Flags().StringVar(&projListSearch, "search", "", "按关键字搜索项目")
	projectListCmd.Flags().StringVar(&projListOrderBy, "order-by", "", "排序字段（id, name, path, created_at, updated_at, last_activity_at）")
	projectListCmd.Flags().StringVar(&projListSort, "sort", "", "排序方向（asc, desc）")
	projectListCmd.Flags().IntVar(&projListPage, "page", 0, "页码")
	projectListCmd.Flags().IntVar(&projListPerPage, "per-page", 0, "每页条数")

	// project create flags
	projectCreateCmd.Flags().StringVar(&projCreateName, "name", "", "项目名称（必需）")
	projectCreateCmd.Flags().StringVar(&projCreatePath, "path", "", "项目路径")
	projectCreateCmd.Flags().StringVar(&projCreateDescription, "description", "", "项目描述")
	projectCreateCmd.Flags().IntVar(&projCreateNamespaceID, "namespace-id", 0, "命名空间 ID")

	// project members flags
	projectMembersCmd.Flags().IntVar(&projMembersPage, "page", 0, "页码")
	projectMembersCmd.Flags().IntVar(&projMembersPerPage, "per-page", 0, "每页条数")

	projectCmd.AddCommand(projectListCmd, projectShowCmd, projectCreateCmd, projectSearchCmd, projectMembersCmd)
	rootCmd.AddCommand(projectCmd)
}
