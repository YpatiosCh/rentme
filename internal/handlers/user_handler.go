package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/services"
)

type userHandler struct {
	service services.Services
	tmpl    *template.Template
}

func NewUserHandler(service services.Services, tmpl *template.Template) UserHandler {
	return &userHandler{
		service: service,
		tmpl:    tmpl,
	}
}

func (u *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := u.service.User().GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
