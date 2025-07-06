package services

import "github.com/YpatiosCh/rentme/internal/repositories"

type Dependency struct {
	Repos     repositories.Repositories
	JwtSecret string
}

type ServiceContainer struct {
	userService         UserService
	authService         AuthService
	subscriptionService SubscriptionService
}

func NewServiceContainer(d Dependency) Services {
	return &ServiceContainer{
		userService:         NewUserService(d.Repos.User()),
		authService:         NewAuthService(d.Repos.User(), d.JwtSecret),
		subscriptionService: NewSubscriptionService(d.Repos.User()),
	}
}

func (s *ServiceContainer) User() UserService {
	return s.userService
}

func (s *ServiceContainer) Auth() AuthService {
	return s.authService
}

func (s *ServiceContainer) Subscription() SubscriptionService {
	return s.subscriptionService
}
