package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRepoCreateFileWithEncoding(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["encoding"] != "base64" {
			t.Errorf("expected encoding=base64, got %v", body["encoding"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"file_path": "test.bin", "branch_name": "main",
		})
	})

	_, err := executeCommand("repo", "create-file",
		"--file-path", "test.bin",
		"--branch-name", "main",
		"--content", "SGVsbG8=",
		"--commit-message", "add binary",
		"--encoding", "base64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepoCompareWithStraight(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/compare", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("straight") != "true" {
			t.Errorf("expected straight=true, got %s", q.Get("straight"))
		}
		if q.Get("from") != "main" {
			t.Errorf("expected from=main, got %s", q.Get("from"))
		}
		if q.Get("to") != "feature" {
			t.Errorf("expected to=feature, got %s", q.Get("to"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"commits": []interface{}{}, "diffs": []interface{}{},
		})
	})

	_, err := executeCommand("repo", "compare", "--from", "main", "--to", "feature", "--straight")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepoRaw(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/blobs/abc123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Write([]byte("raw file content"))
	})

	_, err := executeCommand("repo", "raw", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepoCommitBlob(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/repository/commits/abc123/blob", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("filepath") != "src/main.go" {
			t.Errorf("expected filepath=src/main.go, got %s", r.URL.Query().Get("filepath"))
		}
		w.Write([]byte("package main"))
	})

	_, err := executeCommand("repo", "commit-blob", "abc123", "--filepath", "src/main.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
