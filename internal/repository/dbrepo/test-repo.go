package dbrepo

import (
	"errors"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

func (m *testDBRepo) AllUser() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestrictions(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool ,error) {

	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRoom(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	layout := "2006-01-02"
	if start.Format(layout) == "2050-10-01" && end.Format(layout) == "2050-10-02" {
		return rooms, nil
	} else if start.Format(layout) == "2050-11-11" && end.Format(layout) == "2050-11-12" {
		return rooms, errors.New("error")
	}
	rooms = append(rooms,models.Room{})
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room,nil
}

// GetUserByID returns user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user,nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0,"",nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations,nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations,nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation
	return res,nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}