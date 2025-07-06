package services

import (
	"fmt"
	"time"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/repositories"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/subscription"
)

type subscriptionService struct {
	userRepo repositories.UserRepository
}

func NewSubscriptionService(userRepo repositories.UserRepository) SubscriptionService {
	return &subscriptionService{
		userRepo: userRepo,
	}
}

// GetCustomerSubscriptions gets all subscriptions for a Stripe customer
func (s *subscriptionService) GetCustomerSubscriptions(customerID string) ([]*stripe.Subscription, error) {
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
	}

	iter := subscription.List(params)
	var subscriptions []*stripe.Subscription

	for iter.Next() {
		subscriptions = append(subscriptions, iter.Subscription())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subscriptions, nil
}

// CancelSubscription cancels a user's subscription
func (s *subscriptionService) CancelSubscription(userID string) error {
	// Get user from database
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.SubscriptionID == nil {
		return fmt.Errorf("user has no active subscription")
	}

	// Cancel subscription in Stripe
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err = subscription.Update(*user.SubscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription in Stripe: %w", err)
	}

	// Update user status
	canceledStatus := "canceled"
	user.SubscriptionStatus = &canceledStatus

	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateSubscriptionPlan changes a user's subscription plan
func (s *subscriptionService) UpdateSubscriptionPlan(userID, newPlanID string) error {
	// Get user from database
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.SubscriptionID == nil {
		return fmt.Errorf("user has no active subscription")
	}

	// Get new plan configuration
	plans := config.GetSubscriptionPlans()
	newPlan, exists := plans[newPlanID]
	if !exists {
		return fmt.Errorf("invalid plan: %s", newPlanID)
	}

	// Get current subscription from Stripe
	stripeSubscription, err := subscription.Get(*user.SubscriptionID, nil)
	if err != nil {
		return fmt.Errorf("failed to get subscription from Stripe: %w", err)
	}

	// Update subscription in Stripe
	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(stripeSubscription.Items.Data[0].ID),
				Price: stripe.String(newPlan.PriceID),
			},
		},
		ProrationBehavior: stripe.String("create_prorations"),
	}

	_, err = subscription.Update(*user.SubscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to update subscription in Stripe: %w", err)
	}

	// Update user in database
	user.SubscriptionPlan = &newPlanID
	user.MaxProducts = newPlan.MaxProducts
	user.MaxPhotosPerProduct = newPlan.MaxPhotos
	user.FeaturedListingsLimit = newPlan.FeaturedListings

	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// CreateCheckoutSessionWithUserData creates a Stripe checkout session with user data in metadata
func (s *subscriptionService) CreateCheckoutSessionWithUserData(userEmail, planID, successURL, cancelURL string, userData map[string]string) (string, error) {
	// Get plan configuration
	plans := config.GetSubscriptionPlans()
	plan, exists := plans[planID]
	if !exists {
		return "", fmt.Errorf("invalid plan: %s", planID)
	}

	// Create checkout session parameters
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		CustomerEmail:      stripe.String(userEmail),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(plan.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: userData,
		},
		Metadata: userData, // Also store in session metadata
	}

	// Create the session
	checkoutSession, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return checkoutSession.URL, nil
}

// HandleSuccessfulPaymentFromSession processes payment using session ID
func (s *subscriptionService) HandleSuccessfulPaymentFromSession(sessionID string) (*models.User, error) {
	// Get session details from Stripe
	stripeSession, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get session from Stripe: %w", err)
	}

	// Extract user data from metadata
	metadata := stripeSession.Metadata

	// Get plan configuration
	planID := metadata["plan_id"]
	plans := config.GetSubscriptionPlans()
	plan, exists := plans[planID]
	if !exists {
		return nil, fmt.Errorf("invalid plan: %s", planID)
	}

	// Create new user with all data
	user := &models.User{
		Email:        metadata["email"],
		PasswordHash: "", // Will be set later when user sets password

		// Personal Information
		FirstName: metadata["first_name"],
		LastName:  metadata["last_name"],
		Phone:     metadata["phone"],

		// Address Information
		Address:    metadata["address"],
		City:       metadata["city"],
		Region:     metadata["region"],
		PostalCode: metadata["postal_code"],
		Country:    "Greece",

		// Subscription data
		SubscriptionPlan:   &planID,
		SubscriptionStatus: &[]string{"active"}[0],
		SubscriptionID:     &stripeSession.Subscription.ID,
		CustomerID:         &stripeSession.Customer.ID,

		// Plan limits
		MaxProducts:           plan.MaxProducts,
		MaxPhotosPerProduct:   plan.MaxPhotos,
		FeaturedListingsLimit: plan.FeaturedListings,
		FeaturedListingsUsed:  0,

		// Account status
		IsActive:      true,
		EmailVerified: true, // Auto-verify since they paid

		// Dates
		PlanStartDate: &time.Time{},
		PlanEndDate:   &time.Time{},
	}

	// Add business name if provided
	if businessName, exists := metadata["business_name"]; exists && businessName != "" {
		user.BusinessName = &businessName
	}

	// Set plan dates from subscription
	if stripeSession.Subscription != nil {
		planStart := time.Unix(stripeSession.Subscription.CurrentPeriodStart, 0)
		planEnd := time.Unix(stripeSession.Subscription.CurrentPeriodEnd, 0)
		user.PlanStartDate = &planStart
		user.PlanEndDate = &planEnd
	}

	// Save user to database
	return s.userRepo.CreateUser(user)
}
