package pg

import (
	"fmt"
	"github.com/mgerasimchuk/space-trouble/internal/usecase/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
)

type BookingRepository struct {
	db *gorm.DB
}

const tableName = "bookings"

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(domainBooking *model.Booking) (*model.Booking, error) {
	domainBooking.SetID(uuid.New().String())
	b := fromDomainBooking(domainBooking)

	return domainBooking, r.db.Create(b).Error
}

func (r *BookingRepository) Save(domainBooking *model.Booking) error {
	b := fromDomainBooking(domainBooking)

	return r.db.Save(b).Error
}

func (r *BookingRepository) GetList(status *string, limit, offset *int) ([]*model.Booking, error) {
	var bookings []*booking

	query := r.db
	if status != nil {
		query = query.Where("status = ?", status)
	}

	if limit == nil {
		limit = new(int)
		*limit = repository.DefaultBookingGetListLimit
	}
	query = query.Limit(*limit)

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}

	var domainBookings []*model.Booking

	for i := range bookings {
		domainBookings = append(domainBookings, toDomainBooking(bookings[i]))
	}

	return domainBookings, nil
}

func (r *BookingRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&booking{}).Error
}

func (r *BookingRepository) GetAndModify(searchStatus, modifyStatus string) (*model.Booking, error) {
	res := r.db.Raw(fmt.Sprintf(`
UPDATE %s
SET status = '%s'
    WHERE  id = (
	SELECT id
	FROM %s
	WHERE status = '%s'
	LIMIT 1
	FOR UPDATE SKIP LOCKED
)
RETURNING *;
`, tableName, modifyStatus, tableName, searchStatus))

	if res.Error != nil {
		return nil, res.Error
	}

	b := &booking{}
	res.Scan(b)

	if b.ID == "" {
		return nil, nil
	}

	return toDomainBooking(b), nil
}

type booking struct {
	ID            string    `gorm:"primary_key;column:id"`
	Status        string    `gorm:"column:status"`
	StatusReason  string    `gorm:"column:status_reason"`
	FirstName     string    `gorm:"column:first_name"`
	LastName      string    `gorm:"column:last_name"`
	Gender        string    `gorm:"column:gender"`
	Birthday      time.Time `gorm:"column:birthday"`
	LaunchpadID   string    `gorm:"column:launchpad_id"`
	DestinationID string    `gorm:"column:destination_id"`
	LaunchDate    time.Time `gorm:"column:launch_date"`
}

func fromDomainBooking(b *model.Booking) *booking {
	return &booking{
		ID:            b.ID(),
		Status:        b.Status(),
		StatusReason:  b.StatusReason(),
		FirstName:     b.FirstName(),
		LastName:      b.LastName(),
		Gender:        b.Gender(),
		Birthday:      b.Birthday(),
		LaunchpadID:   b.LaunchpadID(),
		DestinationID: b.DestinationID(),
		LaunchDate:    b.LaunchDate(),
	}
}

func toDomainBooking(b *booking) *model.Booking {
	return model.LoadBooking(
		b.ID, b.Status, b.StatusReason, b.FirstName, b.LastName, b.Gender,
		b.Birthday, b.LaunchpadID, b.DestinationID, b.LaunchDate,
	)
}
