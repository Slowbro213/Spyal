package middleware

import (
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

		// Check if client accepts Brotli
		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			brPath := filepath.Join(dir, r.URL.Path) + ".br"
			if _, err := os.Stat(brPath); err == nil {
				w.Header().Set("Content-Encoding", "br")
				w.Header().Set("Content-Type", detectMimeType(r.URL.Path))
				http.ServeFile(w, r, brPath)
				return
			}
		}

		// Fallback to normal file
		path := filepath.Join(dir, r.URL.Path)
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
