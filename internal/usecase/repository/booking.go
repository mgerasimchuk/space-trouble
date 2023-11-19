package repository

import (
	"github.com/mgerasimchuk/space-trouble/internal/entity"
)

//go:generate mockgen -source=booking.go -destination=../../adapter/repository/mock/booking.go -package=mock

const DefaultBookingGetListLimit = 20

type BookingRepository interface {
	Create(b *entity.Booking) (*entity.Booking, error)
	Save(b *entity.Booking) error
	GetList(status *string, limit, offset *int) ([]*entity.Booking, error)
	Delete(id string) error
	// GetAndModify provides way to find single entity by searchStatus, and set this status immediately to modifyStatus
	GetAndModify(searchStatus, modifyStatus string) (*entity.Booking, error)
}
