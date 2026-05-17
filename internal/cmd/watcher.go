package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// watcher flags
var (
	waFlagPage    int
	waFlagPerPage int
	waFlagMute    bool
)

var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "项目关注管理",
}

var watcherListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取关注者列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		opts := &gongfeng.ListWatchersOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    waFlagPage,
				PerPage: waFlagPerPage,
			},
		}

		watchers, _, err := apiClient.Watchers.ListWatchers(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, watchers, !flagPretty)
	},
}

var watcherStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查询当前用户是否关注项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		watched, _, err := apiClient.Watchers.GetWatchStatus(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, map[string]bool{"watched": watched}, !flagPretty)
	},
}

var watcherWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "关注项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()
		var opts *gongfeng.WatchProjectOptions
		if waFlagMute {
			opts = &gongfeng.WatchProjectOptions{
				Mute: gongfeng.Ptr(true),
			}
		}

		_, _, err := apiClient.Watchers.WatchProject(ctx, projectID(), opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

var watcherUnwatchCmd = &cobra.Command{
	Use:   "unwatch",
	Short: "取消关注项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireProjectID()

		ctx := context.Background()

		_, err := apiClient.Watchers.UnwatchProject(ctx, projectID())
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, &output.SuccessResponse{Success: true}, !flagPretty)
	},
}

func init() {
	// watcher list flags
	watcherListCmd.Flags().IntVar(&waFlagPage, "page", 0, "页码")
	watcherListCmd.Flags().IntVar(&waFlagPerPage, "per-page", 0, "每页条数")

	// watcher watch flags
	watcherWatchCmd.Flags().BoolVar(&waFlagMute, "mute", false, "静默关注（不接收通知）")

	watcherCmd.AddCommand(watcherListCmd, watcherStatusCmd, watcherWatchCmd, watcherUnwatchCmd)
	rootCmd.AddCommand(watcherCmd)
}
