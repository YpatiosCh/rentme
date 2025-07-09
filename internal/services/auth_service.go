package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/YpatiosCh/rentme/internal/err"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/repositories"
	"github.com/YpatiosCh/rentme/pkg/hash"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwt string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwt,
	}
}

func (s *authService) RegisterUser(user *models.User, plainPassword string) (*models.User, *err.Error) {
	var error err.Error
	// Check if the user already exists
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		error.Code = 500
		error.Message = "Could not get user"
		return nil, &error
	}
	if existingUser != nil {
		error.Code = 400
		error.Message = "user already exists"
		return nil, &error
	}

	// hash the password
	hashedPassword, err := hash.Password(plainPassword)
	if err != nil {
		error.Code = 500
		error.Message = "failed to hash password"
		return nil, &error
	}
	user.PasswordHash = hashedPassword

	// create the user in the repository
	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		error.Code = 500
		error.Message = "failed to create user in database"
		return nil, &error
	}

	return user, nil
}

func (s *authService) LoginUser(email, plainPassword string) (*models.User, *err.Error) {
	var error err.Error
	// Fetch the user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		error.Code = 500
		error.Message = "error fetching user"
		return nil, &error
	}
	if user == nil {
		error.Code = 400
		error.Message = "invalid credentials: user not found"
		return nil, &error
	}

	// Verify the password
	if err = hash.Check(plainPassword, user.PasswordHash); err != nil {
		error.Code = 400
		error.Message = "invalid credentials"
		return nil, &error
	}

	// Save the updated user
	user, err = s.userRepo.UpdateUser(user)
	if err != nil {
		error.Code = 500
		error.Message = "failed to update user in database"
		return nil, &error
	}

	return user, nil
}

func (s *authService) GenerateToken(userID uuid.UUID) (string, *err.Error) {
	var error err.Error
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		error.Code = 500
		error.Message = "failed to sign token"
		return "", &error
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*uuid.UUID, *err.Error) {
	var er err.Error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		er.Code = 500
		er.Message = "unexpected signing method"
		return nil, &er
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			er.Code = 500
			er.Message = "invalid token claims"
			return nil, &er
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			er.Code = 500
			er.Message = "invalid user ID in token"
			return nil, &er
		}

		return &userID, nil
	}
	er.Code = 500
	er.Message = "invalid token"
	return nil, &er
}
