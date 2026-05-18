package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// webhook flags
var (
	whFlagURL           string
	whFlagPushEvents    bool
	whFlagIssuesEvents  bool
	whFlagMREvents      bool
	whFlagTagPushEvents bool
	whFlagNoteEvents    bool
	whFlagSSL           bool
	whFlagPage          int
	whFlagPerPage       int
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Webhook 管理",
}

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取 Webhook 列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListWebhooksOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    whFlagPage,
				PerPage: whFlagPerPage,
			},
		}

		hooks, _, err := apiClient.Webhooks.ListWebhooks(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, hooks, !flagPretty)
	},
}

var webhookShowCmd = &cobra.Command{
	Use:   "show <hook_id>",
	Short: "获取 Webhook 详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		hookID := atoi(args[0], "hook_id")
		ctx := context.Background()

		hook, _, err := apiClient.Webhooks.GetWebhook(ctx, projectID(), hookID)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, hook, !flagPretty)
	},
}

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建 Webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		if whFlagURL == "" {
			exitWithParamError("--url is required", "Provide a webhook URL with --url")
		}

		ctx := context.Background()
		opts := &gongfeng.AddWebhookOptions{
			URL:                   gongfeng.Ptr(whFlagURL),
			PushEvents:            gongfeng.Ptr(whFlagPushEvents),
			IssuesEvents:          gongfeng.Ptr(whFlagIssuesEvents),
			MergeRequestsEvents:   gongfeng.Ptr(whFlagMREvents),
			TagPushEvents:         gongfeng.Ptr(whFlagTagPushEvents),
			NoteEvents:            gongfeng.Ptr(whFlagNoteEvents),
			EnableSSLVerification: gongfeng.Ptr(whFlagSSL),
		}

		hook, _, err := apiClient.Webhooks.AddWebhook(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, hook, !flagPretty)
	},
}

var webhookUpdateCmd = &cobra.Command{
	Use:   "update <hook_id>",
	Short: "编辑 Webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		hookID := atoi(args[0], "hook_id")
		ctx := context.Background()

		opts := &gongfeng.EditWebhookOptions{}
		if cmd.Flags().Changed("url") {
			opts.URL = gongfeng.Ptr(whFlagURL)
		}
		if cmd.Flags().Changed("push-events") {
			opts.PushEvents = gongfeng.Ptr(whFlagPushEvents)
		}
		if cmd.Flags().Changed("issues-events") {
			opts.IssuesEvents = gongfeng.Ptr(whFlagIssuesEvents)
		}
		if cmd.Flags().Changed("merge-requests-events") {
			opts.MergeRequestsEvents = gongfeng.Ptr(whFlagMREvents)
		}
		if cmd.Flags().Changed("tag-push-events") {
			opts.TagPushEvents = gongfeng.Ptr(whFlagTagPushEvents)
		}
		if cmd.Flags().Changed("note-events") {
			opts.NoteEvents = gongfeng.Ptr(whFlagNoteEvents)
		}
		if cmd.Flags().Changed("enable-ssl-verification") {
			opts.EnableSSLVerification = gongfeng.Ptr(whFlagSSL)
		}

		hook, _, err := apiClient.Webhooks.EditWebhook(ctx, projectID(), hookID, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, hook, !flagPretty)
	},
}

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete <hook_id>",
	Short: "删除 Webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		hookID := atoi(args[0], "hook_id")
		ctx := context.Background()

		_, err := apiClient.Webhooks.DeleteWebhook(ctx, projectID(), hookID)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// webhook list flags
	webhookListCmd.Flags().IntVar(&whFlagPage, "page", 0, "页码")
	webhookListCmd.Flags().IntVar(&whFlagPerPage, "per-page", 0, "每页条数")

	// webhook create flags
	webhookCreateCmd.Flags().StringVar(&whFlagURL, "url", "", "Webhook URL（必需）")
	webhookCreateCmd.Flags().BoolVar(&whFlagPushEvents, "push-events", false, "推送事件触发")
	webhookCreateCmd.Flags().BoolVar(&whFlagIssuesEvents, "issues-events", false, "Issue 事件触发")
	webhookCreateCmd.Flags().BoolVar(&whFlagMREvents, "merge-requests-events", false, "合并请求事件触发")
	webhookCreateCmd.Flags().BoolVar(&whFlagTagPushEvents, "tag-push-events", false, "Tag 推送事件触发")
	webhookCreateCmd.Flags().BoolVar(&whFlagNoteEvents, "note-events", false, "评论事件触发")
	webhookCreateCmd.Flags().BoolVar(&whFlagSSL, "enable-ssl-verification", false, "启用 SSL 验证")

	// webhook update flags
	webhookUpdateCmd.Flags().StringVar(&whFlagURL, "url", "", "Webhook URL")
	webhookUpdateCmd.Flags().BoolVar(&whFlagPushEvents, "push-events", false, "推送事件触发")
	webhookUpdateCmd.Flags().BoolVar(&whFlagIssuesEvents, "issues-events", false, "Issue 事件触发")
	webhookUpdateCmd.Flags().BoolVar(&whFlagMREvents, "merge-requests-events", false, "合并请求事件触发")
	webhookUpdateCmd.Flags().BoolVar(&whFlagTagPushEvents, "tag-push-events", false, "Tag 推送事件触发")
	webhookUpdateCmd.Flags().BoolVar(&whFlagNoteEvents, "note-events", false, "评论事件触发")
	webhookUpdateCmd.Flags().BoolVar(&whFlagSSL, "enable-ssl-verification", false, "启用 SSL 验证")

	webhookCmd.AddCommand(webhookListCmd, webhookShowCmd, webhookCreateCmd, webhookUpdateCmd, webhookDeleteCmd)
	rootCmd.AddCommand(webhookCmd)
}
