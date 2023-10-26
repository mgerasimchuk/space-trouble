package controller

import (
	"github.com/mgerasimchuk/space-trouble/internal/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mgerasimchuk/space-trouble/internal/usecase"
	"github.com/mgerasimchuk/space-trouble/pkg/util"
	"github.com/sirupsen/logrus"
)

type BookingController struct {
	bookingUsecase *usecase.BookingUsecase
	logger         *logrus.Logger
}

func NewBookingController(bookingUsecase *usecase.BookingUsecase, logger *logrus.Logger) *BookingController {
	return &BookingController{bookingUsecase: bookingUsecase, logger: logger}
}

func (c *BookingController) CreateBooking(ctx *gin.Context) {
	var b booking
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdDomainBooking, err := c.bookingUsecase.CreateBooking(
		b.FirstName, b.LastName, b.Gender, b.Birthday.Time, b.LaunchpadID, b.DestinationID, b.LaunchDate.Time,
	)

	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdBooking := fromDomainBooking(createdDomainBooking)

	ctx.JSON(http.StatusCreated, createdBooking)
}

func (c *BookingController) GetBookings(ctx *gin.Context) {
	var limit, offset *int

	if ctx.Query("limit") != "" {
		if l, err := strconv.Atoi(ctx.Query("limit")); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "\"limit\" param should be a number"})
			return
		} else {
			limit = &l
		}
	}

	if ctx.Query("offset") != "" {
		if o, err := strconv.Atoi(ctx.Query("offset")); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "\"offset\" param should be a number"})
			return
		} else {
			offset = &o
		}
	}

	domainBookings, err := c.bookingUsecase.GetBookings(limit, offset)

	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookings []*booking
	bookings = make([]*booking, 0) // for fixing "null" response body

	for i := range domainBookings {
		bookings = append(bookings, fromDomainBooking(domainBookings[i]))
	}

	ctx.JSON(http.StatusOK, bookings)
}

func (c *BookingController) DeleteBooking(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "\"id\" param can't be empty"})
		return
	}

	err := c.bookingUsecase.DeleteBooking(id)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

type booking struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	StatusReason  string    `json:"statusReason"`
	FirstName     string    `json:"firstName" binding:"required"`
	LastName      string    `json:"lastName" binding:"required"`
	Gender        string    `json:"gender" binding:"required"`
	Birthday      util.Date `json:"birthday" binding:"required"`
	LaunchpadID   string    `json:"launchpadId" binding:"required"`
	DestinationID string    `json:"destinationId" binding:"required"`
	LaunchDate    util.Date `json:"launchDate" binding:"required"`
}

func fromDomainBooking(b *entity.Booking) *booking {
	return &booking{
		ID:            b.ID(),
		Status:        b.Status(),
		StatusReason:  b.StatusReason(),
		FirstName:     b.FirstName(),
		LastName:      b.LastName(),
		Gender:        b.Gender(),
		Birthday:      util.Date{Time: b.Birthday()},
		LaunchpadID:   b.LaunchpadID(),
		DestinationID: b.DestinationID(),
		LaunchDate:    util.Date{Time: b.LaunchDate()},
	}
}
