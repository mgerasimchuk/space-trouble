package service

import (
	"github.com/mgerasimchuk/space-trouble/internal/entity"
	"github.com/mgerasimchuk/space-trouble/internal/entity/service/dto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBookingService_VerifyBooking(t *testing.T) {
	validBooking := entity.CreateBooking(
		"John", "Doe", "male", time.Now().Add(-300*time.Hour),
		"5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour),
	)
	//invalidBooking := model.CreateBooking(
	//	"", "Doe", "male", time.Now().Add(-300*time.Hour),
	//	"5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour),
	//)

	// different types needs to make tests cases more readable

	type args struct {
		booking                 *entity.Booking
		bookingServiceVerifyRes error
		bookingVerifyDTO        dto.BookingVerifyDTO
	}
	type want struct {
		status string
		err    error
	}
	tests := []struct {
		args args
		want want
	}{
		// Positive cases
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{true, true, true, true, true},
			},
			want{entity.StatusApproved, nil},
		},

		// Negative cases
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{false, true, true, true, true},
			},
			want{entity.StatusDeclined, errors.New("launchpad doesn't exists")},
		},
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{true, false, true, true, true},
			},
			want{entity.StatusDeclined, errors.New("launchpad is not active")},
		},
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{true, true, false, true, true},
			},
			want{entity.StatusDeclined, errors.New("booking date is not available")},
		},
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{true, true, true, false, true},
			},
			want{entity.StatusDeclined, errors.New("landpad doesn't exists")},
		},
		{
			args{
				validBooking, nil,
				dto.BookingVerifyDTO{true, true, true, true, false},
			},
			want{entity.StatusDeclined, errors.New("landpad is not active")},
		},
	}

	s := NewBookingVerifierService()

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b := *tt.args.booking
			gotErr := s.Verify(&b, tt.args.bookingVerifyDTO)

			assert.Equal(t, tt.want.status, b.Status())
			if tt.want.err != nil {
				assert.Error(t, gotErr)
				assert.EqualError(t, tt.want.err, gotErr.Error())
				assert.Equal(t, tt.want.err.Error(), b.StatusReason())
			}
		})
	}
}
