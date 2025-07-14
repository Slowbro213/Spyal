package template

import (
	"html/template"
	"net/http"
	"os"
	"spyal/pkg/pages"

	"path/filepath"

	"go.uber.org/zap"
)

type Renderer struct {
	logger    *zap.Logger
	viewsDir  string
}


func NewRenderer(l *zap.Logger, viewsDir string) *Renderer{
	return &Renderer{
		logger: l,
		viewsDir: viewsDir,
	}
}

func (tr *Renderer) Render(w http.ResponseWriter, isFragment bool, props map[string]any, templates ...string) {
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "development"
	}
	props["Stage"] = stage

	if isFragment {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates[0] = pages.LayoutEmpty
	}

	var files []string
	for _, tpl := range templates {
		files = append(files, filepath.Join(tr.viewsDir, tpl))
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		tr.logger.Error("Renderer: template parse error", zap.Error(err), zap.Strings("files", files))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", props); err != nil {
		tr.logger.Error("Renderer: template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
