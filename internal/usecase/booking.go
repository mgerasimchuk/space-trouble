package usecase

import (
	"errors"
	"time"

	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
	"github.com/mgerasimchuk/space-trouble/internal/domain/repository"
	"github.com/mgerasimchuk/space-trouble/internal/domain/service"
	"github.com/sirupsen/logrus"
)

type BookingUsecase struct {
	bookingService *service.BookingService
	bookingRepo    repository.BookingRepository
	logger         *logrus.Logger
}

var internalError = errors.New("internal error")

func NewBookingUsecase(bookingService *service.BookingService, bookingRepo repository.BookingRepository, logger *logrus.Logger) *BookingUsecase {
	return &BookingUsecase{bookingService: bookingService, bookingRepo: bookingRepo, logger: logger}
}

func (u *BookingUsecase) CreateBooking(
	firstName string, lastName string, gender string, birthday time.Time,
	launchpadID string, destinationID string, launchDate time.Time,
) (*model.Booking, error) {
	b := model.CreateBooking(firstName, lastName, gender, birthday, launchpadID, destinationID, launchDate)
	if err := b.Validate(); err != nil {
		return nil, err
	}

	b, err := u.bookingRepo.Create(b)
	if err != nil {
		u.logger.Error(err)

		return nil, internalError
	}

	return b, nil
}

func (u *BookingUsecase) GetBookings(limit, offset *int) ([]*model.Booking, error) {
	if limit != nil && *limit < 1 {
		return nil, errors.New("\"limit\" param should be greater than 1")
	}
	if offset != nil && *offset < 0 {
		return nil, errors.New("\"offset\" param should be greater or equal 0")
	}

	bookings, err := u.bookingRepo.GetList(nil, limit, offset)
	if err != nil {
		u.logger.Error(err)

		return nil, internalError
	}

	return bookings, nil
}

func (u *BookingUsecase) DeleteBooking(id string) error {
	if id == "" {
		return errors.New("\"id\" can't be empty")
	}

	err := u.bookingRepo.Delete(id)
	if err != nil {
		u.logger.Error(err)

		return internalError
	}

	return nil
}

func (u *BookingUsecase) VerifyFirstAvailableBooking() {
	b, err := u.bookingRepo.GetAndModify(model.StatusCreated, model.StatusPending)
	if err != nil {
		u.logger.Error(err)

		return
	}

	if b == nil {
		return
	}

	err = u.bookingService.VerifyBooking(b)
	if err != nil {
		u.logger.Error(err)

		return
	}
}
