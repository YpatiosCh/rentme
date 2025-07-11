package config

import (
	"fmt"
	"os"
	"strconv"

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

// GetSubscriptionPlans returns all available plans with Price IDs from environment
func GetSubscriptionPlans() map[string]PlanConfig {
	return map[string]PlanConfig{
		"basic": {
			ID:               "basic",
			Name:             "Basic Plan",
			PriceID:          GetEnvOrDefault("STRIPE_PRICE_BASIC", "price_1RhoN2B374sFE2QKRcrkRb6n"),
			Price:            GetFloatEnvOrDefault("BASIC_PLAN_PRICE", 9.99),
			MaxProducts:      GetIntEnvOrDefault("BASIC_MAX_PRODUCTS", 5),
			MaxPhotos:        GetIntEnvOrDefault("BASIC_MAX_PHOTOS", 5),
			FeaturedListings: GetIntEnvOrDefault("BASIC_FEATURED_LISTINGS", 0),
			Features: []string{
				"3 Ενεργά προϊόντα",
				"3 Φωτογραφίες/προϊόν",
				"Βασικό συμβόλαιο ενοικίασης",
			},
		},
		"professional": {
			ID:               "professional",
			Name:             "Professional Plan",
			PriceID:          GetEnvOrDefault("STRIPE_PRICE_PROFESSIONAL", "price_1RhoRzB374sFE2QKicJZICyQ"),
			Price:            GetFloatEnvOrDefault("PROFESSIONAL_PLAN_PRICE", 19.99),
			MaxProducts:      GetIntEnvOrDefault("PROFESSIONAL_MAX_PRODUCTS", 10),
			MaxPhotos:        GetIntEnvOrDefault("PROFESSIONAL_MAX_PHOTOS", 10),
			FeaturedListings: GetIntEnvOrDefault("PROFESSIONAL_FEATURED_LISTINGS", 2),
			Features: []string{
				"10 Ενεργά προϊόντα",
				"10 Φωτογραφίες/προϊόν",
				"Προσαρμόσιμο συμβόλαιο",
				"2 Featured listings/μήνα",
			},
		},
		"business": {
			ID:               "business",
			Name:             "Business Plan",
			PriceID:          GetEnvOrDefault("STRIPE_PRICE_BUSINESS", "price_1RhoUQB374sFE2QK5Po5AuND"),
			Price:            GetFloatEnvOrDefault("BUSINESS_PLAN_PRICE", 39.99),
			MaxProducts:      GetIntEnvOrDefault("BUSINESS_MAX_PRODUCTS", 999999),
			MaxPhotos:        GetIntEnvOrDefault("BUSINESS_MAX_PHOTOS", 999999),
			FeaturedListings: GetIntEnvOrDefault("BUSINESS_FEATURED_LISTINGS", 10),
			Features: []string{
				"Απεριόριστα προϊόντα",
				"Απεριόριστες φωτογραφίες",
				"Πλήρως custom συμβόλαιο",
				"10 Featured listings/μήνα",
			},
		},
	}
}

// GetPlanByID returns a specific plan configuration
func GetPlanByID(planID string) (PlanConfig, bool) {
	plans := GetSubscriptionPlans()
	plan, exists := plans[planID]
	return plan, exists
}

// GetPlanByPriceID returns plan configuration by Stripe Price ID
func GetPlanByPriceID(priceID string) (PlanConfig, bool) {
	plans := GetSubscriptionPlans()
	for _, plan := range plans {
		if plan.PriceID == priceID {
			return plan, true
		}
	}
	return PlanConfig{}, false
}

// ValidatePlanExists checks if a plan ID exists
func ValidatePlanExists(planID string) bool {
	_, exists := GetPlanByID(planID)
	return exists
}

// Helper functions for environment variables with defaults
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetFloatEnvOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// Configuration validation
func ValidateStripeConfig() error {
	required := []string{
		"STRIPE_SECRET_KEY",
		"STRIPE_PUBLISHABLE_KEY",
		"STRIPE_WEBHOOK_SECRET",
		"STRIPE_PRICE_BASIC",
		"STRIPE_PRICE_PROFESSIONAL",
		"STRIPE_PRICE_BUSINESS",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return nil
}
