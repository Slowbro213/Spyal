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
	"go.uber.org/zap"
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
		if err := godotenv.Load(".env." + env); err != nil {
			log.Fatalf("❌ Error loading .env.%s: %v", env, err)
		}
	}
}

func initLoggerAndMetrics() (*zap.Logger, *metrics.Metrics) {
	myLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Error loading logger: %v", err)
	}
	metrics := metrics.New()
	return myLogger, metrics
}

func setupRouter(myLogger *zap.Logger, metrics *metrics.Metrics) http.Handler {
	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	rh := renderer.NewRenderHandler(myLogger, viewsDir)
	gh := handlers.NewGameHandler(myLogger, viewsDir)
	lh := handlers.NewLogHandler(myLogger)

	mux := http.NewServeMux()

	mux.Handle("/public/", http.StripPrefix("/public/", middleware.BrotliStatic(publicDir)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		rh.RenderPage(w, r)
	})

	mux.HandleFunc("/create", gh.CreateGame)
	mux.HandleFunc("/create/remote", gh.CreateRemoteGame)

	mux.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir(viewsDir))))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "favicon.ico"))
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	mux.Handle("/metrics", middleware.UsernamePassword(username, password, promhttp.Handler(), *myLogger))

	mux.HandleFunc("/api/log", lh.LogFrontend)

	handler := middleware.MinifyGzipMiddleware(mux)
	handler = middleware.TrackMetrics(metrics, handler)
	handler = middleware.RateLimitMiddleware(handler)

	return handler
}

func startServer(handler http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      handler,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Printf("✅ Server running at http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	loadEnv()

	myLogger, metrics := initLoggerAndMetrics()

	handler := setupRouter(myLogger, metrics)

	startServer(handler)
}
