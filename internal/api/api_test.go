package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, map[string]string) {
	t.Helper()
	rec := httptest.NewRecorder()
	NewRouter(Info{Version: "1.0.0", Commit: "abc1234"}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if ct := rec.Header().Get("Content-Type"); rec.Code != http.StatusNotFound && ct != "application/json" {
		t.Fatalf("%s: Content-Type = %q", path, ct)
	}
	body := map[string]string{}
	if rec.Code != http.StatusNotFound {
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: bad JSON %q: %v", path, rec.Body.String(), err)
		}
	}
	return rec.Code, body
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		path   string
		status int
		key    string
		want   string
	}{
		{"/healthz", http.StatusOK, "status", "ok"},
		{"/version", http.StatusOK, "version", "1.0.0"},
		{"/greet?name=Rohith", http.StatusOK, "message", "hello, Rohith"},
		{"/greet?name=%20Sai%20", http.StatusOK, "message", "hello, Sai"},
		{"/greet", http.StatusBadRequest, "error", "name is required"},
		{"/greet?name=" + strings.Repeat("a", MaxNameLength+1), http.StatusBadRequest, "error", "name is too long"},
		{"/missing", http.StatusNotFound, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			status, body := get(t, tt.path)
			if status != tt.status {
				t.Fatalf("status = %d, want %d", status, tt.status)
			}
			if tt.key != "" && body[tt.key] != tt.want {
				t.Fatalf("%s = %q, want %q", tt.key, body[tt.key], tt.want)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	if got, err := ValidateName(strings.Repeat("é", MaxNameLength)); err != nil || len([]rune(got)) != MaxNameLength {
		t.Fatalf("64 runes: got %q, %v", got, err)
	}
	if _, err := ValidateName("   "); !errors.Is(err, errNameMissing) {
		t.Fatalf("blank: err = %v", err)
	}
}
