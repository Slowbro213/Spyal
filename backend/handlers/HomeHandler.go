package handlers

import (
	"net/http"

	"spyal/core"
	"spyal/pkg/pages"
	"spyal/pkg/utils/template"

	"go.uber.org/zap"
)

type HomeHandler struct {
	core.Handler
	viewsDir string
}

func NewHomeHandler(l *zap.Logger, viewsDir string) *HomeHandler {
	return &HomeHandler{
		Handler: core.Handler{
			Log: l,
		},
		viewsDir: viewsDir,
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

	renderer := template.NewRenderer(hh.Log,hh.viewsDir)

	renderer.Render( w, isFragment, props,
		pages.LayoutBase,
		pages.PageHome,
		pages.CompRoom,
	)
}
