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
	projListSearch           string
	projListOrderBy          string
	projListSort             string
	projListPage             int
	projListPerPage          int
	projListArchived         bool
	projListWithArchived     bool
	projListWithPush         bool
	projListAbandoned        bool
	projListVisibilityLevels string
)

// project create flags
var (
	projCreateName            string
	projCreatePath            string
	projCreateDescription     string
	projCreateNamespaceID     int
	projCreateForkEnabled     bool
	projCreateVisibilityLevel int
)

// project update flags
var (
	projUpdateName                 string
	projUpdateDescription          string
	projUpdateDefaultBranch        string
	projUpdateIssuesEnabled        bool
	projUpdateMergeRequestsEnabled bool
	projUpdateWikiEnabled          bool
	projUpdateReviewEnabled        bool
	projUpdateForkEnabled          bool
	projUpdateVisibilityLevel      int
)

// project members flags
var (
	projMembersPage    int
	projMembersPerPage int
	projMembersQuery   string
)

// project share flags
var (
	projShareGroupID     int
	projShareGroupAccess int
)

// project events flags
var (
	projEventsPage    int
	projEventsPerPage int
	projEventsUser    string
)

// project stars flags
var (
	projStarsPage    int
	projStarsPerPage int
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
		if projListArchived {
			opts.Archived = gongfeng.Ptr(true)
		}
		if projListWithArchived {
			opts.WithArchived = gongfeng.Ptr(true)
		}
		if projListWithPush {
			opts.WithPush = gongfeng.Ptr(true)
		}
		if projListAbandoned {
			opts.Abandoned = gongfeng.Ptr(true)
		}
		if projListVisibilityLevels != "" {
			opts.VisibilityLevels = gongfeng.Ptr(projListVisibilityLevels)
		}

		projects, _, err := apiClient.Projects.ListProjects(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, projects, !flagPretty)
	},
}

var projectOwnedCmd = &cobra.Command{
	Use:   "owned",
	Short: "获取当前用户拥有的项目列表",
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

		projects, _, err := apiClient.Projects.ListOwnedProjects(ctx, opts)
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
		if projCreateForkEnabled {
			opts.ForkEnabled = gongfeng.Ptr(true)
		}
		if projCreateVisibilityLevel != 0 {
			opts.VisibilityLevel = gongfeng.Ptr(projCreateVisibilityLevel)
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

var projectUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新项目设置",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		opts := &gongfeng.UpdateProjectOptions{}
		if projUpdateName != "" {
			opts.Name = gongfeng.Ptr(projUpdateName)
		}
		if projUpdateDescription != "" {
			opts.Description = gongfeng.Ptr(projUpdateDescription)
		}
		if projUpdateDefaultBranch != "" {
			opts.DefaultBranch = gongfeng.Ptr(projUpdateDefaultBranch)
		}
		if cmd.Flags().Changed("issues-enabled") {
			opts.IssuesEnabled = gongfeng.Ptr(projUpdateIssuesEnabled)
		}
		if cmd.Flags().Changed("merge-requests-enabled") {
			opts.MergeRequestsEnabled = gongfeng.Ptr(projUpdateMergeRequestsEnabled)
		}
		if cmd.Flags().Changed("wiki-enabled") {
			opts.WikiEnabled = gongfeng.Ptr(projUpdateWikiEnabled)
		}
		if cmd.Flags().Changed("review-enabled") {
			opts.ReviewEnabled = gongfeng.Ptr(projUpdateReviewEnabled)
		}
		if cmd.Flags().Changed("fork-enabled") {
			opts.ForkEnabled = gongfeng.Ptr(projUpdateForkEnabled)
		}
		if projUpdateVisibilityLevel != 0 {
			opts.VisibilityLevel = gongfeng.Ptr(projUpdateVisibilityLevel)
		}

		project, _, err := apiClient.Projects.UpdateProject(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(project, "description")
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		_, err := apiClient.Projects.DeleteProject(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var projectSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索项目",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListProjectsOptions{
			Search: gongfeng.Ptr(args[0]),
		}
		projects, _, err := apiClient.Projects.ListProjects(ctx, opts)
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
		if projMembersQuery != "" {
			opts.Query = gongfeng.Ptr(projMembersQuery)
		}

		members, _, err := apiClient.Projects.ListProjectMembers(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, members, !flagPretty)
	},
}

var projectMemberShowCmd = &cobra.Command{
	Use:   "member-show <user_id>",
	Short: "获取项目成员详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		userID := atoi(args[0], "user_id")
		ctx := context.Background()
		member, _, err := apiClient.Projects.GetProjectMember(ctx, projectID(), userID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, member, !flagPretty)
	},
}

var projectShareCmd = &cobra.Command{
	Use:   "share",
	Short: "将项目共享给组",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		if projShareGroupID == 0 {
			exitWithParamError("--group-id is required", "Specify the group ID to share with")
		}
		ctx := context.Background()
		opts := &gongfeng.ShareProjectOptions{
			GroupID: gongfeng.Ptr(projShareGroupID),
		}
		if projShareGroupAccess != 0 {
			opts.GroupAccess = gongfeng.Ptr(projShareGroupAccess)
		}
		_, err := apiClient.Projects.ShareProject(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var projectSharesCmd = &cobra.Command{
	Use:   "shares",
	Short: "获取项目共享的组列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		shares, _, err := apiClient.Projects.ListProjectShares(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, shares, !flagPretty)
	},
}

var projectUnshareCmd = &cobra.Command{
	Use:   "unshare <group_id>",
	Short: "删除项目共享关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		groupID := atoi(args[0], "group_id")
		ctx := context.Background()
		_, err := apiClient.Projects.DeleteProjectShare(ctx, projectID(), groupID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var projectEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "获取项目事件列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		opts := &gongfeng.ListProjectEventsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    projEventsPage,
				PerPage: projEventsPerPage,
			},
		}
		if projEventsUser != "" {
			opts.UserIDOrName = projEventsUser
		}
		events, _, err := apiClient.Projects.ListProjectEvents(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, events, !flagPretty)
	},
}

var projectStarCmd = &cobra.Command{
	Use:   "star",
	Short: "标星项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		_, _, err := apiClient.Projects.StarProject(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var projectUnstarCmd = &cobra.Command{
	Use:   "unstar",
	Short: "取消标星项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		_, err := apiClient.Projects.UnstarProject(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var projectStarStatusCmd = &cobra.Command{
	Use:   "star-status",
	Short: "查询项目标星状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		starred, _, err := apiClient.Projects.GetStarStatus(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, map[string]bool{"starred": starred}, !flagPretty)
	},
}

var projectStarsCmd = &cobra.Command{
	Use:   "stars",
	Short: "获取项目标星用户列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()
		ctx := context.Background()
		opts := &gongfeng.ListProjectStarsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    projStarsPage,
				PerPage: projStarsPerPage,
			},
		}
		stars, _, err := apiClient.Projects.ListProjectStars(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, stars, !flagPretty)
	},
}

func init() {
	// project list flags
	projectListCmd.Flags().StringVar(&projListSearch, "search", "", "按关键字搜索项目")
	projectListCmd.Flags().StringVar(&projListOrderBy, "order-by", "", "排序字段（id, name, path, created_at, updated_at, last_activity_at）")
	projectListCmd.Flags().StringVar(&projListSort, "sort", "", "排序方向（asc, desc）")
	projectListCmd.Flags().IntVar(&projListPage, "page", 0, "页码")
	projectListCmd.Flags().IntVar(&projListPerPage, "per-page", 0, "每页条数")
	projectListCmd.Flags().BoolVar(&projListArchived, "archived", false, "仅显示已归档项目")
	projectListCmd.Flags().BoolVar(&projListWithArchived, "with-archived", false, "包含已归档项目")
	projectListCmd.Flags().BoolVar(&projListWithPush, "with-push", false, "仅显示有推送的项目")
	projectListCmd.Flags().BoolVar(&projListAbandoned, "abandoned", false, "仅显示废弃项目")
	projectListCmd.Flags().StringVar(&projListVisibilityLevels, "visibility-levels", "", "可见性级别过滤")

	// project owned flags
	projectOwnedCmd.Flags().IntVar(&projListPage, "page", 0, "页码")
	projectOwnedCmd.Flags().IntVar(&projListPerPage, "per-page", 0, "每页条数")
	projectOwnedCmd.Flags().StringVar(&projListSearch, "search", "", "按关键字搜索")

	// project create flags
	projectCreateCmd.Flags().StringVar(&projCreateName, "name", "", "项目名称（必需）")
	projectCreateCmd.Flags().StringVar(&projCreatePath, "path", "", "项目路径")
	projectCreateCmd.Flags().StringVar(&projCreateDescription, "description", "", "项目描述")
	projectCreateCmd.Flags().IntVar(&projCreateNamespaceID, "namespace-id", 0, "命名空间 ID")
	projectCreateCmd.Flags().BoolVar(&projCreateForkEnabled, "fork-enabled", false, "允许 Fork")
	projectCreateCmd.Flags().IntVar(&projCreateVisibilityLevel, "visibility-level", 0, "可见性级别")

	// project update flags
	projectUpdateCmd.Flags().StringVar(&projUpdateName, "name", "", "项目名称")
	projectUpdateCmd.Flags().StringVar(&projUpdateDescription, "description", "", "项目描述")
	projectUpdateCmd.Flags().StringVar(&projUpdateDefaultBranch, "default-branch", "", "默认分支")
	projectUpdateCmd.Flags().BoolVar(&projUpdateIssuesEnabled, "issues-enabled", false, "启用缺陷管理")
	projectUpdateCmd.Flags().BoolVar(&projUpdateMergeRequestsEnabled, "merge-requests-enabled", false, "启用合并请求")
	projectUpdateCmd.Flags().BoolVar(&projUpdateWikiEnabled, "wiki-enabled", false, "启用 Wiki")
	projectUpdateCmd.Flags().BoolVar(&projUpdateReviewEnabled, "review-enabled", false, "启用代码评审")
	projectUpdateCmd.Flags().BoolVar(&projUpdateForkEnabled, "fork-enabled", false, "允许 Fork")
	projectUpdateCmd.Flags().IntVar(&projUpdateVisibilityLevel, "visibility-level", 0, "可见性级别")

	// project members flags
	projectMembersCmd.Flags().IntVar(&projMembersPage, "page", 0, "页码")
	projectMembersCmd.Flags().IntVar(&projMembersPerPage, "per-page", 0, "每页条数")
	projectMembersCmd.Flags().StringVar(&projMembersQuery, "query", "", "按用户名搜索")

	// project share flags
	projectShareCmd.Flags().IntVar(&projShareGroupID, "group-id", 0, "组 ID（必需）")
	projectShareCmd.Flags().IntVar(&projShareGroupAccess, "group-access", 0, "组访问级别")

	// project events flags
	projectEventsCmd.Flags().IntVar(&projEventsPage, "page", 0, "页码")
	projectEventsCmd.Flags().IntVar(&projEventsPerPage, "per-page", 0, "每页条数")
	projectEventsCmd.Flags().StringVar(&projEventsUser, "user", "", "按用户 ID 或用户名过滤")

	// project stars flags
	projectStarsCmd.Flags().IntVar(&projStarsPage, "page", 0, "页码")
	projectStarsCmd.Flags().IntVar(&projStarsPerPage, "per-page", 0, "每页条数")

	projectCmd.AddCommand(
		projectListCmd, projectOwnedCmd, projectShowCmd, projectCreateCmd, projectUpdateCmd, projectDeleteCmd,
		projectSearchCmd, projectMembersCmd, projectMemberShowCmd,
		projectShareCmd, projectSharesCmd, projectUnshareCmd,
		projectEventsCmd, projectStarCmd, projectUnstarCmd, projectStarStatusCmd, projectStarsCmd,
	)
	rootCmd.AddCommand(projectCmd)
}
