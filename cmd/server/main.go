package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"spyal/renderer"
)

const (
	ReadTimeout = 5
	WriteTimeout = 10 
	IdleTimeout = 120
)

func main() {
	// Serve static assets
	fs := http.FileServer(http.Dir("../../assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	cs := http.FileServer(http.Dir("../../views/components/"))
	http.Handle("/components/", http.StripPrefix("/components/", cs))
	vs := http.FileServer(http.Dir("../../views/pages/"))
	http.Handle("/", http.StripPrefix("/", vs))


	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"../../assets/favicon.ico")
	})


	logger := log.New(os.Stdout, "INFO ", log.LstdFlags)
	rh := renderer.NewRenderHandler(logger)

	http.HandleFunc("/component/room", func(w http.ResponseWriter, _ *http.Request) {
		props := map[string]any{
			"RoomID":     "ABC123",
			"PlayerName": "Ardi",
		}
		rh.RenderComponent(w, "../../views/components/room.html", props)
	})

	log.Println("Server running on http://localhost:8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  ReadTimeout * time.Second,
		WriteTimeout: WriteTimeout * time.Second,
		IdleTimeout:  IdleTimeout * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
