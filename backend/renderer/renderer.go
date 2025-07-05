package renderer

import (
	"html/template"
	"log"
	"net/http"
)

type RenderHandler struct {
	l *log.Logger
}

func NewRenderHandler(l *log.Logger) *RenderHandler {
	return &RenderHandler{l: l}
}

func (rh *RenderHandler) RenderComponent(w http.ResponseWriter, componentPath string, props map[string]any) {
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
