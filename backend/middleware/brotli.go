package middleware

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func BrotliStatic(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only handle GET and HEAD
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		safeDir := os.Getenv("PUBLIC_DIR")
		// Check if client accepts Brotli
		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			brPath := filepath.Join(dir, r.URL.Path) + ".br"

			if !strings.HasPrefix(brPath, safeDir) {
				http.Error(w, "Invalid file name", http.StatusBadRequest)
				log.Printf("Incorrect File: %s prefix: %s\n",brPath, safeDir)
				return
			}

			if _, err := os.Stat(brPath); err == nil {
				w.Header().Set("Content-Encoding", "br")
				w.Header().Set("Content-Type", detectMimeType(r.URL.Path))
				http.ServeFile(w, r, brPath)
				return
			}
		}

		path := filepath.Join(dir, r.URL.Path)

		if !strings.HasPrefix(path, safeDir) {
			http.Error(w, "Invalid file name", http.StatusBadRequest)
			log.Printf("Incorrect File: %s prefix: %s\n",path, safeDir)
			return
		}

		if _, err := os.Stat(path); err == nil {
			http.ServeFile(w, r, path)
		} else {
			http.NotFound(w, r)
		}
	})
}

func detectMimeType(path string) string {
	switch {
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".html"):
		return "text/html"
	default:
		return "application/octet-stream"
	}
}
