package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"spyal/bootstrap"
	"spyal/core"
	"spyal/events"
	"spyal/handlers"
	"spyal/middleware"
	"spyal/pkg/pages"
	"spyal/pkg/utils/logger"
	"spyal/pkg/utils/metrics"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
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

func inits() (*zap.Logger, *metrics.Metrics) {
	myLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Error loading logger: %v", err)
	}
	core.Logger = myLogger
	metrics := metrics.New()
	err = pages.Init()
	if err != nil {
		myLogger.Error("Error while Initing Pages: ", zap.Error(err))
		log.Fatalf("Error loading Initing Pages: %v", err)
	}
	return myLogger, metrics
}

//nolint
func setupRouter(myLogger *zap.Logger, metrics *metrics.Metrics) http.Handler {
	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	hh := handlers.NewHomeHandler(myLogger)
	gh := handlers.NewGameHandler(myLogger)
	rh := handlers.NewRoomHandler(myLogger)
	lh := handlers.NewLogHandler(myLogger)

	router := core.NewRouter()

	router.Get("/public/", http.StripPrefix("/public/", middleware.BrotliStatic(publicDir)).ServeHTTP)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		hh.HomePage(w, r)
	})

	router.Get("/create", gh.CreateGamePage)
	router.Get("/create/remote", gh.CreateRemoteGamePage)

	router.Post("/create/remote", gh.CreateRemoteGame)

	router.Get("/room/", rh.Show)

	router.Get("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir(viewsDir))).ServeHTTP)

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "favicon.ico"))
	})

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	router.Get("/metrics", middleware.UsernamePassword(username, password, promhttp.Handler(), *myLogger).ServeHTTP)

	router.Post("/api/log", lh.LogFrontend)

	eServer := core.PokedServer{
		Log: myLogger,
		EmitEvent: events.NewEchoEvent,
	}

	err := bootstrap.InitAll()
	if err != nil {
		myLogger.Error("Error initializing " , zap.Error(err))
		log.Fatalf("Error initializing %v",err)
	}

	router.Get("/echo", eServer.StartWSServer)

	handler := middleware.MinifyGzipMiddleware(router)
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

	myLogger, metrics := inits()

	handler := setupRouter(myLogger, metrics)

	startServer(handler)
}
