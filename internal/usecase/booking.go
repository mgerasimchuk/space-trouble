package usecase

import (
	"errors"
	"github.com/mgerasimchuk/space-trouble/internal/domain/service"
	"github.com/mgerasimchuk/space-trouble/internal/domain/service/dto"
	"github.com/mgerasimchuk/space-trouble/internal/usecase/repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
)

type BookingUsecase struct {
	bookingVerifierSvc *service.BookingVerifierService
	bookingRepo        repository.BookingRepository
	launchpadRepo      repository.LaunchpadRepository
	landpadRepo        repository.LandpadRepository
}

func NewBookingUsecase(bookingSvc *service.BookingVerifierService, bookingRepo repository.BookingRepository, launchpadRepo repository.LaunchpadRepository, landpadRepo repository.LandpadRepository) *BookingUsecase {
	return &BookingUsecase{bookingVerifierSvc: bookingSvc, bookingRepo: bookingRepo, launchpadRepo: launchpadRepo, landpadRepo: landpadRepo}
}

var internalError = errors.New("internal error")

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
		return internalError
	}

	return nil
}

func (u *BookingUsecase) VerifyFirstAvailableBooking() error {
	b, err := u.bookingRepo.GetAndModify(model.StatusCreated, model.StatusPending)
	if err != nil {
		return err
	}

	bookingVerifyDTO := dto.BookingVerifyDTO{}
	g := errgroup.Group{}
	g.Go(func() (err error) {
		bookingVerifyDTO.IsLaunchpadExists, err = u.launchpadRepo.IsExists(b.LaunchpadID())
		return err
	})
	g.Go(func() (err error) {
		bookingVerifyDTO.IsLaunchpadActive, err = u.launchpadRepo.IsActive(b.LaunchpadID())
		return err
	})
	g.Go(func() (err error) {
		bookingVerifyDTO.IsDateAvailableForLaunch, err = u.launchpadRepo.IsDateAvailableForLaunch(b.LaunchpadID(), b.LaunchDate())
		return err
	})
	g.Go(func() (err error) {
		bookingVerifyDTO.IsLandpadExists, err = u.landpadRepo.IsExists(b.DestinationID())
		return err
	})
	g.Go(func() (err error) {
		bookingVerifyDTO.IsLandpadActive, err = u.landpadRepo.IsActive(b.DestinationID())
		return err
	})

	if err = g.Wait(); err != nil {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("Unknown reason") // hide errors from the adapters layer
		return err
	}

	err = u.bookingVerifierSvc.Verify(b, bookingVerifyDTO)
	if err != nil {
		return err
	}

	return u.bookingRepo.Save(b)
}
