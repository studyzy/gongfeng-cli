package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestIssueListWithNewFilters(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/issues", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("resolve_state") != "resolved" {
			t.Errorf("expected resolve_state=resolved, got %s", q.Get("resolve_state"))
		}
		if q.Get("grade") != "3" {
			t.Errorf("expected grade=3, got %s", q.Get("grade"))
		}
		if q.Get("created_after") != "2025-01-01T00:00:00Z" {
			t.Errorf("expected created_after param, got %s", q.Get("created_after"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("issue", "list", "--resolve-state", "resolved", "--grade", "3", "--created-after", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueMyList(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/issues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "title": "my issue"},
		})
	})

	_, err := executeCommand("issue", "my-list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCreateWithGrade(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/issues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["grade"] != float64(2) {
			t.Errorf("expected grade=2, got %v", body["grade"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "title": "new issue", "grade": 2,
		})
	})

	_, err := executeCommand("issue", "create", "--title", "new issue", "--grade", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueSubscribe(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/issues/42/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("issue", "subscribe", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueUnsubscribe(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/issues/42/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("issue", "unsubscribe", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
