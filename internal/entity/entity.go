package entity

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null"`
	Name  string `json:"name" gorm:"not null"`
	Phone string `json:"phone"`
}

// Cafe represents a cafe/restaurant in the system
type Cafe struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name" gorm:"not null"`
	Address string  `json:"address" gorm:"not null"`
	Phone   string  `json:"phone"`
	Tables  []Table `json:"tables,omitempty" gorm:"foreignKey:CafeID"`
}

// Table represents a table in a cafe
type Table struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	CafeID   uint      `json:"cafe_id" gorm:"not null"`
	Number   int       `json:"number" gorm:"not null"`
	Seats    int       `json:"seats" gorm:"not null"`
	Bookings []Booking `json:"bookings,omitempty" gorm:"foreignKey:TableID"`
}

// Booking represents a reservation made by a user
type Booking struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	TableID  uint      `json:"table_id" gorm:"not null"`
	CafeID   uint      `json:"cafe_id" gorm:"not null"`
	DateTime time.Time `json:"date_time" gorm:"not null"`
	Status   string    `json:"status" gorm:"not null;default:'pending'"` // pending, confirmed, cancelled
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Table    Table     `json:"table,omitempty" gorm:"foreignKey:TableID"`
}
