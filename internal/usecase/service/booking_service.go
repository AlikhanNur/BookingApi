package service

import (
	"booking-api/internal/entity"
	"booking-api/internal/repo"
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrTableNotFound    = errors.New("table not found")
	ErrTableUnavailable = errors.New("table is not available for the selected time")
	ErrInvalidDateTime  = errors.New("invalid date and time")
	ErrInvalidStatus    = errors.New("invalid booking status")
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
	// Validate date and time
	if dateTime.Before(time.Now()) {
		return nil, ErrInvalidDateTime
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if table exists and belongs to the cafe
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil || table.CafeID != cafeID {
		return nil, ErrTableNotFound
	}

	// Check table availability
	available, err := s.IsTableAvailable(ctx, tableID, dateTime)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrTableUnavailable
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
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Only allow cancellation of pending or confirmed bookings
	if booking.Status != "pending" && booking.Status != "confirmed" {
		return ErrInvalidStatus
	}

	return s.bookingRepo.UpdateStatus(ctx, bookingID, "cancelled")
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uint) ([]entity.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}

// New methods

func (s *bookingService) IsTableAvailable(ctx context.Context, tableID uint, dateTime time.Time) (bool, error) {
	// Get all bookings for the table on the specified date
	bookings, err := s.bookingRepo.ListByDateAndCafe(ctx, tableID, dateTime)
	if err != nil {
		return false, err
	}

	// Check if there are any active bookings for the specified time
	for _, booking := range bookings {
		if booking.Status != "cancelled" && booking.DateTime.Equal(dateTime) {
			return false, nil
		}
	}

	return true, nil
}

func (s *bookingService) UpdateBooking(ctx context.Context, bookingID uint, dateTime time.Time, status string) error {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Validate status
	if status != "pending" && status != "confirmed" && status != "cancelled" {
		return ErrInvalidStatus
	}

	// If changing date/time, check availability
	if !dateTime.Equal(booking.DateTime) {
		available, err := s.IsTableAvailable(ctx, booking.TableID, dateTime)
		if err != nil {
			return err
		}
		if !available {
			return ErrTableUnavailable
		}
	}

	booking.DateTime = dateTime
	booking.Status = status
	return s.bookingRepo.UpdateStatus(ctx, bookingID, status)
}

func (s *bookingService) GetBookingsByDateRange(ctx context.Context, cafeID uint, startDate, endDate time.Time) ([]entity.Booking, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidDateTime
	}

	var allBookings []entity.Booking
	currentDate := startDate
	for !currentDate.After(endDate) {
		dayBookings, err := s.bookingRepo.ListByDateAndCafe(ctx, cafeID, currentDate)
		if err != nil {
			return nil, err
		}
		allBookings = append(allBookings, dayBookings...)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return allBookings, nil
}

func (s *bookingService) GetBookingsByCafe(ctx context.Context, cafeID uint) ([]entity.Booking, error) {
	// Get all bookings for today
	return s.bookingRepo.ListByDateAndCafe(ctx, cafeID, time.Now())
}
