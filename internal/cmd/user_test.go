package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestUserWatched(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/watched", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "project1"},
		})
	})

	_, err := executeCommand("user", "watched")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserSSHKeys(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "title": "my-key"},
		})
	})

	_, err := executeCommand("user", "ssh-keys")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserSSHKeyCreate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "test-key" {
			t.Errorf("expected title=test-key, got %v", body["title"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "title": "test-key",
		})
	})

	_, err := executeCommand("user", "ssh-key-create", "--title", "test-key", "--key", "ssh-rsa AAAA...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserSSHKeyDelete(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/keys/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := executeCommand("user", "ssh-key-delete", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserEmails(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/emails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "email": "user@example.com"},
		})
	})

	_, err := executeCommand("user", "emails")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserEmailCreate(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/emails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["email"] != "new@example.com" {
			t.Errorf("expected email=new@example.com, got %v", body["email"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "email": "new@example.com",
		})
	})

	_, err := executeCommand("user", "email-create", "--email", "new@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserFindByEmail(t *testing.T) {
	mux := setupTest(t)

	mux.HandleFunc("/api/v3/user/email", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("email") != "test@example.com" {
			t.Errorf("expected email=test@example.com, got %s", r.URL.Query().Get("email"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "username": "testuser",
		})
	})

	_, err := executeCommand("user", "find-by-email", "--email", "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
