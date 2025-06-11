package postgres

import (
	"booking-api/internal/entity"
	"context"

	"gorm.io/gorm"
)

type CafeRepo struct {
	db *gorm.DB
}

func NewCafeRepo(db *gorm.DB) *CafeRepo {
	return &CafeRepo{db: db}
}

func (r *CafeRepo) Create(ctx context.Context, cafe *entity.Cafe) error {
	return r.db.WithContext(ctx).Create(cafe).Error
}

func (r *CafeRepo) GetByID(ctx context.Context, id uint) (*entity.Cafe, error) {
	var cafe entity.Cafe
	err := r.db.WithContext(ctx).Preload("Tables").First(&cafe, id).Error
	return &cafe, err
}

func (r *CafeRepo) GetAll(ctx context.Context) ([]entity.Cafe, error) {
	var cafes []entity.Cafe
	err := r.db.WithContext(ctx).Preload("Tables").Find(&cafes).Error
	return cafes, err
}

func (r *CafeRepo) Update(ctx context.Context, cafe *entity.Cafe) error {
	return r.db.WithContext(ctx).Save(cafe).Error
}

func (r *CafeRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Cafe{}, id).Error
}
