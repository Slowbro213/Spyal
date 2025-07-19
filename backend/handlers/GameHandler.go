package handlers

import (
	"context"
	"time"

	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"spyal/core"
	"spyal/db"
	"spyal/pkg/pages"
	"spyal/pkg/game"
	"spyal/pkg/utils"
	"spyal/pkg/utils/renderer"
)

type GameHandler struct {
	core.Handler
}


func NewGameHandler(l *zap.Logger) *GameHandler {
	return &GameHandler{
		Handler: core.Handler{
			Log: l,
		},
	}
}

func (gh *GameHandler) CreateGamePage(w http.ResponseWriter, r *http.Request) {
	isFragment := pages.IsFragment(r)
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë"}
	err := renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageCreate,
	)
	if err != nil {
		gh.Log.Error("Error while rendering Create Game Page: ", zap.Error(err))
	}
}

func (gh *GameHandler) CreateRemoteGamePage(w http.ResponseWriter, r *http.Request) {
	isFragment := pages.IsFragment(r)
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë Online"}
	err := renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageRemote,
	)
	if err != nil {
		gh.Log.Error("Error while rendering Create Remote Game Page: ", zap.Error(err))
	}
}

func (gh *GameHandler) CreateRemoteGame(w http.ResponseWriter, r *http.Request) {
	var body game.GameForm

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		gh.Log.Error("failed to decode JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	roomID, err := utils.GenerateRoomID()
	if err != nil {
		gh.Log.Error("Error while generating Room ID: ", zap.Error(err))
	}

	game := game.Game{
		RoomID:      roomID,
		Players:     []string{body.Params.PlayerName},
		RoomName:    body.Params.GameName,
		IsPublic:    !body.Params.IsPrivate,
		MaxPlayers:  body.Params.MaxNumbers,
		GameStarted: false,
		CreatedAt:   time.Now().Unix(),
	}

	gameJSON, err := json.Marshal(game)
	if err != nil {
		gh.Log.Error("failed to marshal game to JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	client, err := db.GetRedis()
	if err != nil {
		gh.Log.Error("failed to connect to Redis", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	err = client.Set(ctx, "game:"+roomID, gameJSON, 0).Err()
	if err != nil {
		gh.Log.Error("failed to store game in Redis", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	gh.Log.Info("stored game in Redis", zap.String("roomID", roomID))

	resp := map[string]any{"roomID": roomID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		gh.Log.Error("Error while Encoding Json: ", zap.Error(err))
	}
}
