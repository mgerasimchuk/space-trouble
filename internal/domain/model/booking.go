package model

import (
	"errors"
	"time"
)

const (
	StatusCreated  = "created"
	StatusPending  = "pending"
	StatusDeclined = "declined"
	StatusApproved = "approved"
)

type Booking struct {
	id            string
	status        string
	statusReason  string
	firstName     string
	lastName      string
	gender        string
	birthday      time.Time
	launchpadID   string
	destinationID string
	launchDate    time.Time
}

func CreateBooking(firstName, lastName, gender string, birthday time.Time, launchpadID, destinationID string, launchDate time.Time) *Booking {
	b := &Booking{
		firstName: firstName, lastName: lastName, gender: gender, birthday: birthday,
		launchpadID: launchpadID, destinationID: destinationID, launchDate: launchDate,
	}
	b.status = StatusCreated

	return b
}

func LoadBooking(id, status, statusReason, firstName, lastName, gender string, birthday time.Time, launchpadID, destinationID string, launchDate time.Time) *Booking {
	b := &Booking{
		id: id, status: status, statusReason: statusReason, firstName: firstName, lastName: lastName, gender: gender,
		birthday: birthday, launchpadID: launchpadID, destinationID: destinationID, launchDate: launchDate,
	}

	return b
}

func (b *Booking) ID() string {
	return b.id
}

func (b *Booking) SetID(id string) {
	b.id = id
}

func (b *Booking) Status() string {
	return b.status
}

func (b *Booking) SetStatus(status string) {
	b.status = status
}

func (b *Booking) StatusReason() string {
	return b.statusReason
}

func (b *Booking) SetStatusReason(statusReason string) {
	b.statusReason = statusReason
}

func (b *Booking) FirstName() string {
	return b.firstName
}

func (b *Booking) SetFirstName(firstName string) {
	b.firstName = firstName
}

func (b *Booking) LastName() string {
	return b.lastName
}

func (b *Booking) SetLastName(lastName string) {
	b.lastName = lastName
}

func (b *Booking) Gender() string {
	return b.gender
}

func (b *Booking) SetGender(gender string) {
	b.gender = gender
}

func (b *Booking) Birthday() time.Time {
	return b.birthday
}

func (b *Booking) SetBirthday(birthday time.Time) {
	b.birthday = birthday
}

func (b *Booking) LaunchpadID() string {
	return b.launchpadID
}

func (b *Booking) SetLaunchpadID(launchpadID string) {
	b.launchpadID = launchpadID
}

func (b *Booking) DestinationID() string {
	return b.destinationID
}

func (b *Booking) SetDestinationID(destinationID string) {
	b.destinationID = destinationID
}

func (b *Booking) LaunchDate() time.Time {
	return b.launchDate
}

func (b *Booking) SetLaunchDate(launchDate time.Time) {
	b.launchDate = launchDate
}

func (b Booking) Validate() error {
	if err := b.validateStatus(); err != nil {
		return err
	}
	if err := b.validateFirstName(); err != nil {
		return err
	}
	if err := b.validateLastName(); err != nil {
		return err
	}
	if err := b.validateGender(); err != nil {
		return err
	}
	if err := b.validateBirthday(); err != nil {
		return err
	}
	if err := b.validateLaunchpadID(); err != nil {
		return err
	}
	if err := b.validateDestinationID(); err != nil {
		return err
	}
	if err := b.validateLaunchDate(); err != nil {
		return err
	}

	return nil
}

func (b Booking) validateStatus() error {
	if b.Status() == "" {
		return errors.New("value of \"status\" field can't be empty")
	}

	return nil
}

func (b Booking) validateFirstName() error {
	if b.FirstName() == "" {
		return errors.New("value of \"firstName\" field can't be empty")
	}

	return nil
}

func (b Booking) validateLastName() error {
	if b.LastName() == "" {
		return errors.New("value of \"lastName\" field can't be empty")
	}

	return nil
}

func (b Booking) validateGender() error {
	if b.Gender() == "" {
		return errors.New("value of \"gender\" field can't be empty")
	}

	return nil
}

func (b Booking) validateBirthday() error {
	if b.Birthday().After(time.Now()) {
		// actually should be more than at least 18 years
		return errors.New("value of \"birthday\" field should be in the past")
	}

	return nil
}

func (b Booking) validateLaunchpadID() error {
	if b.LaunchpadID() == "" {
		return errors.New("value of \"launchpadID\" field can't be empty")
	}

	return nil
}

func (b Booking) validateDestinationID() error {
	if b.DestinationID() == "" {
		return errors.New("value of \"destinationID\" field can't be empty")
	}

	return nil
}

func (b Booking) validateLaunchDate() error {
	if b.LaunchDate().Before(time.Now()) {
		return errors.New("value of \"launchDate\" field should be in the future")
	}

	return nil
}
