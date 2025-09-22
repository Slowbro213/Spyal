package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
	"net/http"
	"spyal/cache"
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

	isHost := game.HostID == userID

	props := map[string]any{
		"RoomID":      game.RoomID,
		"RoomName":    game.Name,
		"IsPublic":    !game.Private,
		"MaxPlayers":  game.MaxPlayers,
		"GameStarted": round.Status != "waiting",
		"Players":     users,
		"Spies":       spies,
		"UserID":      userID,
		"IsHost":      isHost,
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

// nolint
func (rh *RoomHandler) Leave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raw, ok := ctx.Value("id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := raw

	type LeaveResult struct {
		RoomID      string `db:"room_id"`
		HostID      int64  `db:"host_id"`
		RoundID     int64  `db:"round_id"`
		HostDeleted bool   `db:"host_deleted"`
	}

	q := `
	WITH target AS (
		SELECT r.id AS round_id, r.game_id, g.room_id, g.host_id
		FROM rounds r
		JOIN games g ON g.id = r.game_id
		JOIN game_participants gp ON gp.round_id = r.id
		WHERE gp.user_id = $1
		  AND r.status = 'waiting'
		ORDER BY r.created_at DESC
		LIMIT 1
	),
	deleted_participant AS (
		DELETE FROM game_participants gp
		WHERE gp.user_id = $1
		  AND gp.round_id IN (SELECT round_id FROM target)
		RETURNING round_id
	),
	deleted_round AS (
		DELETE FROM rounds r
		WHERE r.id IN (SELECT round_id FROM target WHERE host_id = $1)
		RETURNING id
	)
	SELECT t.room_id, t.host_id, t.round_id, dr.id IS NOT NULL AS host_deleted
	FROM target t
	LEFT JOIN deleted_round dr ON dr.id = t.round_id;
	`

	db := rh.roundRepo.DB()
	var result LeaveResult
	err := sqlx.GetContext(ctx, db, &result, q, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no waiting round found", http.StatusNotFound)
			return
		}
		rh.Log.Error("leave db error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	rh.Log.Info("roomID: " + result.RoomID)

	if err = cache.Delete(ctx, "players:"+result.RoomID); err != nil {
		rh.Log.Error("leave cache error", zap.Error(err))
	}

	if err = cache.Delete(ctx, "public_active_waiting"); err != nil {
		rh.Log.Error("leave cache error", zap.Error(err))
	}

	data := map[string]any{
		"user_id":  userID,
		"topic":    result.RoomID,
	}

	if result.HostDeleted {
		if err = cache.Delete(ctx, "round:"+result.RoomID); err != nil {
			rh.Log.Error("leave cache round error", zap.Error(err))
		}
		core.Dispatch(events.NewGameEndEvent(data))
	} else {
		core.Dispatch(events.NewLeftEvent(data))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":      true,
		"room_id":      result.RoomID,
		"host_left":    userID == result.HostID,
		"host_deleted": result.HostDeleted,
	})

}
