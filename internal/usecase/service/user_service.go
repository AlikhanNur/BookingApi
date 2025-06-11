package service

import (
	"booking-api/internal/entity"
	"booking-api/internal/repo"
	"context"
)

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	return s.userRepo.Create(ctx, user)
}
