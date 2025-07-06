package services

import (
	"github.com/YpatiosCh/rentme/internal/err"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/repositories"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(serID uuid.UUID) (*models.User, *err.Error) {
	var error err.Error
	user, err := s.userRepo.GetUserByID(serID.String())
	if err != nil {
		error.Code = 500
		error.Message = err.Error()
		error.Err = err
		return nil, &error
	}
	if user == nil {
		error.Code = 404
		error.Message = "user not found"
		error.Err = nil
		return nil, &error
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, *err.Error) {
	var error err.Error
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		error.Code = 500
		error.Message = err.Error()
		error.Err = err
		return nil, &error
	}
	if user == nil {
		error.Code = 404
		error.Message = "user not found"
		error.Err = nil
		return nil, &error
	}
	return user, nil
}

func (s *userService) UpdateUser(user *models.User) (*models.User, *err.Error) {
	var error err.Error
	user, err := s.userRepo.UpdateUser(user)
	if err != nil {
		error.Code = 500
		error.Message = err.Error()
		error.Err = err
		return nil, &error
	}
	if user == nil {
		error.Code = 404
		error.Message = "user not found"
		error.Err = nil
		return nil, &error
	}
	return user, nil
}

func (s *userService) CreateUser(user *models.User) (*models.User, *err.Error) {
	var error err.Error
	user, err := s.userRepo.CreateUser(user)
	if err != nil {
		error.Code = 500
		error.Message = err.Error()
		error.Err = err
		return nil, &error
	}
	if user == nil {
		error.Code = 404
		error.Message = "user not found"
		error.Err = nil
		return nil, &error
	}
	return user, nil
}

func (s *userService) GetAllUsers() (*[]models.User, *err.Error) {
	var error err.Error
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		error.Code = 500
		error.Message = err.Error()
		error.Err = err
		return nil, &error
	}
	if users == nil {
		error.Code = 404
		error.Message = "no users found"
		error.Err = nil
		return nil, &error
	}
	return users, nil
}
