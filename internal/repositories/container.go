package repositories

import "gorm.io/gorm"

type RepositoryContainer struct {
	userRepo UserRepository
}

func NewRepositoryContainer(db *gorm.DB) Repositories {
	return &RepositoryContainer{
		userRepo: NewUserRepository(db),
	}
}

func (r *RepositoryContainer) User() UserRepository {
	return r.userRepo
}
