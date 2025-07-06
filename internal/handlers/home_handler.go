package handlers

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/services"
)

type homeHandler struct {
	service services.Services
	tmpl    *template.Template
}

func NewHomeHandler(services services.Services, tmpl *template.Template) HomeHandler {
	return &homeHandler{
		service: services,
		tmpl:    tmpl,
	}
}

func (h *homeHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
