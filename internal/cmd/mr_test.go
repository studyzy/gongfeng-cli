package cmd

import (
	"net/http"
	"testing"
)

func TestMRDownloadFiles(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_requests/1/changed_files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("diff content here"))
	})

	_, err := executeCommand("mr", "download-files", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
