package handlers

import (
	"fmt"
	"slices"
	"strings"

	"net/http"
	"spyal/core"
	"spyal/events"
	"spyal/models"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"
	"spyal/repos"

	"go.uber.org/zap"
)

type RoomHandler struct {
	core.Handler
	gameRepo  repos.GameRepository
	roundRepo repos.RoundRepository
	userRepo  repos.UserRepository
}

func NewRoomHandler(l *zap.Logger, g repos.GameRepository, r repos.RoundRepository, u repos.UserRepository) *RoomHandler {
	return &RoomHandler{
		Handler: core.Handler{
			Log: l,
		},
		gameRepo:  g,
		roundRepo: r,
		userRepo:  u,
	}
}

func (rh *RoomHandler) Show(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[2] == "" {
		rh.Log.Error(fmt.Sprintf("Room ID problem. parts : %v", parts))
		http.Error(w, "Room ID is required", http.StatusInternalServerError)
		return
	}
	roomID := parts[2]

	ctx := r.Context()

	var game models.Game
	if err := rh.gameRepo.GetBy(ctx, &game, "room_id", roomID); err != nil {
		rh.Log.Error("Error retrieving game from room_id: ", zap.Error(err))
		http.Error(w, "There was a problem retrieving your game", http.StatusInternalServerError)
		return
	}

	var round models.Round

	if err := rh.roundRepo.GetBy(ctx, &round, "game_id", game.ID); err != nil {
		rh.Log.Error("Error fetching round: ", zap.Error(err))
		http.Error(w, "There was a problem retrieving your game", http.StatusInternalServerError)
		return
	}

	players, err := rh.roundRepo.GetPlayers(ctx, round.ID)

	if err != nil {
		rh.Log.Error("Error Fetching players: ", zap.Error(err))
		http.Error(w, "There was a problem retrieving your game", http.StatusInternalServerError)
		return
	}

	users := make(map[int]*models.User)
	for _, player := range players {
		var user models.User
		if err = rh.userRepo.GetBy(ctx, &user, "id", player.UserID); err != nil {
			rh.Log.Error("Error fetching user: ", zap.Error(err))
			http.Error(w, "There was a problem retrieving your game", http.StatusInternalServerError)
			return
		}
		users[int(player.UserID)] = &user
	}

	raw, ok := ctx.Value("id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := raw

	exists := false
	for _, u := range users {
		if u.ID == userID {
			exists = true
			break
		}
	}

	rawName, ok := ctx.Value("username").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	username := rawName
	data := map[string]any{
		"user_id":  userID,
		"username": username,
		"topic":    roomID,
	}
	if !exists {
		core.Dispatch(events.NewUserJoinedEvent(data))
		if err := rh.roundRepo.AddParticipant(ctx, round.ID, userID, false); err != nil {
			rh.Log.Error("Error adding participant", zap.Error(err))
			http.Error(w, "Kishte nje problem ne shtimin tuaj ne loj", http.StatusInternalServerError)
			return
		}
		var user models.User
		if err := rh.userRepo.GetBy(ctx, &user, "id", userID); err != nil {
			rh.Log.Error("Error retrieving user", zap.Error(err))
			http.Error(w, "Kishte nje problem ne shtimin tuaj ne loj", http.StatusInternalServerError)
			return
		}
		users[int(userID)] = &user
	}

	spies := slices.Collect(
		func(yield func(int64) bool) {
			for _, p := range players {
				if p.IsSpy {
					if !yield(p.UserID) {
						return
					}
				}
			}
		},
	)

	props := map[string]any{
		"RoomID":      game.RoomID,
		"RoomName":    game.Name,
		"IsPublic":    !game.Private,
		"MaxPlayers":  game.MaxPlayers,
		"GameStarted": round.Status != "waiting",
		"Players":     users,
		"Spies":       spies,
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
