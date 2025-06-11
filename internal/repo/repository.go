package repo

import (
	"context"
	"time"

	"booking-api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type TableRepository interface {
	Create(ctx context.Context, table *entity.Table) error
	GetAvailableByCafeID(ctx context.Context, cafeID uint) ([]entity.Table, error)
	GetByID(ctx context.Context, id uint) (*entity.Table, error)
	ListByCafe(ctx context.Context, cafeID uint) ([]entity.Table, error)
	Update(ctx context.Context, table *entity.Table) error
	Delete(ctx context.Context, id uint) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByUserID(ctx context.Context, userID uint) ([]entity.Booking, error)
	GetByID(ctx context.Context, id uint) (*entity.Booking, error)
	ListByUser(ctx context.Context, userID uint) ([]entity.Booking, error)
	ListByDateAndCafe(ctx context.Context, cafeID uint, date time.Time) ([]entity.Booking, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Delete(ctx context.Context, id uint) error
}

type CafeRepository interface {
	Create(ctx context.Context, cafe *entity.Cafe) error
	GetByID(ctx context.Context, id uint) (*entity.Cafe, error)
	GetAll(ctx context.Context) ([]entity.Cafe, error)
	Update(ctx context.Context, cafe *entity.Cafe) error
	Delete(ctx context.Context, id uint) error
}
