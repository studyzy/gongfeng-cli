package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCommitDiffWithPath(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits/abc123/diff", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("path") != "src/main.go" {
			t.Errorf("expected path=src/main.go, got %s", r.URL.Query().Get("path"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"old_path": "src/main.go", "new_path": "src/main.go", "diff": "@@ -1 +1 @@\n-old\n+new"},
		})
	})

	_, err := executeCommand("commit", "diff", "abc123", "--path", "src/main.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitDiffIgnoreWhiteSpace(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits/abc123/diff", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ignore_white_space") != "true" {
			t.Errorf("expected ignore_white_space=true, got %s", r.URL.Query().Get("ignore_white_space"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("commit", "diff", "abc123", "--ignore-white-space")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitListWithPath(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "README.md" {
			t.Errorf("expected path=README.md, got %s", r.URL.Query().Get("path"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("commit", "list", "--path", "README.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitRefsWithType(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits/abc123/refs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("type") != "branch" {
			t.Errorf("expected type=branch, got %s", r.URL.Query().Get("type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("commit", "refs", "abc123", "--type", "branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitCreateComment(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits/abc123/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["note"] != "test comment" {
			t.Errorf("expected note=test comment, got %v", body["note"])
		}
		if body["path"] != "main.go" {
			t.Errorf("expected path=main.go, got %v", body["path"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"note": "test comment", "path": "main.go", "line": 10,
		})
	})

	_, err := executeCommand("commit", "create-comment", "abc123", "--note", "test comment", "--path", "main.go", "--line", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
