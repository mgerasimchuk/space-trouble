package service

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
	"github.com/mgerasimchuk/space-trouble/internal/domain/repository"
)

type BookingService struct {
	bookingRepo   repository.BookingRepository
	launchpadRepo repository.LaunchpadRepository
	landpadRepo   repository.LandpadRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, launchpadRepo repository.LaunchpadRepository, landpadRepo repository.LandpadRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo, launchpadRepo: launchpadRepo, landpadRepo: landpadRepo}
}

func (s BookingService) VerifyBooking(b *model.Booking) error {
	if err := b.Validate(); err != nil {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason(err.Error())

		return s.bookingRepo.Save(b)
	}

	var isLaunchpadExists, isLaunchpadActive, isDateAvailableForLaunch, isLandpadExists, isLandpadActive bool
	var isLaunchpadExistsErr, isLaunchpadActiveErr, isDateAvailableForLaunchErr, isLandpadExistsErr, isLandpadActiveErr error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		isLaunchpadExists, isLaunchpadExistsErr = s.launchpadRepo.IsExists(b.LaunchpadID())
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		isLaunchpadActive, isLaunchpadActiveErr = s.launchpadRepo.IsActive(b.LaunchpadID())
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		isDateAvailableForLaunch, isDateAvailableForLaunchErr = s.launchpadRepo.IsDateAvailableForLaunch(b.LaunchpadID(), b.LaunchDate())
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		isLandpadExists, isLandpadExistsErr = s.landpadRepo.IsExists(b.DestinationID())
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		isLandpadActive, isLandpadActiveErr = s.landpadRepo.IsExists(b.DestinationID())
	}()
	wg.Wait()

	for _, err := range []error{
		isLaunchpadExistsErr, isLaunchpadActiveErr, isDateAvailableForLaunchErr, isLandpadExistsErr, isLandpadActiveErr,
	} {
		if err != nil {
			b.SetStatus(model.StatusDeclined)
			b.SetStatusReason("Unknown reason")
			saveErr := s.bookingRepo.Save(b)
			if saveErr != nil {
				err = errors.Wrap(err, saveErr.Error())
			}

			return err
		}
	}

	if !isLaunchpadExists {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("launchpad doesn't exists")

		return s.bookingRepo.Save(b)
	}
	if !isLaunchpadActive {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("launchpad is not active")

		return s.bookingRepo.Save(b)
	}
	if !isDateAvailableForLaunch {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("booking date is not available")

		return s.bookingRepo.Save(b)
	}
	if !isLandpadExists {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("landpad doesn't exists")

		return s.bookingRepo.Save(b)
	}
	if !isLandpadActive {
		b.SetStatus(model.StatusDeclined)
		b.SetStatusReason("landpad is not active")

		return s.bookingRepo.Save(b)
	}

	b.SetStatus(model.StatusApproved)
	b.SetStatusReason("")

	return s.bookingRepo.Save(b)
}
