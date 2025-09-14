package handlers

import (
	"net/http"

	"spyal/core"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"
	"spyal/repos"

	"go.uber.org/zap"
)

type HomeHandler struct {
	core.Handler
	gameRepo repos.GameRepository
}

func NewHomeHandler(l *zap.Logger, g repos.GameRepository) *HomeHandler {
	return &HomeHandler{
		Handler: core.Handler{
			Log: l,
		},
		gameRepo: g,
	}
}

func (hh *HomeHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	games, err := hh.gameRepo.GetPublicActive(ctx, "")
	if err != nil {
		hh.Log.Error("Error while fetching public active games: ", zap.Error(err))
		if pages.IsFragment(r) {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("Kishte një problem në marrjen e lojrave."))
			if err != nil {
				http.Error(w, "Kishte një problem në marrjen e lojrave.", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Kishte një problem në marrjen e lojrave.", http.StatusInternalServerError)
		}
		return
	}

	roomsCap := 3
	props := map[string]any{
		"Rooms":    games,
		"RoomsCap": roomsCap,
	}

	isFragment := pages.IsFragment(r)
	err = renderer.Render(w, isFragment, props,
		pages.LayoutBase,
		pages.PageHome,
	)
	if err != nil {
		hh.Log.Error("Error rendering Page: ", zap.Error(err))
		if !pages.IsFragment(r) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
