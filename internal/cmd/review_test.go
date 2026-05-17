package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestReviewMRShow(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_request/1/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 10, "state": "approved",
		})
	})

	_, err := executeCommand("review", "mr-show", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewMRReopen(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_request/1/review/reopen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 10, "state": "opened"})
	})

	_, err := executeCommand("review", "mr-reopen", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewMRCancel(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_request/1/review/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("review", "mr-cancel", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewCommitListWithAuthorID(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("author_id") != "5" {
			t.Errorf("expected author_id=5, got %s", r.URL.Query().Get("author_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("review", "commit-list", "--author-id", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewCommitCreateWithNewParams(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "review title" {
			t.Errorf("expected title='review title', got %v", body["title"])
		}
		if body["source_commit"] != "abc123" {
			t.Errorf("expected source_commit=abc123, got %v", body["source_commit"])
		}
		if body["source_branch"] != "feature" {
			t.Errorf("expected source_branch=feature, got %v", body["source_branch"])
		}
		if body["target_branch"] != "main" {
			t.Errorf("expected target_branch=main, got %v", body["target_branch"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "title": "review title",
		})
	})

	_, err := executeCommand("review", "commit-create", "--title", "review title",
		"--source-commit", "abc123", "--source-branch", "feature", "--target-branch", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewCommitUpdate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/review/10", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "updated title" {
			t.Errorf("expected title='updated title', got %v", body["title"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 10, "title": "updated title",
		})
	})

	_, err := executeCommand("review", "commit-update", "10", "--title", "updated title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewCommitReopen(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/review/10/reopen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 10, "state": "opened"})
	})

	_, err := executeCommand("review", "commit-reopen", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
