package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/services"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
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

// CreateSubscription creates a Stripe Subscription for recurring payments
func (h *authHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlanID string `json:"plan_id"`
		Email  string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get plan details
	plan, exists := config.GetPlanByID(req.PlanID)
	if !exists {
		http.Error(w, "Invalid plan", http.StatusBadRequest)
		return
	}

	// Create Stripe customer
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
	}
	customer, err := customer.New(customerParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(plan.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
		Expand: []*string{
			stripe.String("latest_invoice.payment_intent"),
		},
	}

	subscription, err := subscription.New(subscriptionParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"clientSecret":   subscription.LatestInvoice.PaymentIntent.ClientSecret,
		"subscriptionId": subscription.ID,
		"customerId":     customer.ID,
	})
}

// GetStripeConfig returns the Stripe publishable key for the frontend
func (h *authHandler) GetStripeConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"publishable_key": config.GetPublishableKey(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		SubscriptionStatus: stringPtr("active"), // Payment succeeded
		SubscriptionID:     &req.SubscriptionID, // Stripe subscription ID
		CustomerID:         &req.CustomerID,     // Stripe customer ID
		PlanStartDate:      timePtr(time.Now()), // Start now
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
