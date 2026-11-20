// Package api is the HTTP API of the sample service.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// MaxNameLength is the longest name /greet accepts.
const MaxNameLength = 64

// Info is what /version reports.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// NewRouter returns the service routes.
func NewRouter(info Info) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", health)
	r.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, info)
	})
	r.Get("/greet", greet)
	return r
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// greet answers "hello, <name>"; name is required, trimmed and at most MaxNameLength characters.
func greet(w http.ResponseWriter, r *http.Request) {
	name, err := ValidateName(r.URL.Query().Get("name"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello, " + name})
}

// ValidateName trims name and checks it is present and not too long.
func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	switch {
	case name == "":
		return "", errNameMissing
	case len([]rune(name)) > MaxNameLength:
		return "", errNameTooLong
	}
	return name, nil
}

type apiError string

func (e apiError) Error() string { return string(e) }

const (
	errNameMissing = apiError("name is required")
	errNameTooLong = apiError("name is too long")
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
