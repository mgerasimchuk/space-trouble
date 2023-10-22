package usecase

import (
	"github.com/mgerasimchuk/space-trouble/internal/domain/service"
	"github.com/pkg/errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
	"github.com/mgerasimchuk/space-trouble/pkg/utils"

	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestBookingUsecase_VerifyBooking(t *testing.T) {
	validBooking := model.CreateBooking(
		"John", "Doe", "male", time.Now().Add(-300*time.Hour),
		"5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour),
	)

	// different types needs to make tests cases more readable
	type launchpadIsExistsBollErrRes struct {
		res bool
		err error
	}
	type launchpadIsActiveBollErrRes struct {
		res bool
		err error
	}
	type launchpadIsDateAvailableForLaunchBollErrRes struct {
		res bool
		err error
	}
	type landpadIsExistsBollErrRes struct {
		res bool
		err error
	}
	type landpadIsActiveBollErrRes struct {
		res bool
		err error
	}

	type args struct {
		booking                                  *model.Booking
		bookingRepoSaveRes                       error
		launchpadRepoIsExistsRes                 launchpadIsExistsBollErrRes
		launchpadRepoIsActiveRes                 launchpadIsActiveBollErrRes
		launchpadRepoIsDateAvailableForLaunchRes launchpadIsDateAvailableForLaunchBollErrRes
		landpadRepoIsExistsRes                   landpadIsExistsBollErrRes
		landpadRepoIsActiveRes                   landpadIsActiveBollErrRes
	}
	type want struct {
		status       string
		statusReason string
		err          error
	}
	tests := []struct {
		args args
		want want
	}{
		// Positive cases
		{
			args{
				validBooking, nil,
				launchpadIsExistsBollErrRes{true, nil}, launchpadIsActiveBollErrRes{true, nil},
				launchpadIsDateAvailableForLaunchBollErrRes{true, nil},
				landpadIsExistsBollErrRes{true, nil}, landpadIsActiveBollErrRes{true, nil},
			},
			want{model.StatusApproved, "", nil},
		},

		// Negative cases
		{
			args{
				validBooking, nil,
				launchpadIsExistsBollErrRes{false, errors.New("no such host")}, launchpadIsActiveBollErrRes{true, nil},
				launchpadIsDateAvailableForLaunchBollErrRes{true, nil},
				landpadIsExistsBollErrRes{true, nil}, landpadIsActiveBollErrRes{true, nil},
			},
			want{model.StatusDeclined, "Unknown reason", errors.New("no such host")},
		},
		{
			args{
				validBooking, nil,
				launchpadIsExistsBollErrRes{false, nil}, launchpadIsActiveBollErrRes{true, nil},
				launchpadIsDateAvailableForLaunchBollErrRes{true, nil},
				landpadIsExistsBollErrRes{true, nil}, landpadIsActiveBollErrRes{true, nil},
			},
			want{model.StatusDeclined, "launchpad doesn't exists", errors.New("launchpad doesn't exists")},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotBookingInSave := model.Booking{}
			bookingRepo := mock.NewMockBookingRepository(gomock.NewController(t))
			bookingRepo.EXPECT().Save(gomock.Any()).Do(func(b *model.Booking) {
				gotBookingInSave = *b
			}).Return(tt.args.bookingRepoSaveRes).AnyTimes()
			b := *tt.args.booking
			bookingRepo.EXPECT().GetAndModify(gomock.Any(), gomock.Any()).Return(&b, nil).AnyTimes()

			launchpadRepo := mock.NewMockLaunchpadRepository(gomock.NewController(t))
			launchpadRepo.EXPECT().IsExists(tt.args.booking.LaunchpadID()).
				Return(tt.args.launchpadRepoIsExistsRes.res, tt.args.launchpadRepoIsExistsRes.err).AnyTimes()
			launchpadRepo.EXPECT().IsActive(tt.args.booking.LaunchpadID()).
				Return(tt.args.launchpadRepoIsActiveRes.res, tt.args.launchpadRepoIsActiveRes.err).AnyTimes()
			launchpadRepo.EXPECT().IsDateAvailableForLaunch(tt.args.booking.LaunchpadID(), tt.args.booking.LaunchDate()).
				Return(tt.args.launchpadRepoIsDateAvailableForLaunchRes.res, tt.args.launchpadRepoIsDateAvailableForLaunchRes.err).AnyTimes()

			landpadRepo := mock.NewMockLaunchpadRepository(gomock.NewController(t))
			landpadRepo.EXPECT().IsExists(tt.args.booking.DestinationID()).
				Return(tt.args.landpadRepoIsExistsRes.res, tt.args.landpadRepoIsExistsRes.err).AnyTimes()
			landpadRepo.EXPECT().IsActive(tt.args.booking.DestinationID()).
				Return(tt.args.landpadRepoIsActiveRes.res, tt.args.landpadRepoIsActiveRes.err).AnyTimes()

			bookingVerifierSvc := service.NewBookingVerifierService()
			bookingUC := NewBookingUsecase(bookingVerifierSvc, bookingRepo, launchpadRepo, landpadRepo)
			gotErr := bookingUC.VerifyFirstAvailableBooking()

			utils.EqualError(t, tt.want.err, gotErr)
			assert.Equal(t, tt.want.status, b.Status())
			assert.Equal(t, tt.want.statusReason, b.StatusReason())
			if tt.want.err == nil {
				assert.Equal(t, tt.want.status, gotBookingInSave.Status())
				assert.Equal(t, tt.want.statusReason, gotBookingInSave.StatusReason())
			}
		})
	}
}
