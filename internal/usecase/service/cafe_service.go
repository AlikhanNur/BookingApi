package service

import (
	"booking-api/internal/entity"
	"booking-api/internal/repo"
	"context"
)

type cafeService struct {
	cafeRepo repo.CafeRepository
}

func NewCafeService(cafeRepo repo.CafeRepository) *cafeService {
	return &cafeService{cafeRepo: cafeRepo}
}

func (s *cafeService) GetAll(ctx context.Context) ([]entity.Cafe, error) {
	return s.cafeRepo.GetAll(ctx)
}

func (s *cafeService) GetByID(ctx context.Context, id uint) (*entity.Cafe, error) {
	return s.cafeRepo.GetByID(ctx, id)
}
