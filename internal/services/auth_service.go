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

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterUser handles user creation and password hashing
func (s *authService) RegisterUser(user *models.User, plainPassword string) (*models.User, *err.Error) {
	var e err.Error

	// Check if user already exists
	existingUser, dbErr := s.userRepo.GetUserByEmail(user.Email)
	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		e.Code = 500
		e.Message = "could not check existing user"
		return nil, &e
	}
	if existingUser != nil {
		e.Code = 400
		e.Message = "user already exists"
		return nil, &e
	}

	// Hash password
	hashedPassword, hashErr := hash.Password(plainPassword)
	if hashErr != nil {
		e.Code = 500
		e.Message = "failed to hash password"
		return nil, &e
	}
	user.PasswordHash = hashedPassword

	// Save new user
	createdUser, saveErr := s.userRepo.CreateUser(user)
	if saveErr != nil {
		e.Code = 500
		e.Message = "failed to create user"
		return nil, &e
	}

	return createdUser, nil
}

// LoginUser validates user credentials and returns JWT
func (s *authService) LoginUser(email, plainPassword string) (*models.User, string, *err.Error) {
	var e err.Error

	user, dbErr := s.userRepo.GetUserByEmail(email)
	if dbErr != nil {
		e.Code = 500
		e.Message = "error fetching user"
		return nil, "", &e
	}
	if user == nil {
		e.Code = 400
		e.Message = "invalid credentials"
		return nil, "", &e
	}

	if passErr := hash.Check(plainPassword, user.PasswordHash); passErr != nil {
		e.Code = 400
		e.Message = "invalid credentials"
		return nil, "", &e
	}

	token, tokenErr := s.GenerateToken(user.ID)
	if tokenErr != nil {
		return nil, "", tokenErr
	}

	return user, token, nil
}

// GenerateToken creates a JWT for the given user ID
func (s *authService) GenerateToken(userID uuid.UUID) (string, *err.Error) {
	var e err.Error

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, signErr := token.SignedString([]byte(s.jwtSecret))
	if signErr != nil {
		e.Code = 500
		e.Message = "failed to sign token"
		return "", &e
	}

	return tokenString, nil
}

// ValidateToken verifies the JWT and extracts the user ID
func (s *authService) ValidateToken(tokenString string) (*uuid.UUID, *err.Error) {
	var e err.Error

	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if parseErr != nil {
		e.Code = 401
		e.Message = "invalid token"
		return nil, &e
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		e.Code = 401
		e.Message = "invalid token claims"
		return nil, &e
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		e.Code = 401
		e.Message = "invalid token payload"
		return nil, &e
	}

	userID, uuidErr := uuid.Parse(userIDStr)
	if uuidErr != nil {
		e.Code = 401
		e.Message = "invalid user ID in token"
		return nil, &e
	}

	return &userID, nil
}
