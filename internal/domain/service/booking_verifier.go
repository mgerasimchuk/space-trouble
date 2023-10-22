package service

import (
	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
	"github.com/mgerasimchuk/space-trouble/internal/domain/service/dto"
	"github.com/pkg/errors"
)

type BookingVerifierService struct {
}

func NewBookingVerifierService() *BookingVerifierService {
	return &BookingVerifierService{}
}

func (s BookingVerifierService) Verify(b *model.Booking, dto dto.BookingVerifyDTO) error {
	if !dto.IsLaunchpadExists {
		b.SetStatus(model.StatusDeclined)
		err := errors.New("launchpad doesn't exists")
		b.SetStatusReason(err.Error())
		return err
	}
	if !dto.IsLaunchpadActive {
		b.SetStatus(model.StatusDeclined)
		err := errors.New("launchpad is not active")
		b.SetStatusReason(err.Error())
		return err
	}
	if !dto.IsDateAvailableForLaunch {
		b.SetStatus(model.StatusDeclined)
		err := errors.New("booking date is not available")
		b.SetStatusReason(err.Error())
		return err
	}
	if !dto.IsLandpadExists {
		b.SetStatus(model.StatusDeclined)
		err := errors.New("landpad doesn't exists")
		b.SetStatusReason(err.Error())
		return err
	}
	if !dto.IsLandpadActive {
		b.SetStatus(model.StatusDeclined)
		err := errors.New("landpad is not active")
		b.SetStatusReason(err.Error())
		return err
	}

	b.SetStatus(model.StatusApproved)
	b.SetStatusReason("")
	return nil
}
