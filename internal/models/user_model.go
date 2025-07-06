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

	// Personal Information
	FirstName string  `gorm:"type:varchar(100);not null"`
	LastName  string  `gorm:"type:varchar(100);not null"`
	Phone     string  `gorm:"type:varchar(20);not null"` // Required for contact
	Website   *string `gorm:"type:varchar(255)"`         // Optional business website

	// Business Information (for owners)
	BusinessName        *string `gorm:"type:varchar(255)"` // Optional business name
	BusinessDescription *string `gorm:"type:text"`         // Business description for profile

	// Address Information (required for product location context)
	Address    string   `gorm:"type:varchar(255);not null"`
	City       string   `gorm:"type:varchar(100);not null"`
	Region     string   `gorm:"type:varchar(100);not null"` // Περιοχή
	PostalCode string   `gorm:"type:varchar(10);not null"`
	Prefecture string   `gorm:"type:varchar(100);not null"` // Νομός
	Country    string   `gorm:"type:varchar(50);not null;default:'Greece'"`
	Latitude   *float64 `gorm:"type:decimal(10,8)"` // For map functionality
	Longitude  *float64 `gorm:"type:decimal(11,8)"` // For map functionality

	// Account Status
	IsActive    bool       `gorm:"default:true;index"`
	IsBanned    bool       `gorm:"default:false;index"`
	BannedUntil *time.Time `gorm:"type:timestamp"`
	BanReason   *string    `gorm:"type:varchar(500)"`

	// Email Verification
	EmailVerified   bool       `gorm:"default:false;index"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp"`

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

	// Activity Tracking
	LastLoginAt       *time.Time `gorm:"type:timestamp;index"`
	LoginCount        int        `gorm:"default:0"`
	ProductViewsTotal int        `gorm:"default:0"` // Total views across all products
	ContactsTotal     int        `gorm:"default:0"` // Total contacts received

	// Communication Preferences
	EmailNotifications bool `gorm:"default:true"`  // Receive email notifications
	WeeklyReports      bool `gorm:"default:true"`  // Weekly analytics reports
	MarketingEmails    bool `gorm:"default:false"` // Marketing communications

	// Profile Settings
	ProfilePublic     bool `gorm:"default:true"` // Show profile to visitors
	ShowPhone         bool `gorm:"default:true"` // Show phone in product listings
	ShowEmail         bool `gorm:"default:true"` // Show email in product listings
	ShowWebsite       bool `gorm:"default:true"` // Show website in product listings
	ResponseTimeHours *int `gorm:"default:24"`   // Expected response time

	// Analytics Preferences
	AllowAnalytics bool `gorm:"default:true"` // Allow platform to track analytics

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

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) GetDisplayName() string {
	if u.BusinessName != nil && *u.BusinessName != "" {
		return *u.BusinessName
	}
	return u.GetFullName()
}

func (u *User) GetLocationString() string {
	return u.City + ", " + u.Region
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
