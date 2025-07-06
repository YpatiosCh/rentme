package services

import (
	"github.com/YpatiosCh/rentme/internal/err"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
)

type Services interface {
	User() UserService
	Auth() AuthService
	Subscription() SubscriptionService
}

type UserService interface {
	GetUserByID(serID uuid.UUID) (*models.User, *err.Error)
	GetUserByEmail(email string) (*models.User, *err.Error)
	UpdateUser(user *models.User) (*models.User, *err.Error)
	CreateUser(user *models.User) (*models.User, *err.Error)
	GetAllUsers() (*[]models.User, *err.Error)
}

type SubscriptionService interface {
	GetCustomerSubscriptions(customerID string) ([]*stripe.Subscription, error)
	CancelSubscription(userID string) error
	UpdateSubscriptionPlan(userID, newPlanID string) error
	CreateCheckoutSessionWithUserData(userEmail, planID, successURL, cancelURL string, metadata map[string]string) (string, error)
	HandleSuccessfulPaymentFromSession(sessionID string) (*models.User, error)
}

type AuthService interface {
	RegisterUser(user *models.User, plainPassword string) (*models.User, error)
}
