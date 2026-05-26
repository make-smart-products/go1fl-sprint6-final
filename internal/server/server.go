package server

import (
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// New creates an HTTP server for the sports statistics site.
func New(addr string) *http.Server {
	mux := http.NewServeMux()
	handlers.Register(mux)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
