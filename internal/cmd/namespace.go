package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"

	"github.com/studyzy/gongfeng-cli/internal/output"
)

// namespace list flags
var (
	nsListSearch  string
	nsListPage    int
	nsListPerPage int
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "命名空间管理",
}

var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取命名空间列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := &gongfeng.ListNamespacesOptions{
			ListOptions: gongfeng.ListOptions{
				Page:    nsListPage,
				PerPage: nsListPerPage,
			},
		}
		if nsListSearch != "" {
			opts.Search = gongfeng.Ptr(nsListSearch)
		}

		namespaces, _, err := apiClient.Namespaces.ListNamespaces(ctx, opts)
		if err != nil {
			exitWithAPIError(err)
		}

		return output.PrintJSON(os.Stdout, namespaces, !flagPretty)
	},
}

func init() {
	// namespace list flags
	namespaceListCmd.Flags().StringVar(&nsListSearch, "search", "", "按关键字搜索命名空间")
	namespaceListCmd.Flags().IntVar(&nsListPage, "page", 0, "页码")
	namespaceListCmd.Flags().IntVar(&nsListPerPage, "per-page", 0, "每页条数")

	namespaceCmd.AddCommand(namespaceListCmd)
	rootCmd.AddCommand(namespaceCmd)
}
