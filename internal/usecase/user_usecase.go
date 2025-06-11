package usecase

import (
	"booking-api/internal/entity"
	"context"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}
