package handlers

import (
	"html/template"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"go.uber.org/zap"
)

type GameHandler struct {
	l        *zap.Logger
	viewsDir string
	games    map[string]Game
}

type Game struct {
	RoomID    string
	IsLocal   bool
	Players   []string
	CreatedAt int64
}

// Layout constants.
const (
	LayoutBase   = "base.html"
	LayoutEmpty  = "empty.html"
	PageCreate   = "create.html"
	PageOnline   = "remote.html"
)

// NewGameHandler creates a new GameHandler.
func NewGameHandler(l *zap.Logger, vd string) *GameHandler {
	return &GameHandler{
		l:        l,
		viewsDir: vd,
		games:    make(map[string]Game),
	}
}

// CreateGame handles creation of local games (renders the create page).
func (gh *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	isFragment := r.Header.Get("X-Smart-Link") == "true"
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë"}
	gh.renderTemplate(w, LayoutBase, PageCreate, isFragment, props)
}

// CreateRemoteGame handles creation of remote games (renders the online page).
func (gh *GameHandler) CreateRemoteGame(w http.ResponseWriter, r *http.Request) {
	isFragment := r.Header.Get("X-Smart-Link") == "true"
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë Online"}
	gh.renderTemplate(w, LayoutBase, PageOnline, isFragment, props)
}
// renderTemplate renders a template with the given layout, page, and props.
// If isFragment is true, uses the empty layout.
func (gh *GameHandler) renderTemplate(w http.ResponseWriter, layout, page string, isFragment bool, props map[string]any) {
	var tmpl *template.Template
	var err error

	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "development" // fallback default.
	}
	props["Stage"] = stage

	layoutName := layout
	if isFragment {
		layoutName = LayoutEmpty
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	tmpl, err = template.ParseFiles(
		filepath.Join(gh.viewsDir, "layouts", layoutName),
		filepath.Join(gh.viewsDir, "pages", page),
	)
	if err != nil {
		gh.l.Error(fmt.Sprintf("Error parsing templates: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", props); err != nil {
		gh.l.Error(fmt.Sprintf("Error executing template: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
