package postgres

import (
	"context"

	"booking-api/internal/entity"
	"booking-api/internal/repo"

	"gorm.io/gorm"
)

type tableRepo struct {
	db *gorm.DB
}

func NewTableRepo(db *gorm.DB) repo.TableRepository {
	return &tableRepo{db: db}
}

func (r *tableRepo) Create(ctx context.Context, table *entity.Table) error {
	return r.db.WithContext(ctx).Create(table).Error
}

func (r *tableRepo) GetByID(ctx context.Context, id uint) (*entity.Table, error) {
	var t entity.Table
	err := r.db.WithContext(ctx).First(&t, id).Error
	return &t, err
}
func (r *tableRepo) GetAvailableByCafeID(ctx context.Context, cafeID uint) ([]entity.Table, error) {
	var tables []entity.Table
	// Предположим, что доступность — это отсутствие бронирований
	err := r.db.WithContext(ctx).
		Where("cafe_id = ?", cafeID).
		Preload("Bookings", "status = ?", "confirmed").
		Find(&tables).Error
	if err != nil {
		return nil, err
	}

	// Фильтруем только те столы, у которых нет активных бронирований
	var available []entity.Table
	for _, table := range tables {
		if len(table.Bookings) == 0 {
			available = append(available, table)
		}
	}
	return available, nil
}

func (r *tableRepo) ListByCafe(ctx context.Context, cafeID uint) ([]entity.Table, error) {
	var tables []entity.Table
	err := r.db.WithContext(ctx).Where("cafe_id = ?", cafeID).Find(&tables).Error
	return tables, err
}

func (r *tableRepo) Update(ctx context.Context, table *entity.Table) error {
	return r.db.WithContext(ctx).Save(table).Error
}

func (r *tableRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Table{}, id).Error
}
