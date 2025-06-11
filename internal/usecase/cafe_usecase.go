package usecase

import (
	"booking-api/internal/entity"
	"context"
)

type CafeUsecase interface {
	GetAll(ctx context.Context) ([]entity.Cafe, error)
	GetByID(ctx context.Context, id uint) (*entity.Cafe, error)
}
