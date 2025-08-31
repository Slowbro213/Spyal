package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"spyal/broadcasting"
	"spyal/core"
	"strings"

	"github.com/tdewolff/minify/v2"
	minhtml "github.com/tdewolff/minify/v2/html"
)

type bufferedWriter struct {
	core.Middleware
	buf *bytes.Buffer
}

func (w *bufferedWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *bufferedWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// MinifyGzipMiddleware minifies + gzips HTML responses.
func MinifyGzipMiddleware(next http.Handler) http.Handler {
	m := minify.New()
	m.AddFunc("text/html", minhtml.Minify)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		if broadcasting.IsWebSocketRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Only apply to HTML
		bw := &bufferedWriter{
			Middleware: core.Middleware{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			},
			buf: &bytes.Buffer{},
		}

		next.ServeHTTP(bw, r)

		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			// not HTML, send raw
			w.WriteHeader(bw.StatusCode)
			_, err := w.Write(bw.buf.Bytes())
			if err != nil {
				http.Error(w, "Minification header writing failed", http.StatusInternalServerError)
			}
			return
		}

		minified, err := m.Bytes("text/html", bw.buf.Bytes())
		if err != nil {
			http.Error(w, "Minification failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length") // gzip will change it
		w.WriteHeader(bw.StatusCode)

		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		_, err = gzw.Write(minified)
		if err != nil {
			http.Error(w, "gzip Writing Failed", http.StatusInternalServerError)
			return
		}
	})
}
