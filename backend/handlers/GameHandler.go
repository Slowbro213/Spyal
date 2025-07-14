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
	"spyal/pkg/utils"
	"spyal/pkg/utils/template"
)

type GameHandler struct {
	core.Handler
	viewsDir string
}

type GameForm struct {
	Params struct {
		PlayerName string `json:"playerName"`
		GameName   string `json:"gameName"`
		Spies      int  	`json:"spyNumber"`
		MaxNumbers int    `json:"maxNumbers"`
		IsPrivate  bool   `json:"isPrivate"`
	} `json:"params"`
}

type Round struct {
	Number     int      `json:"number"`
	Word       string   `json:"word"`
	SpyWord    string   `json:"spyWord"`
	Spies      []string `json:"spies"` 
	Winner     string   `json:"winner"`  
}



type Game struct {
	RoomID      string   `json:"roomID"`
	Players     []string `json:"players"`
	Word        string   `json:"word"`
	SpyWord     string   `json:"spyWord"`
	Spies       []string `json:"spies"`
	RoomName    string   `json:"roomName"`
	IsPublic    bool     `json:"isPublic"`
	MaxPlayers  int      `json:"maxPlayers"`
	GameStarted bool     `json:"gameStarted"`
	Rounds      []Round  `json:"rounds"`
	CreatedAt   int64    `json:"createdAt"`
}


func NewGameHandler(l *zap.Logger, vd string) *GameHandler {
	return &GameHandler{
		Handler: core.Handler{
			Log: l,
		},
		viewsDir: vd,
	}
}

func (gh *GameHandler) CreateGamePage(w http.ResponseWriter, r *http.Request) {
	isFragment := pages.IsFragment(r)
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë"}
	renderer := template.NewRenderer(gh.Log, gh.viewsDir)
	renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageCreate,
	)
}

func (gh *GameHandler) CreateRemoteGamePage(w http.ResponseWriter, r *http.Request) {
	isFragment := pages.IsFragment(r)
	props := map[string]any{"title": "Spyfall Shqip - Krijo Lojë Online"}
	renderer := template.NewRenderer(gh.Log, gh.viewsDir)
	renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageRemote,
	)
}

func (gh *GameHandler) CreateRemoteGame(w http.ResponseWriter, r *http.Request) {
	var body GameForm

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

	game := Game{
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
