package services

import "github.com/YpatiosCh/rentme/internal/repositories"

type subscriptionService struct {
	userRepo repositories.UserRepository
}

func NewSubscriptionService(userRepo repositories.UserRepository) SubscriptionService {
	return &subscriptionService{
		userRepo: userRepo,
	}
}
