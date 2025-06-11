package usecase

import (
	"booking-api/internal/entity"
	"context"
)

type TableUsecase interface {
	GetAvailableTables(ctx context.Context, cafeID uint) ([]entity.Table, error)
	GetByID(ctx context.Context, id uint) (*entity.Table, error)
}
