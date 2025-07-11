package handlers

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/services"
)

type itemHandler struct {
	services services.Services
	tmpl     *template.Template
}

func NewItemHandler(services services.Services, tmpl *template.Template) ItemHandler {
	return &itemHandler{
		services: services,
		tmpl:     tmpl,
	}
}

func (i *itemHandler) CreateItemForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	if err := i.tmpl.ExecuteTemplate(w, "add-item.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
