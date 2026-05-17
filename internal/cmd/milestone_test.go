package cmd

import (
	"net/http"
	"testing"
)

func TestMilestoneDelete(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/milestones/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("milestone", "delete", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
