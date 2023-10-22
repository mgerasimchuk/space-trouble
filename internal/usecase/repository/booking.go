package repository

import (
	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
)

const DefaultBookingGetListLimit = 20

type BookingRepository interface {
	Create(b *model.Booking) (*model.Booking, error)
	Save(b *model.Booking) error
	GetList(status *string, limit, offset *int) ([]*model.Booking, error)
	Delete(id string) error
	// GetAndModify provides way to find single entity by searchStatus, and set this status immediately to modifyStatus
	GetAndModify(searchStatus, modifyStatus string) (*model.Booking, error)
}
