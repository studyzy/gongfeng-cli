package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestProjectListWithNewFilters(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("archived") != "true" {
			t.Errorf("expected archived=true, got %s", q.Get("archived"))
		}
		if q.Get("with_push") != "true" {
			t.Errorf("expected with_push=true, got %s", q.Get("with_push"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("project", "list", "--archived", "--with-push")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectOwned(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/owned", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "my-project"},
		})
	})

	_, err := executeCommand("project", "owned")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectUpdate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "new-name" {
			t.Errorf("expected name=new-name, got %v", body["name"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 123, "name": "new-name",
		})
	})

	_, err := executeCommand("project", "update", "--name", "new-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectDelete(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("project", "delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectStar(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/star", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(true)
	})

	_, err := executeCommand("project", "star")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectUnstar(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/star", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := executeCommand("project", "unstar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectShare(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/share", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["group_id"] != float64(5) {
			t.Errorf("expected group_id=5, got %v", body["group_id"])
		}
		w.WriteHeader(http.StatusCreated)
	})

	_, err := executeCommand("project", "share", "--group-id", "5", "--group-access", "30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectEvents(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("project", "events")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProjectMembersWithQuery(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/projects/123/members", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("query") != "john" {
			t.Errorf("expected query=john, got %s", r.URL.Query().Get("query"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})

	_, err := executeCommand("project", "members", "--query", "john")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
