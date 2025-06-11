package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"booking-api/internal/usecase"

	"github.com/gorilla/mux"
)

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

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   uint      `json:"user_id"`
		CafeID   uint      `json:"cafe_id"`
		TableID  uint      `json:"table_id"`
		DateTime time.Time `json:"date_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	booking, err := h.Booking.CreateBooking(r.Context(), req.UserID, req.CafeID, req.TableID, req.DateTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"id": booking.ID, "status": booking.Status})
}

func (h *Handler) GetCafes(w http.ResponseWriter, r *http.Request) {
	cafes, err := h.Cafe.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cafes)
}

func (h *Handler) GetAvailableTables(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	tables, err := h.Table.GetAvailableTables(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tables)
}

func (h *Handler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	bookings, err := h.Booking.GetUserBookings(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	err := h.Booking.CancelBooking(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled"})
}
