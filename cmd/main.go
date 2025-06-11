package main

import (
	handler "booking-api/internal/api"
	// "booking-api/internal/repo"
	postgresrepo "booking-api/internal/repo/postgres"

	"booking-api/internal/usecase/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Подключение к БД
	dsn := "host=localhost user=postgres password=1234 dbname=bookingdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 2. Репозитории
	// bookingRepo := repo.NewBookingRepo(db)
	// tableRepo := repo.NewTableRepo(db)
	// userRepo := repo.NewUserRepo(db)
	// cafeRepo := repo.NewCafeRepo(db)
	bookingRepo := postgresrepo.NewBookingRepo(db)
	tableRepo := postgresrepo.NewTableRepo(db)
	userRepo := postgresrepo.NewUserRepo(db)
	cafeRepo := postgresrepo.NewCafeRepo(db)

	// 3. Сервисы (usecase)
	bookingService := service.NewBookingService(bookingRepo, tableRepo, userRepo)
	tableService := service.NewTableService(tableRepo)
	userService := service.NewUserService(userRepo)
	cafeService := service.NewCafeService(cafeRepo)

	// 4. HTTP handler
	h := handler.NewHandler(bookingService, tableService, cafeService, userService)

	// 5. Router (mux)
	r := mux.NewRouter()

	// Booking routes
	r.HandleFunc("/api/bookings", h.CreateBooking).Methods("POST")
	r.HandleFunc("/api/bookings/{id}", h.CancelBooking).Methods("DELETE")
	r.HandleFunc("/api/users/{id}/bookings", h.GetUserBookings).Methods("GET")

	// Cafe and tables
	r.HandleFunc("/api/cafes", h.GetCafes).Methods("GET")
	r.HandleFunc("/api/cafes/{id}/available-tables", h.GetAvailableTables).Methods("GET")

	// 6. Запуск сервера
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
