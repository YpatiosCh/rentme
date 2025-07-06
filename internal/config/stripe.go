package config

import (
	"os"

	"github.com/stripe/stripe-go/v76"
)

// InitStripe initializes Stripe with the secret key
func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// GetPublishableKey returns the publishable key for frontend
func GetPublishableKey() string {
	return os.Getenv("STRIPE_PUBLISHABLE_KEY")
}

// GetWebhookSecret returns the webhook secret for signature verification
func GetWebhookSecret() string {
	return os.Getenv("STRIPE_WEBHOOK_SECRET")
}

// Subscription plans configuration
type PlanConfig struct {
	ID               string
	Name             string
	PriceID          string  // Stripe Price ID
	Price            float64 // Price in EUR
	MaxProducts      int
	MaxPhotos        int
	FeaturedListings int
	Features         []string
}

// GetSubscriptionPlans returns all available plans
func GetSubscriptionPlans() map[string]PlanConfig {
	return map[string]PlanConfig{
		"basic": {
			ID:               "basic",
			Name:             "Basic Plan",
			PriceID:          "price_1RhoN2B374sFE2QKRcrkRb6n",
			Price:            9.99,
			MaxProducts:      5,
			MaxPhotos:        5,
			FeaturedListings: 0,
			Features: []string{
				"5 Active Products",
				"5 Photos per Product",
				"Basic Contract Template",
				"Basic Analytics",
			},
		},
		"professional": {
			ID:               "professional",
			Name:             "Professional Plan",
			PriceID:          "price_1RhoRzB374sFE2QKicJZICyQ",
			Price:            19.99,
			MaxProducts:      15,
			MaxPhotos:        10,
			FeaturedListings: 2,
			Features: []string{
				"15 Active Products",
				"10 Photos per Product",
				"Customizable Contract",
				"Advanced Analytics",
				"2 Featured Listings/month",
			},
		},
		"business": {
			ID:               "business",
			Name:             "Business Plan",
			PriceID:          "price_1RhoUQB374sFE2QK5Po5AuND",
			Price:            39.99,
			MaxProducts:      999999, 
			MaxPhotos:        999999, 
			FeaturedListings: 10,
			Features: []string{
				"Unlimited Products",
				"Unlimited Photos",
				"Fully Custom Contract",
				"Full Analytics + Export",
				"10 Featured Listings/month",
			},
		},
	}
}
