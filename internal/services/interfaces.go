package services

import (
	"github.com/YpatiosCh/rentme/internal/err"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/google/uuid"
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
}

type AuthService interface {
	RegisterUser(user *models.User, plainPassword string) (*models.User, *err.Error)
}
