package handlers

import (
	"context"
	"strings"

	"encoding/json"
	"net/http"
	"spyal/cache"
	"spyal/core"
	"spyal/pkg/game"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"

	"go.uber.org/zap"
)

type RoomHandler struct {
	core.Handler
}


func NewRoomHandler(l *zap.Logger) *RoomHandler {
	return &RoomHandler{
		Handler: core.Handler{
			Log: l,
		},
	}
}

func (rh *RoomHandler) Show(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path                  
	parts := strings.Split(path, "/")  

	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	roomID := parts[2]

	ctx := context.Background()
	gameJSON, err := cache.Get(ctx, "game:"+roomID)
	if err != nil {
		rh.Log.Error("failed to get game from Memory", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var game game.Game
	if err := json.Unmarshal([]byte(gameJSON), &game); err != nil {
		rh.Log.Error("failed to unmarshal game JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	props := map[string]any{
		"RoomID":      game.RoomID,
		"RoomName":    game.RoomName,
		"IsPublic":    game.IsPublic,
		"MaxPlayers":  game.MaxPlayers,
		"GameStarted": game.GameStarted,
		"Players":     game.Players,
		"Spies":       game.Spies,
		"IsHost":      true,
		"CreatedAt":   game.CreatedAt,
	}

	isFragment := pages.IsFragment(r)
	err = renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageRoom,
	)

	if err != nil {
		rh.Log.Error("Error Rendering Room: ", zap.Error(err))
	}
}
