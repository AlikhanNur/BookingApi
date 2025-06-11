package main

import (
	api "booking-api/internal/api"
	postgresrepo "booking-api/internal/repo/postgres"
	"booking-api/internal/usecase/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Connect to database
	dsn := "host=localhost user=postgres password=1234 dbname=bookingdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// // 2. Repositories
	// bookingRepo := repo.NewBookingRepo(db)
	// tableRepo := repo.NewTableRepo(db)
	// userRepo := repo.NewUserRepo(db)
	// cafeRepo := repo.NewCafeRepo(db)

	bookingRepo := postgresrepo.NewBookingRepo(db)
	tableRepo := postgresrepo.NewTableRepo(db)
	userRepo := postgresrepo.NewUserRepo(db)
	cafeRepo := postgresrepo.NewCafeRepo(db)

	// 3. Services
	bookingService := service.NewBookingService(bookingRepo, tableRepo, userRepo)
	tableService := service.NewTableService(tableRepo)
	userService := service.NewUserService(userRepo)
	cafeService := service.NewCafeService(cafeRepo)

	// 4. Handler
	h := api.NewHandler(bookingService, tableService, cafeService, userService)

	// 5. Gin Router
	r := gin.Default()

	// Booking routes
	r.POST("/api/bookings", h.CreateBooking)
	r.PUT("/api/bookings/:id", h.UpdateBooking) // New: update booking
	r.DELETE("/api/bookings/:id", h.CancelBooking)
	r.GET("/api/users/:id/bookings", h.GetUserBookings)
	r.POST("/api/bookings/check-availability", h.CheckTableAvailability) // New: check table availability

	// Cafe and tables
	r.GET("/api/cafes", h.GetCafes)
	r.GET("/api/cafes/:cafe_id/available-tables", h.GetAvailableTables) // Use cafe_id for consistency
	r.GET("/api/cafes/:cafe_id/bookings", h.GetBookingsByCafe) // New: get bookings by cafe
	r.POST("/api/cafes/:cafe_id/bookings/date-range", h.GetBookingsByDateRange) // New: get bookings by date range

	// 6. Run server
	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
