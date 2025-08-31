package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"spyal/bootstrap"
	"spyal/broadcasting"
	"spyal/core"
	"spyal/handlers"
	"spyal/middleware"
	"spyal/pkg/utils/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func Routes(myLogger *zap.Logger, metrics *metrics.Metrics) http.Handler {
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

	eServer := broadcasting.PokedServer{
		Log: myLogger,
	}

	err := bootstrap.InitAll()
	if err != nil {
		myLogger.Error("Error initializing ", zap.Error(err))
		log.Fatalf("Error initializing %v", err)
	}

	router.Get("/echo", eServer.StartWSServer)

	handler := middleware.MinifyGzipMiddleware(router)
	handler = middleware.TrackMetrics(metrics, handler)
	handler = middleware.RateLimitMiddleware(handler)

	return handler
}
