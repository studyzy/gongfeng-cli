package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// user list flags
var (
	userListSearch  string
	userListPage    int
	userListPerPage int
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "用户管理",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取用户列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListUsersOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    userListPage,
				PerPage: userListPerPage,
			},
		}
		if userListSearch != "" {
			opts.Search = gongfeng.Ptr(userListSearch)
		}

		users, _, err := apiClient.Users.ListUsers(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, users, !flagPretty)
	},
}

var userShowCmd = &cobra.Command{
	Use:   "show <user_id>",
	Short: "获取用户详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		uid := atoi(args[0], "user_id")

		user, _, err := apiClient.Users.GetUser(ctx, uid)
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(user, "bio")
	},
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "获取当前认证用户信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		user, _, err := apiClient.Users.GetCurrentUser(ctx)
		if err != nil {
			exitWithAPIError(err)
		}

		return printDetail(user, "bio")
	},
}

func init() {
	// user list flags
	userListCmd.Flags().StringVar(&userListSearch, "search", "", "按关键字搜索用户")
	userListCmd.Flags().IntVar(&userListPage, "page", 0, "页码")
	userListCmd.Flags().IntVar(&userListPerPage, "per-page", 0, "每页条数")

	userCmd.AddCommand(userListCmd, userShowCmd, userMeCmd)
	rootCmd.AddCommand(userCmd)
}
