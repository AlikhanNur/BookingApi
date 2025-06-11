package service

import (
	"booking-api/internal/entity"
	"booking-api/internal/repo"
	"context"
	"errors"
	"time"
)

type bookingService struct {
	bookingRepo repo.BookingRepository
	tableRepo   repo.TableRepository
	userRepo    repo.UserRepository
}

func NewBookingService(
	bookingRepo repo.BookingRepository,
	tableRepo repo.TableRepository,
	userRepo repo.UserRepository,
) *bookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		tableRepo:   tableRepo,
		userRepo:    userRepo,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID, cafeID, tableID uint, dateTime time.Time) (*entity.Booking, error) {
	// Проверим, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверим, существует ли столик
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil || table.CafeID != cafeID {
		return nil, errors.New("invalid table or cafe")
	}

	booking := &entity.Booking{
		UserID:   userID,
		TableID:  tableID,
		CafeID:   cafeID,
		DateTime: dateTime,
		Status:   "pending",
	}

	err = s.bookingRepo.Create(ctx, booking)
	return booking, err
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID uint) error {
	return s.bookingRepo.UpdateStatus(ctx, bookingID, "cancelled")
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uint) ([]entity.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}
