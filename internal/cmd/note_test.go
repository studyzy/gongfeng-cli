package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestMRNoteCreateWithNewParams(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_requests/1/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["body"] != "inline comment" {
			t.Errorf("expected body='inline comment', got %v", body["body"])
		}
		if body["path"] != "src/main.go" {
			t.Errorf("expected path=src/main.go, got %v", body["path"])
		}
		if body["line"] != "42" {
			t.Errorf("expected line=42, got %v", body["line"])
		}
		if body["line_type"] != "new" {
			t.Errorf("expected line_type=new, got %v", body["line_type"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "body": "inline comment",
		})
	})

	_, err := executeCommand("note", "mr-create", "1", "--body", "inline comment",
		"--path", "src/main.go", "--line", "42", "--line-type", "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMRNoteUpdateWithReviewerState(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/merge_requests/1/notes/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["reviewer_state"] != "resolved" {
			t.Errorf("expected reviewer_state=resolved, got %v", body["reviewer_state"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 5, "body": "updated",
		})
	})

	_, err := executeCommand("note", "mr-update", "1", "5", "--body", "updated", "--reviewer-state", "resolved")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewNoteList(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/reviews/10/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "body": "note 1"},
		})
	})

	_, err := executeCommand("note", "review-list", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewNoteCreate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/reviews/10/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["body"] != "review comment" {
			t.Errorf("expected body='review comment', got %v", body["body"])
		}
		if body["path"] != "lib.go" {
			t.Errorf("expected path=lib.go, got %v", body["path"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "body": "review comment",
		})
	})

	_, err := executeCommand("note", "review-create", "10", "--body", "review comment", "--path", "lib.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewNoteUpdate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/reviews/10/notes/2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["body"] != "updated review comment" {
			t.Errorf("expected body='updated review comment', got %v", body["body"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "body": "updated review comment",
		})
	})

	_, err := executeCommand("note", "review-update", "10", "2", "--body", "updated review comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
