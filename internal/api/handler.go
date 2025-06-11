package handler

import (
	"net/http"
	"strconv"
	"time"

	"booking-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success response structure
type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper function to handle errors
func handleError(c *gin.Context, err error) {
	switch err.Error() {
	case "user not found":
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case "table not found":
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case "table is not available for the selected time":
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
	case "invalid date and time":
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case "invalid booking status":
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
	}
}

type Handler struct {
	Booking usecase.BookingUsecase
	Table   usecase.TableUsecase
	Cafe    usecase.CafeUsecase
	User    usecase.UserUsecase
}

func NewHandler(booking usecase.BookingUsecase, table usecase.TableUsecase, cafe usecase.CafeUsecase, user usecase.UserUsecase) *Handler {
	return &Handler{
		Booking: booking,
		Table:   table,
		Cafe:    cafe,
		User:    user,
	}
}

func (h *Handler) CreateBooking(c *gin.Context) {
	var req struct {
		UserID   uint      `json:"user_id" binding:"required"`
		CafeID   uint      `json:"cafe_id" binding:"required"`
		TableID  uint      `json:"table_id" binding:"required"`
		DateTime time.Time `json:"date_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	booking, err := h.Booking.CreateBooking(
		c.Request.Context(),
		req.UserID,
		req.CafeID,
		req.TableID,
		req.DateTime,
	)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Booking created successfully",
		Data:    booking,
	})
}

func (h *Handler) GetCafes(c *gin.Context) {
	cafes, err := h.Cafe.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: cafes})
}

func (h *Handler) GetAvailableTables(c *gin.Context) {
	cafeIDStr := c.Param("cafe_id")
	cafeID, err := strconv.ParseUint(cafeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid cafe ID"})
		return
	}

	tables, err := h.Table.GetAvailableTables(c.Request.Context(), uint(cafeID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: tables})
}

func (h *Handler) GetUserBookings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	bookings, err := h.Booking.GetUserBookings(c.Request.Context(), uint(id))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: bookings})
}

func (h *Handler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid booking ID"})
		return
	}

	err = h.Booking.CancelBooking(c.Request.Context(), uint(id))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Booking cancelled successfully"})
}

// New booking endpoints

func (h *Handler) UpdateBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid booking ID"})
		return
	}

	var req struct {
		DateTime time.Time `json:"date_time" binding:"required"`
		Status   string    `json:"status" binding:"required,oneof=pending confirmed cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	err = h.Booking.UpdateBooking(c.Request.Context(), uint(id), req.DateTime, req.Status)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Booking updated successfully"})
}

func (h *Handler) CheckTableAvailability(c *gin.Context) {
	var req struct {
		TableID  uint      `json:"table_id" binding:"required"`
		DateTime time.Time `json:"date_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	available, err := h.Booking.IsTableAvailable(c.Request.Context(), req.TableID, req.DateTime)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: gin.H{"available": available},
	})
}

func (h *Handler) GetBookingsByCafe(c *gin.Context) {
	cafeIDStr := c.Param("cafe_id")
	cafeID, err := strconv.ParseUint(cafeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid cafe ID"})
		return
	}

	bookings, err := h.Booking.GetBookingsByCafe(c.Request.Context(), uint(cafeID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: bookings})
}

func (h *Handler) GetBookingsByDateRange(c *gin.Context) {
	cafeIDStr := c.Param("cafe_id")
	cafeID, err := strconv.ParseUint(cafeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid cafe ID"})
		return
	}

	var req struct {
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	bookings, err := h.Booking.GetBookingsByDateRange(c.Request.Context(), uint(cafeID), req.StartDate, req.EndDate)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: bookings})
}
