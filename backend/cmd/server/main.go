package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"spyal/core"
	"spyal/db"
	"spyal/cache"
	"spyal/pkg/pages"
	"spyal/pkg/utils/logger"
	"spyal/pkg/utils/metrics"

	"github.com/grafana/pyroscope-go"
	"github.com/joho/godotenv"
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

func inits() (*zap.Logger, *metrics.Metrics, *db.DB) {
	pyroscopeURL := os.Getenv("PYROSCOPE_URL")
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
	_, err = pyroscope.Start(pyroscope.Config{
		ApplicationName: "spyal",
		ServerAddress: pyroscopeURL,
		Logger: pyroscope.StandardLogger,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
		},
	})

	if err != nil {
		myLogger.Error("Error while Initing Pyroscope: ", zap.Error(err))
		log.Fatalf("Error loading Initing Pyroscope: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}


	database, err := db.Connect(db.Config{DSN: dsn})
	if err != nil {
		log.Fatal(err)
	}

	cache.Init()

	return myLogger, metrics, database
}

func startServer(handler http.Handler) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         host + ":" + port,
		Handler:      handler,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Printf("✅ Server running at http://%s:%s\n", host, port)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	loadEnv()

	myLogger, metrics, database := inits()
	defer database.Close()

	handler := Routes(myLogger, metrics,database)

	startServer(handler)
}
