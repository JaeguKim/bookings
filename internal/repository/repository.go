package repository

import (
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

type DatabaseRepo interface {
	AllUser() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestrictions(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool ,error)
	SearchAvailabilityForAllRoom(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
}