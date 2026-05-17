package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"
	gongfeng "github.com/studyzy/gongfeng-sdk-go"
)

// setupTest creates a mock HTTP server and configures the global apiClient.
// It bypasses the PersistentPreRunE authentication check.
func setupTest(t *testing.T) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := gongfeng.NewClient("test-token", gongfeng.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	apiClient = client

	// Set a default project ID for tests
	flagProjectID = "123"
	flagPretty = false
	flagJSON = false
	flagToken = "test-token"

	// Override PersistentPreRunE to skip auth
	rootCmd.PersistentPreRunE = nil

	t.Cleanup(func() {
		// Restore PersistentPreRunE after test
		rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "login" {
				return nil
			}
			if v, _ := cmd.Flags().GetBool("version"); v {
				return nil
			}
			return initClientAndConfig(cmd)
		}
		flagToken = ""
	})

	return mux
}

// executeCommand executes a cobra command with the given args and returns stdout output.
func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()

	// Reset args
	rootCmd.SetArgs(nil)

	return buf.String(), err
}
