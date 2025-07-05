package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"spyal/renderer"

	"github.com/joho/godotenv"
)

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

func main() {
	// ─── Load Environment ──────────────────────────────────────
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("❌ Error loading .env.%s: %v", env, err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")
	// pagesDir := os.Getenv("PAGES_DIR")

	logger := log.New(os.Stdout, "INFO ", log.LstdFlags)
	rh := renderer.NewRenderHandler(logger, viewsDir)

	// ─── Set Up Router ─────────────────────────────────────────
	mux := http.NewServeMux()

	// Static files: /public/*
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir))))

	// Pages: /
	mux.HandleFunc("/", rh.RenderPage)

	// Components: /views/*
	mux.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir(viewsDir))))

	// Render dynamic components: /components/*
	mux.HandleFunc("/components/", rh.RenderComponent)

	// Favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "favicon.ico"))
	})

	// Healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ , err := w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// ─── Start Server ──────────────────────────────────────────
	logger.Printf("✅ Server running at http://localhost:%s", port)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      mux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
