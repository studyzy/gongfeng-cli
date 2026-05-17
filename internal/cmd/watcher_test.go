package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestWatcherStatus(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/watch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(true)
	})

	_, err := executeCommand("watcher", "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatcherWatchWithMute(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/watch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["mute"] != true {
			t.Errorf("expected mute=true, got %v", body["mute"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"user_id": 1, "project_id": 123})
	})

	_, err := executeCommand("watcher", "watch", "--mute")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
