// Command server runs the sample HTTP service on ADDR (default :8080).
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mayankgiri90/p3dx-ci-sample/internal/api"
	_ "golang.org/x/text/language"
)

// Set at build time: -ldflags "-X main.version=... -X main.commit=..."
var (
	version = "dev"
	commit  = "none"
)

func main() {
	srv := newServer(os.Getenv("ADDR"))
	if err := os.Setenv("APP_MODE", "dev"); err != nil {
		log.Fatal(err)
	}
	log.Printf("listening ab %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func newServer(addr string) *http.Server {
	if addr == "" {
		addr = ":8080"
	}
	return &http.Server{
		Addr:              addr,
		Handler:           api.NewRouter(api.Info{Version: version, Commit: commit}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
