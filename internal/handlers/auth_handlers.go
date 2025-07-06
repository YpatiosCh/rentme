package handlers

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/services"
)

type authHandler struct {
	services services.Services
	tmpl     *template.Template
}

func NewAuthHandler(services services.Services, tmpl *template.Template) AuthHandler {
	return &authHandler{
		services: services,
		tmpl:     tmpl,
	}
}

// ShowRegistrationForm displays the multi-step registration form
func (h *authHandler) ShowRegistrationForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get available subscription plans for Step 1
	plans := config.GetSubscriptionPlans()

	data := struct {
		Plans map[string]config.PlanConfig
		Title string
	}{
		Plans: plans,
		Title: "Εγγραφή - RentMe",
	}

	if err := h.tmpl.ExecuteTemplate(w, "registration.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
