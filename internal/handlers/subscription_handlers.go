package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/services"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateCheckoutSessionRequest represents the request body for creating checkout session
type CreateCheckoutSessionRequest struct {
	PlanID       string  `json:"plan_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	BusinessName *string `json:"business_name"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	Region       string  `json:"region"`
	PostalCode   string  `json:"postal_code"`
}

// CreateCheckoutSessionHandler creates a Stripe checkout session
func (h *SubscriptionHandler) CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.PlanID == "" || req.FirstName == "" || req.LastName == "" || req.Phone == "" || req.Address == "" || req.City == "" || req.Region == "" || req.PostalCode == "" {
		http.Error(w, "All required fields must be filled", http.StatusBadRequest)
		return
	}

	// Validate plan exists
	plans := config.GetSubscriptionPlans()
	if _, exists := plans[req.PlanID]; !exists {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	// Create success and cancel URLs
	successURL := "http://localhost:8080/subscription/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := "http://localhost:8080/subscription/cancel"

	// Prepare metadata with all user data
	metadata := map[string]string{
		"plan_id":     req.PlanID,
		"first_name":  req.FirstName,
		"last_name":   req.LastName,
		"email":       req.Email,
		"phone":       req.Phone,
		"address":     req.Address,
		"city":        req.City,
		"region":      req.Region,
		"postal_code": req.PostalCode,
	}

	// Add business name if provided
	if req.BusinessName != nil && *req.BusinessName != "" {
		metadata["business_name"] = *req.BusinessName
	}

	// Create checkout session
	checkoutURL, err := h.subscriptionService.CreateCheckoutSessionWithUserData(
		req.Email,
		req.PlanID,
		successURL,
		cancelURL,
		metadata,
	)
	if err != nil {
		http.Error(w, "Failed to create checkout session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return checkout URL
	response := map[string]string{
		"checkout_url": checkoutURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// SubscriptionSuccessHandler handles successful subscription completion
func (h *SubscriptionHandler) SubscriptionSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Missing session_id parameter", http.StatusBadRequest)
		return
	}

	// Process the successful payment and create user
	user, err := h.subscriptionService.HandleSuccessfulPaymentFromSession(sessionID)
	if err != nil {
		http.Error(w, "Error processing payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Welcome to RentMe!</title>
		</head>
		<body>
			<h1>🎉 Welcome to RentMe, ` + user.FirstName + `!</h1>
			<p>Your ` + user.GetPlanDisplayName() + ` subscription has been successfully activated.</p>
			<p>Account created for: ` + user.Email + `</p>
			<p>You can now start adding your products for rent.</p>
			<a href="/dashboard">Go to Dashboard</a>
		</body>
		</html>
	`))
}

// SubscriptionCancelHandler handles subscription cancellation
func (h *SubscriptionHandler) SubscriptionCancelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Subscription Cancelled</title>
		</head>
		<body>
			<h1>Subscription Cancelled</h1>
			<p>You have cancelled your subscription process.</p>
			<p>You can try again anytime by returning to our plans page.</p>
			<a href="/">Try Again</a>
		</body>
		</html>
	`))
}
