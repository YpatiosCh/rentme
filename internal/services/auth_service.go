package services

import (
	"errors"
	"fmt"
	"time"

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

func (s *authService) RegisterUser(user *models.User, plainPassword string) (*models.User, error) {
	// Check if the user already exists
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// hash the password
	hashedPassword, err := hash.Password(plainPassword)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPassword

	// last login time is set to current time
	currentTime := time.Now()
	user.LastLoginAt = &currentTime

	// increment the login count
	user.LoginCount++

	// create the user in the repository
	return s.userRepo.CreateUser(user)
}

func (s *authService) LoginUser(email, plainPassword string) (*models.User, error) {
	// Fetch the user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials: user not found")
	}
	// Check if the user is active or banned
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}
	if user.IsBanned {
		return nil, errors.New("account is banned")
	}

	// Verify the password
	if err = hash.Check(plainPassword, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login time and increment login count
	currentTime := time.Now()
	user.LastLoginAt = &currentTime
	user.LoginCount++

	// Save the updated user
	return s.userRepo.UpdateUser(user)
}

func (s *authService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user ID in token")
		}

		return &userID, nil
	}

	return nil, errors.New("invalid token")
}
