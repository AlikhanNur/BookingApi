package usecase

import (
	"booking-api/internal/entity"
	"context"
	"time"
)

type BookingUsecase interface {
	CreateBooking(ctx context.Context, userID, cafeID, tableID uint, dateTime time.Time) (*entity.Booking, error)
	CancelBooking(ctx context.Context, bookingID uint) error
	GetUserBookings(ctx context.Context, userID uint) ([]entity.Booking, error)
}
