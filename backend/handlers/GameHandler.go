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

func (gh *GameHandler) CreateGame(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles(
		filepath.Join(gh.viewsDir, "layouts", "base.html"),
		filepath.Join(gh.viewsDir, "pages", "create.html"),
	)
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
