package service

import (
	"booking-api/internal/entity"
	"booking-api/internal/repo"
	"context"
)

type tableService struct {
	tableRepo repo.TableRepository
}

func NewTableService(tableRepo repo.TableRepository) *tableService {
	return &tableService{tableRepo: tableRepo}
}

func (s *tableService) GetAvailableTables(ctx context.Context, cafeID uint) ([]entity.Table, error) {
	return s.tableRepo.GetAvailableByCafeID(ctx, cafeID)
}

func (s *tableService) GetByID(ctx context.Context, id uint) (*entity.Table, error) {
	return s.tableRepo.GetByID(ctx, id)
}
