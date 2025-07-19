package handlers

import (
	"net/http"

	"spyal/core"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"

	"go.uber.org/zap"
)

type HomeHandler struct {
	core.Handler
}

func NewHomeHandler(l *zap.Logger) *HomeHandler {
	return &HomeHandler{
		Handler: core.Handler{
			Log: l,
		},
	}
}

func (hh *HomeHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	props := map[string]any{
		"title": "Spyfall Shqip",
		"Room": map[string]any{
			"roomid":     "ABC123",
			"playername": "Thanas Papa",
		},
		"Button": map[string]any{
			"text":  "Krijo Dhome",
			"href":  "/create",
			"class": "btn btn-primary",
			"icon":  "➕",
		},
	}

	isFragment := pages.IsFragment(r)

	err := renderer.Render( w, isFragment, props,
		pages.LayoutBase,
		pages.PageHome,
		pages.CompRoom,
	)

	if err != nil {
		hh.Log.Error("Error while rendering Page: ", zap.Error(err))
	}
}
