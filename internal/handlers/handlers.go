package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Register wires the static sports dashboard routes.
func Register(mux *http.ServeMux) {
	root := projectRoot()
	files := http.FileServer(http.Dir(root))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFile(w, r, filepath.Join(root, "index.html"))
			return
		}

		cleanPath := filepath.Clean(r.URL.Path)
		if strings.HasPrefix(cleanPath, "/.") {
			http.NotFound(w, r)
			return
		}

		if strings.HasPrefix(cleanPath, "/data/") {
			w.Header().Set("Cache-Control", "no-cache")
		}

		files.ServeHTTP(w, r)
	})
}

func projectRoot() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	for _, candidate := range []string{
		workingDirectory,
		filepath.Dir(workingDirectory),
		filepath.Dir(filepath.Dir(workingDirectory)),
	} {
		if _, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil {
			return candidate
		}
	}

	return workingDirectory
}
