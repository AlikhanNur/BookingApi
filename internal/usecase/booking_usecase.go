package usecase

import (
	"booking-api/internal/entity"
	"context"
	"time"
)

type BookingUsecase interface {
	// Basic booking operations
	CreateBooking(ctx context.Context, userID, cafeID, tableID uint, dateTime time.Time) (*entity.Booking, error)
	CancelBooking(ctx context.Context, bookingID uint) error
	GetUserBookings(ctx context.Context, userID uint) ([]entity.Booking, error)

	// New operations
	IsTableAvailable(ctx context.Context, tableID uint, dateTime time.Time) (bool, error)
	UpdateBooking(ctx context.Context, bookingID uint, dateTime time.Time, status string) error
	GetBookingsByDateRange(ctx context.Context, cafeID uint, startDate, endDate time.Time) ([]entity.Booking, error)
	GetBookingsByCafe(ctx context.Context, cafeID uint) ([]entity.Booking, error)
}
