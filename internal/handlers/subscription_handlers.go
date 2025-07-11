package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/services"

	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/subscription"
)

type subHandler struct {
	services services.Services
	tmpl     *template.Template
}

func NewSubHandler(services services.Services, tmpl *template.Template) SubHandler {
	return &subHandler{
		services: services,
		tmpl:     tmpl,
	}
}

// CreateSubscription creates a Stripe Subscription for recurring payments
func (h *subHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
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
func (h *subHandler) GetStripeConfig(w http.ResponseWriter, r *http.Request) {
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

// StripeWebhook handles Stripe webhook events
func (h *subHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Verify webhook signature με ignore API version mismatch
	event, err := webhook.ConstructEventWithOptions(
		payload,
		r.Header.Get("Stripe-Signature"),
		config.GetWebhookSecret(),
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {
	case "invoice.payment_succeeded":
		h.handlePaymentSucceeded(event)
	case "invoice.payment_failed":
		h.handlePaymentFailed(event)
	case "customer.subscription.deleted":
		h.handleSubscriptionDeleted(event)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// handlePaymentSucceeded processes successful payment events
func (h *subHandler) handlePaymentSucceeded(event stripe.Event) {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Printf("Error parsing invoice: %v", err)
		return
	}

	log.Printf("Payment succeeded for customer: %s, amount: %d %s",
		invoice.Customer.ID, invoice.AmountPaid, invoice.Currency)

	// Find user by customer ID
	user, err := h.services.User().GetUserByCustomerID(invoice.Customer.ID)
	if err != nil {
		log.Printf("Error finding user by customer ID %s: %v", invoice.Customer.ID, err)
		return
	}

	// Update subscription status to active
	user.SubscriptionStatus = stringPtr("active")
	user.PlanStartDate = timePtr(time.Now())
	// Don't set PlanEndDate - subscription is recurring

	_, updateErr := h.services.User().UpdateUser(user)
	if updateErr != nil {
		log.Printf("Error updating user subscription status: %v", updateErr)
		return
	}

	log.Printf("Successfully updated user %s subscription status to active", user.Email)
}

// handlePaymentFailed processes failed payment events
func (h *subHandler) handlePaymentFailed(event stripe.Event) {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Printf("Error parsing invoice: %v", err)
		return
	}

	log.Printf("Payment failed for customer: %s, amount: %d %s",
		invoice.Customer.ID, invoice.AmountDue, invoice.Currency)

	// Find user by customer ID
	user, err := h.services.User().GetUserByCustomerID(invoice.Customer.ID)
	if err != nil {
		log.Printf("Error finding user by customer ID %s: %v", invoice.Customer.ID, err)
		return
	}

	// Update subscription status to past_due
	user.SubscriptionStatus = stringPtr("past_due")

	_, updateErr := h.services.User().UpdateUser(user)
	if updateErr != nil {
		log.Printf("Error updating user subscription status: %v", updateErr)
		return
	}

	log.Printf("Successfully updated user %s subscription status to past_due", user.Email)
}

// handleSubscriptionDeleted processes subscription cancellation events
func (h *subHandler) handleSubscriptionDeleted(event stripe.Event) {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		log.Printf("Error parsing subscription: %v", err)
		return
	}

	log.Printf("Subscription deleted for customer: %s", subscription.Customer.ID)

	// Find user by customer ID
	user, err := h.services.User().GetUserByCustomerID(subscription.Customer.ID)
	if err != nil {
		log.Printf("Error finding user by customer ID %s: %v", subscription.Customer.ID, err)
		return
	}

	// Update subscription status to canceled and set end date
	user.SubscriptionStatus = stringPtr("canceled")
	user.PlanEndDate = timePtr(time.Now())

	_, updateErr := h.services.User().UpdateUser(user)
	if updateErr != nil {
		log.Printf("Error updating user subscription status: %v", updateErr)
		return
	}

	log.Printf("Successfully updated user %s subscription status to canceled", user.Email)
}
