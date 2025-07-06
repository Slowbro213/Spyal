package renderer

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type RenderHandler struct {
	l        *log.Logger
	viewsDir string
}

func NewRenderHandler(l *log.Logger, vd string) *RenderHandler {
	return &RenderHandler{
		l:        l,
		viewsDir: vd,
	}
}

func (rh *RenderHandler) RenderPage(w http.ResponseWriter, r *http.Request) {
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

	isFragment := r.Header.Get("X-Smart-Link") == "true"

	var layout string

	if isFragment {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		layout = filepath.Join(rh.viewsDir, "layouts", "empty.html")
	} else {
		layout = filepath.Join(rh.viewsDir, "layouts", "base.html")
	}

	tmpl, err := template.ParseFiles(
		layout,
		filepath.Join(rh.viewsDir, "pages", "index.html"),
		filepath.Join(rh.viewsDir, "components", "button.html"),
		filepath.Join(rh.viewsDir, "components", "room.html"),
	)
	if err != nil {
		rh.l.Printf("Template parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ✅ Render the "base" layout template, which should include others via {{ template }}
	if err := tmpl.ExecuteTemplate(w, "base", props); err != nil {
		rh.l.Printf("Template exec error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (rh *RenderHandler) RenderComponent(w http.ResponseWriter, r *http.Request) {
	componentPath := r.URL.Path // e.g. "/component/room" → "./component/room"

	// Parse query parameters as props
	props := make(map[string]any)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			props[key] = values[0] // Only take the first value
		}
	}

	rh.render(w, rh.viewsDir+componentPath+".html", props) // Add .html if needed}
}

func (rh *RenderHandler) render(w http.ResponseWriter, componentPath string, props map[string]any) {
	tmpl, err := template.ParseFiles(componentPath)
	if err != nil {
		rh.l.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, props)
	if err != nil {
		rh.l.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
