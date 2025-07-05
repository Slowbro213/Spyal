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
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	_ = godotenv.Load(".env." + env) // Loads .env.development or .env.production
	port := os.Getenv("PORT")

	// Resolve paths based on environment
	publicDir := os.Getenv("PUBLIC_DIR")
	viewsDir := os.Getenv("VIEWS_DIR")
	pagesDir := os.Getenv("PAGES_DIR")


	logger := log.New(os.Stdout, "INFO ", log.LstdFlags)
	rh := renderer.NewRenderHandler(logger)

	// Serve static assets at /static/*
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir))))
	// Serve Pages 
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(pagesDir))))

	// Favicon route
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "favicon.ico"))
	})

	// Component route example
	http.HandleFunc("/component/room", func(w http.ResponseWriter, _ *http.Request) {
		props := map[string]any{
			"RoomID":     "ABC123",
			"PlayerName": "Ardi",
		}
		componentPath := filepath.Join(viewsDir, "components", "room.html")
		rh.RenderComponent(w, componentPath, props)
	})

	log.Println("✅ Server running at http://localhost:8080")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
