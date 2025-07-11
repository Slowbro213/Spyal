package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"spyal/handlers"
	"spyal/middleware"
	"spyal/pkg/utils/logger"
	"spyal/pkg/utils/metrics"
	"spyal/renderer"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

func loadEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env != "production" {
		err := godotenv.Load(".env." + env)
		if err != nil {
			log.Fatalf("❌ Error loading .env.%s: %v", env, err)
		}
	}
}

func main() {
	// ─── Load Environment ──────────────────────────────────────
	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	myLogger, err := logger.NewLogger()
	if err != nil {
		return
	}
	metrics := metrics.New()

	// Log startup
	myLogger.Info("Starting server on :8080")
	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")

	logger := log.New(os.Stdout, "INFO ", log.LstdFlags)
	rh := renderer.NewRenderHandler(logger, viewsDir)
	logger = log.New(os.Stdout, "GAMEHANDLER ", log.LstdFlags)
	gh := handlers.NewGameHandler(logger, viewsDir)

	// ─── Set Up Router ─────────────────────────────────────────
	mux := http.NewServeMux()

	mux.Handle("/public/", http.StripPrefix("/public/", middleware.BrotliStatic(publicDir)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Handle root page
		rh.RenderPage(w, r)
	})
	mux.HandleFunc("/create", gh.CreateGame)
	mux.HandleFunc("/create/remote", gh.CreateRemoteGame)

	mux.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir(viewsDir))))
	// Render dynamic components: /components/*
	mux.HandleFunc("/components/", rh.RenderComponent)

	// Favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "favicon.ico"))
	})

	// Healthcheck
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Metrics endpoint
	mux.Handle("/metrics", middleware.UsernamePassword("Thanas","Thanas24", promhttp.Handler(), *myLogger))

	// ─── Start Server ──────────────────────────────────────────
	logger.Printf("✅ Server running at http://localhost:%s", port)
	wrapped := middleware.MinifyGzipMiddleware(mux)
	wrapped = middleware.TrackMetrics(metrics,wrapped)
	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      wrapped,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
