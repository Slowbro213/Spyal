package renderer

import (
	"html/template"
	"net/http"
	"os"
	"spyal/pkg/pages"
)

func Render(w http.ResponseWriter, isFragment bool, props map[string]any, templates ...string) error {
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "development"
	}
	props["Stage"] = stage

	if isFragment {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates[0] = pages.LayoutEmpty
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "base", props); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return nil
}
