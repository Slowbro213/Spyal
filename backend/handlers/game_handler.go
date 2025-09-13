package handlers

import (
	"bytes"
	"io"
	"time"

	"encoding/json"
	"net/http"
	"spyal/core"
	"spyal/dto/requests"
	"spyal/models"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"
	"spyal/pkg/utils/room"
	"spyal/repos"

	"go.uber.org/zap"
)

type GameHandler struct {
	core.Handler
	gameRepo  repos.GameRepository
	roundRepo repos.RoundRepository
	wordRepo  repos.WordRepository
}

func NewGameHandler(l *zap.Logger, g repos.GameRepository,
	r repos.RoundRepository, w repos.WordRepository) *GameHandler {
	return &GameHandler{
		Handler: core.Handler{
			Log: l,
		},
		gameRepo:  g,
		roundRepo: r,
		wordRepo:  w,
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
	var body requests.RemoteGameForm
	rawBody, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(rawBody))

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		gh.Log.Error("failed to decode JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := body.Validate(); err != nil {
		gh.Log.Info("invalid game form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomID, err := room.GenerateRoomID()
	if err != nil {
		gh.Log.Error("error generating room ID", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	raw := ctx.Value("id")


	userID, ok := raw.(int64)
	if !ok {
		gh.Log.Error("missing or malformed user id in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newGame := &models.Game{
		RoomID:     roomID,
		HostID:     userID,
		Name:       *body.GameName,
		SpyNumber:  body.SpyNumber,
		MaxPlayers: body.MaxNumbers,
		Private:    body.IsPrivate,
		CreatedAt:  time.Now().UTC(),
	}

	if err = gh.gameRepo.Create(ctx, newGame); err != nil {
		gh.Log.Error("failed to save game to db", zap.Error(err))
		http.Error(w, "Sorry, there was a problem creating your game", http.StatusInternalServerError)
		return
	}

	main, related, err := gh.wordRepo.RandomPair(ctx)

	if err != nil {
		gh.Log.Error("failed to retrieve words from db", zap.Error(err))
		http.Error(w, "Sorry, there was a problem creating your game", http.StatusInternalServerError)
		return
	}

	firstRound := &models.Round{
		GameID:  newGame.ID,
		Word:    main.ID,
		SpyWord: related.ID,
	}

	if err = gh.roundRepo.Create(ctx, firstRound); err != nil {
		gh.Log.Error("failed to save round to db", zap.Error(err))
		http.Error(w, "Sorry, there was a problem creating your game", http.StatusInternalServerError)
		return
	}

	if err = gh.roundRepo.AddParticipant(ctx, firstRound.ID, userID, false); err != nil {
		gh.Log.Error("failed to save round to db", zap.Error(err))
		http.Error(w, "Sorry, there was a problem creating your game", http.StatusInternalServerError)
		return
	}

	gh.Log.Info("stored game", zap.String("roomID", roomID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"roomID": roomID})
}
