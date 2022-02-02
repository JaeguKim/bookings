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
}