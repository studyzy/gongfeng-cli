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

// user ssh-key flags
var (
	userSSHKeyTitle string
	userSSHKeyKey   string
)

// user email flags
var (
	userEmailAddr string
)

// user watched flags
var (
	userWatchedPage    int
	userWatchedPerPage int
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

var userWatchedCmd = &cobra.Command{
	Use:   "watched",
	Short: "获取当前用户关注的项目列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListWatchedProjectsOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    userWatchedPage,
				PerPage: userWatchedPerPage,
			},
		}
		projects, _, err := apiClient.Users.ListWatchedProjects(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, projects, !flagPretty)
	},
}

var userSSHKeysCmd = &cobra.Command{
	Use:   "ssh-keys",
	Short: "获取当前用户 SSH Key 列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		keys, _, err := apiClient.Users.ListSSHKeys(ctx)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, keys, !flagPretty)
	},
}

var userSSHKeyShowCmd = &cobra.Command{
	Use:   "ssh-key-show <key_id>",
	Short: "获取指定 SSH Key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		keyID := atoi(args[0], "key_id")
		key, _, err := apiClient.Users.GetSSHKey(ctx, keyID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, key, !flagPretty)
	},
}

var userSSHKeyCreateCmd = &cobra.Command{
	Use:   "ssh-key-create",
	Short: "创建 SSH Key",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSSHKeyTitle == "" {
			exitWithParamError("--title is required", "Specify the SSH key title")
		}
		if userSSHKeyKey == "" {
			exitWithParamError("--key is required", "Specify the SSH public key")
		}
		ctx := context.Background()
		opts := &gongfeng.CreateSSHKeyOptions{
			Title: gongfeng.Ptr(userSSHKeyTitle),
			Key:   gongfeng.Ptr(userSSHKeyKey),
		}
		key, _, err := apiClient.Users.CreateSSHKey(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, key, !flagPretty)
	},
}

var userSSHKeyDeleteCmd = &cobra.Command{
	Use:   "ssh-key-delete <key_id>",
	Short: "删除 SSH Key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		keyID := atoi(args[0], "key_id")
		_, err := apiClient.Users.DeleteSSHKey(ctx, keyID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var userEmailsCmd = &cobra.Command{
	Use:   "emails",
	Short: "获取当前用户邮箱列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		emails, _, err := apiClient.Users.ListEmails(ctx)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, emails, !flagPretty)
	},
}

var userEmailShowCmd = &cobra.Command{
	Use:   "email-show <email_id>",
	Short: "获取指定邮箱",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		emailID := atoi(args[0], "email_id")
		email, _, err := apiClient.Users.GetEmail(ctx, emailID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, email, !flagPretty)
	},
}

var userEmailCreateCmd = &cobra.Command{
	Use:   "email-create",
	Short: "添加邮箱",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userEmailAddr == "" {
			exitWithParamError("--email is required", "Specify the email address")
		}
		ctx := context.Background()
		opts := &gongfeng.CreateEmailOptions{
			Email: gongfeng.Ptr(userEmailAddr),
		}
		email, _, err := apiClient.Users.CreateEmail(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, email, !flagPretty)
	},
}

var userEmailDeleteCmd = &cobra.Command{
	Use:   "email-delete <email_id>",
	Short: "删除邮箱",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		emailID := atoi(args[0], "email_id")
		_, err := apiClient.Users.DeleteEmail(ctx, emailID)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var userFindByEmailCmd = &cobra.Command{
	Use:   "find-by-email",
	Short: "通过邮箱查找用户",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userEmailAddr == "" {
			exitWithParamError("--email is required", "Specify the email address")
		}
		ctx := context.Background()
		opts := &gongfeng.GetUserByEmailOptions{
			Email: gongfeng.Ptr(userEmailAddr),
		}
		user, _, err := apiClient.Users.GetUserByEmail(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}
		return output.PrintJSON(os.Stdout, user, !flagPretty)
	},
}

func init() {
	// user list flags
	userListCmd.Flags().StringVar(&userListSearch, "search", "", "按关键字搜索用户")
	userListCmd.Flags().IntVar(&userListPage, "page", 0, "页码")
	userListCmd.Flags().IntVar(&userListPerPage, "per-page", 0, "每页条数")

	// user watched flags
	userWatchedCmd.Flags().IntVar(&userWatchedPage, "page", 0, "页码")
	userWatchedCmd.Flags().IntVar(&userWatchedPerPage, "per-page", 0, "每页条数")

	// ssh-key-create flags
	userSSHKeyCreateCmd.Flags().StringVar(&userSSHKeyTitle, "title", "", "SSH Key 标题（必需）")
	userSSHKeyCreateCmd.Flags().StringVar(&userSSHKeyKey, "key", "", "SSH 公钥内容（必需）")

	// email-create flags
	userEmailCreateCmd.Flags().StringVar(&userEmailAddr, "email", "", "邮箱地址（必需）")

	// find-by-email flags
	userFindByEmailCmd.Flags().StringVar(&userEmailAddr, "email", "", "邮箱地址（必需）")

	userCmd.AddCommand(
		userListCmd, userShowCmd, userMeCmd, userWatchedCmd,
		userSSHKeysCmd, userSSHKeyShowCmd, userSSHKeyCreateCmd, userSSHKeyDeleteCmd,
		userEmailsCmd, userEmailShowCmd, userEmailCreateCmd, userEmailDeleteCmd,
		userFindByEmailCmd,
	)
	rootCmd.AddCommand(userCmd)
}
