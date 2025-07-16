package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/middleware"
	"github.com/YpatiosCh/rentme/internal/models"
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

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		user, _ := middleware.GetUserFromContext(r)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := h.tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		plainPassword := r.FormValue("password")

		_, token, err := h.services.Auth().LoginUser(email, plainPassword)
		if err != nil {
			http.Error(w, err.Message, err.Code)
			return
		}

		// Set JWT cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			HttpOnly: true,
			Secure:   false, // true σε production
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   7 * 24 * 3600, // 7 days
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "No session to logout", http.StatusUnauthorized)
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CompleteRegistration completes user registration after successful payment
func (h *authHandler) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PlanID          string `json:"plan_id"`
		PaymentIntentID string `json:"payment_intent_id"`
		SubscriptionID  string `json:"subscription_id"`
		CustomerID      string `json:"customer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create user with full subscription info
	user := &models.User{
		Email: req.Email,
		// Subscription fields
		SubscriptionStatus: stringPtr("active"),                  // Payment succeeded
		SubscriptionID:     &req.SubscriptionID,                  // Stripe subscription ID
		CustomerID:         &req.CustomerID,                      // Stripe customer ID
		PlanStartDate:      timePtr(time.Now()),                  // Start now
		NextBillingDate:    timePtr(time.Now().AddDate(0, 1, 0)), // One month from now
	}

	// Set plan limits
	user.SetSubscriptionLimits(req.PlanID)

	// Register user
	createdUser, err := h.services.Auth().RegisterUser(user, req.Password)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}

	// Generate JWT
	token, err := h.services.Auth().GenerateToken(createdUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // true σε production
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7 days
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Registration completed successfully",
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
