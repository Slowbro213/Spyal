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
	"spyal/db"
	"spyal/handlers"
	"spyal/middleware"
	"spyal/pkg/utils/metrics"
	"spyal/repos"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

//nolint
func Routes(myLogger *zap.Logger, metrics *metrics.Metrics, database *db.DB) http.Handler {
	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	userRepo := repos.NewUserRepo(database)
	gameRepo := repos.NewGameRepo(database)
	roundRepo := repos.NewRoundRepo(database)
	wordRepo := repos.NewWordRepo(database)
	hh := handlers.NewHomeHandler(myLogger,gameRepo)
	gh := handlers.NewGameHandler(myLogger,gameRepo,roundRepo,wordRepo)
	rh := handlers.NewRoomHandler(myLogger,gameRepo,roundRepo,userRepo)
	lh := handlers.NewLogHandler(myLogger)

	uh := handlers.NewUserHandler(myLogger,userRepo)

	router := core.NewRouter()

	router.Get("/public/", http.StripPrefix("/public/", middleware.BrotliStatic(publicDir)).ServeHTTP)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		hh.HomePage(w, r)
	})

	router.Get("/login", middleware.GuestMiddleware(http.HandlerFunc(uh.LoginPage)).ServeHTTP)
	router.Post("/login", middleware.GuestMiddleware(http.HandlerFunc(uh.LoginOrRegister)).ServeHTTP)
	router.Post("/logout", middleware.GuestMiddleware(http.HandlerFunc(uh.Logout)).ServeHTTP)

	router.Get("/create", gh.CreateGamePage)
	router.Get("/create/remote", gh.CreateRemoteGamePage)
	router.Get("/games", gh.Index)

	router.Post("/create/remote", middleware.AuthMiddleware(http.HandlerFunc(gh.CreateRemoteGame)).ServeHTTP)

	router.Get("/room/", middleware.AuthMiddleware(http.HandlerFunc(rh.Show)).ServeHTTP)
	router.Post("/leave", middleware.AuthMiddleware(http.HandlerFunc(rh.Leave)).ServeHTTP)

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

	router.Get("/poked", eServer.StartWSServer)

	handler := middleware.MinifyGzipMiddleware(router)
	handler = middleware.TrackMetrics(metrics, handler)
	handler = middleware.RateLimitMiddleware(handler)

	return handler
}
