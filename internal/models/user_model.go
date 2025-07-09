package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`

	// Subscription Information
	SubscriptionPlan   *string    `gorm:"type:varchar(20);index"` // "basic", "professional", "business"
	SubscriptionStatus *string    `gorm:"type:varchar(20);index"` // "active", "canceled", "past_due", "unpaid"
	SubscriptionID     *string    `gorm:"type:varchar(255)"`      // Stripe subscription ID
	CustomerID         *string    `gorm:"type:varchar(255)"`      // Stripe customer ID
	PlanStartDate      *time.Time `gorm:"type:timestamp"`
	PlanEndDate        *time.Time `gorm:"type:timestamp"`
	TrialEndDate       *time.Time `gorm:"type:timestamp"`

	// Usage Limits (based on subscription plan)
	MaxProducts           int `gorm:"default:0"` // Max products allowed
	MaxPhotosPerProduct   int `gorm:"default:0"` // Max photos per product
	FeaturedListingsUsed  int `gorm:"default:0"` // Used featured listings this month
	FeaturedListingsLimit int `gorm:"default:0"` // Max featured listings per month
	ProductCount          int `gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt *gorm.DeletedAt `gorm:"index"` // Soft delete support
}

// BeforeCreate hook to set UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// Helper methods for subscription management
func (u *User) IsSubscribed() bool {
	return u.SubscriptionPlan != nil &&
		u.SubscriptionStatus != nil &&
		*u.SubscriptionStatus == "active" &&
		(u.PlanEndDate == nil || u.PlanEndDate.After(time.Now()))
}

func (u *User) CanAddProduct() bool {
	// Check if user has active subscription and hasn't exceeded product limit
	return u.IsSubscribed() && u.getCurrentProductCount() < u.MaxProducts
}

func (u *User) CanUseFeaturedListing() bool {
	return u.IsSubscribed() && u.FeaturedListingsUsed < u.FeaturedListingsLimit
}

func (u *User) GetPlanDisplayName() string {
	if u.SubscriptionPlan == nil {
		return "No Plan"
	}

	switch *u.SubscriptionPlan {
	case "basic":
		return "Basic Plan"
	case "professional":
		return "Professional Plan"
	case "business":
		return "Business Plan"
	default:
		return "Unknown Plan"
	}
}

// This would need to be implemented with a count query in the repository
func (u *User) getCurrentProductCount() int {
	// This should be implemented in the service layer with a repository call
	// For now, return 0 as placeholder
	return 0
}

// SetSubscriptionLimits sets the appropriate limits based on plan
func (u *User) SetSubscriptionLimits(plan string) {
	u.SubscriptionPlan = &plan

	switch plan {
	case "basic":
		u.MaxProducts = 5
		u.MaxPhotosPerProduct = 5
		u.FeaturedListingsLimit = 0
	case "professional":
		u.MaxProducts = 15
		u.MaxPhotosPerProduct = 10
		u.FeaturedListingsLimit = 2
	case "business":
		u.MaxProducts = 999999         // Unlimited
		u.MaxPhotosPerProduct = 999999 // Unlimited
		u.FeaturedListingsLimit = 10
	}
}
