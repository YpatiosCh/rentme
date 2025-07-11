package handlers

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/middleware"
	"github.com/YpatiosCh/rentme/internal/models"
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

	user, _ := middleware.GetUserFromContext(r)

	data := struct {
		User *models.User
	}{
		User: user,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
