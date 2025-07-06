package repositories

import (
	"github.com/YpatiosCh/rentme/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := u.db.Where("id = ?", userID).Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // not an error, user just doesn't exist
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := u.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
