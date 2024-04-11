package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/HINKOKO/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error with the rooms")
	}
	return 1, nil
}

// Inserts a room restriction -> due to a new reservation
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error here")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID - Search if a room is available for selected daters
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// test for failing at querying, choosing a 'pivot' date somehow
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("date conflict")
	}
	// If choosen date is after 2049-12-31 -> false, no availability
	if start.After(t) {
		return false, nil
	}

	// Otherwise, welcome to our Bed & Breakfast
	return true, nil
}

// SearchAvailabilityForAllRooms - returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// test for failing at querying, choosing a 'pivot' date somehow
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("date conflict")
	}
	// If choosen date is after 2049-12-31 -> false, no availability
	if start.After(t) {
		return rooms, nil
	}

	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)
	// Otherwise, welcome to our Bed & Breakfast
	return rooms, nil
}

// GetRoomByID - Pick a room by id from the DB
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("Some error haha!")
	}

	return room, nil
}

// GetUserID - Retrieve a user by its ID from database
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

// UpdateUser - Update a user in the Database
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate - authenticate a user
// Compare what the user type - hash - compare Hashes
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
