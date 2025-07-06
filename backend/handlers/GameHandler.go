package handlers

import (
	// "crypto/rand".
	// "encoding/hex".
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type GameHandler struct {
	l        *log.Logger
	viewsDir string
	games    map[string]Game
}

type Game struct {
	RoomID    string
	IsLocal   bool
	Players   []string
	CreatedAt int64
}

// const (
// 	idLength = 3
// )

func NewGameHandler(l *log.Logger, vd string) *GameHandler {
	return &GameHandler{
		l:        l,
		viewsDir: vd,
		games:    make(map[string]Game),
	}
}

func (gh *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	isFragment := r.Header.Get("X-Smart-Link") == "true"

	var tmpl *template.Template
	var err error
	if isFragment {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl, err = template.ParseFiles(
			filepath.Join(gh.viewsDir, "layouts", "empty.html"),
			filepath.Join(gh.viewsDir, "pages", "create.html"),
		)
	} else {
		tmpl, err = template.ParseFiles(
			filepath.Join(gh.viewsDir, "layouts", "base.html"),
			filepath.Join(gh.viewsDir, "pages", "create.html"),
		)
	}
	if err != nil {
		gh.l.Printf("Error parsing create templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	props := map[string]any{
		"title": "Spyfall Shqip - Krijo Lojë",
	}

	if err := tmpl.ExecuteTemplate(w, "base", props); err != nil {
		gh.l.Printf("Error executing create template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// generateRoomID creates a 6-character unique ID for game rooms.
// func generateRoomID() (string, error) {
// 	b := make([]byte, idLength) // 3 bytes = 6 hex characters
// 	if _, err := rand.Read(b); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(b), nil
// }
