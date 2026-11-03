package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	if got := newServer("").Addr; got != ":8080" {
		t.Fatalf("default Addr = %q", got)
	}
	srv := newServer("127.0.0.1:9999")
	if srv.Addr != "127.0.0.1:9999" || srv.ReadHeaderTimeout == 0 {
		t.Fatalf("server = %+v", srv)
	}
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/healthz = %d", rec.Code)
	}
}
