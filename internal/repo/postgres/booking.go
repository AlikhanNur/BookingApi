package postgres

import (
	"context"
	"time"

	"booking-api/internal/entity"
	"booking-api/internal/repo"

	"gorm.io/gorm"
)

type bookingRepo struct {
	db *gorm.DB
}

func NewBookingRepo(db *gorm.DB) repo.BookingRepository {
	return &bookingRepo{db: db}
}

func (r *bookingRepo) Create(ctx context.Context, b *entity.Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *bookingRepo) GetByID(ctx context.Context, id uint) (*entity.Booking, error) {
	var b entity.Booking
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Table").
		First(&b, id).Error
	return &b, err
}
func (r *bookingRepo) GetByUserID(ctx context.Context, userID uint) ([]entity.Booking, error) {
	var bookings []entity.Booking
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepo) ListByUser(ctx context.Context, userID uint) ([]entity.Booking, error) {
	var bookings []entity.Booking
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date_time").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) ListByDateAndCafe(ctx context.Context, cafeID uint, date time.Time) ([]entity.Booking, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var bookings []entity.Booking
	err := r.db.WithContext(ctx).
		Where("cafe_id = ? AND date_time BETWEEN ? AND ?", cafeID, start, end).
		Order("date_time").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *bookingRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Booking{}, id).Error
}
