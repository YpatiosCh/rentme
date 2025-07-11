package repositories

import "github.com/YpatiosCh/rentme/internal/models"

type Repositories interface {
	User() UserRepository
}

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetAllUsers() (*[]models.User, error)
	GetUserByCustomerID(customerID string) (*models.User, error)
}
