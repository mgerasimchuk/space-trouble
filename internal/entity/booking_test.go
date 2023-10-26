package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestBooking_Validate(t *testing.T) {
	type args struct {
		booking *Booking
	}
	tests := []struct {
		args args
		want error
	}{
		// Positive cases
		{
			args{CreateBooking("John", "Doe", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			nil,
		},

		// Negative cases
		{
			args{LoadBooking("c7dca667-2753-4839-9da0-629bc606827f", "", "Unknown reason", "John", "Doe", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"status\" field can't be empty"),
		},
		{
			args{CreateBooking("", "Doe", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"firstName\" field can't be empty"),
		},
		{
			args{CreateBooking("John", "", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"lastName\" field can't be empty"),
		},
		{
			args{CreateBooking("John", "Doe", "", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"gender\" field can't be empty"),
		},
		{
			args{CreateBooking("John", "Doe", "male", time.Now().Add(+300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"birthday\" field should be in the past"),
		},
		{
			args{CreateBooking("John", "Doe", "male", time.Now().Add(-300*time.Hour), "", "5e9e3032383ecb267a34e7c7", time.Now().Add(300*time.Hour))},
			errors.New("value of \"launchpadID\" field can't be empty"),
		},
		{
			args{CreateBooking("John", "Doe", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "", time.Now().Add(300*time.Hour))},
			errors.New("value of \"destinationID\" field can't be empty"),
		},
		{
			args{CreateBooking("John", "Doe", "male", time.Now().Add(-300*time.Hour), "5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", time.Now().Add(-300*time.Hour))},
			errors.New("value of \"launchDate\" field should be in the future"),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.args.booking.Validate(); !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Booking.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
