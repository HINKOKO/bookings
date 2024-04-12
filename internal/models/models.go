package models

import (
	"time"
)

// User 'Schema' if you like
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction - restrictions model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	RoomID    int
	StartDate time.Time
	EndDate   time.Time
	Room      Room
	Processed int
}

// RoomRestriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation // supa + If I need them
	Restriction   Restriction // supa + If I need them
}

// MAilData - Holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
