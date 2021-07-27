package main

import (
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/api"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/pg"
	"github.com/mgerasimchuk/space-trouble/internal/domain/model"
	"github.com/mgerasimchuk/space-trouble/internal/domain/repository"
	"github.com/mgerasimchuk/space-trouble/internal/domain/service"
	"github.com/mgerasimchuk/space-trouble/internal/usecase"
	"github.com/sirupsen/logrus"
)

// for playing with components, run from root of repo:
// GO111MODULE=on go run -mod vendor ${PWD}/cmd/development/main.go
func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "postgres",
	)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	bookingRepo := pg.NewBookingRepository(db)
	_ = bookingRepo

	launchpadRepo := api.NewLaunchpadRepository("https://api.spacexdata.com")
	_ = launchpadRepo

	landpadRepo := api.NewLandpadRepository("https://api.spacexdata.com")
	_ = landpadRepo

	bookingService := service.NewBookingService(bookingRepo, launchpadRepo, landpadRepo)

	bookingUsecase := usecase.NewBookingUsecase(bookingService, bookingRepo, logger)

	//examplesBookingRepo(bookingRepo)
	//examplesLaunchpadRepo(launchpadRepo)
	//examplesLandpadRepo(landpadRepo)
	examplesBookingUsecase(bookingUsecase)
}

func examplesBookingRepo(bookingRepo repository.BookingRepository) {
	p := gofakeit.Person()
	b1 := model.CreateBooking(p.FirstName, p.LastName, p.Gender, gofakeit.Date(), uuid.New().String(), uuid.New().String(), gofakeit.Date())
	b1, err := bookingRepo.Create(b1)
	fmt.Printf("After Create First\nBooking: %#v\nError: %#v\n\n", b1, err)

	b1.SetStatus(model.StatusDeclined)
	err = bookingRepo.Save(b1)
	fmt.Printf("After Update First\nBooking: %#v\nError: %#v\n\n", b1, err)

	bookings, err := bookingRepo.GetList(nil, nil, nil)
	fmt.Printf("Bookings: %#v\nError: %#v\n\n", bookings, err)

	p = gofakeit.Person()
	b2 := model.CreateBooking(p.FirstName, p.LastName, p.Gender, gofakeit.Date(), uuid.New().String(), uuid.New().String(), gofakeit.Date())
	b2, err = bookingRepo.Create(b2)
	fmt.Printf("After Create Second\nBooking: %#v\nError: %#v\n\n", b2, err)

	b2, err = bookingRepo.GetAndModify(model.StatusCreated, model.StatusPending)
	fmt.Printf("After GetAndModify\nBooking: %#v\nError: %#v\n\n", b2, err)

	bookings, err = bookingRepo.GetList(nil, nil, nil)
	fmt.Printf("Bookings (all): %#v\nError: %#v\n\n", bookings, err)

	status := model.StatusPending
	bookings, err = bookingRepo.GetList(&status, nil, nil)
	fmt.Printf("Bookings (pending): %#v\nError: %#v\n\n", bookings, err)

	_ = bookingRepo.Delete(b1.ID())
	bookings, err = bookingRepo.GetList(nil, nil, nil)
	fmt.Printf("After Delete First\nBookings: %#v\nError: %#v\n\n", bookings, err)

	_ = bookingRepo.Delete(b2.ID())
	bookings, err = bookingRepo.GetList(nil, nil, nil)
	fmt.Printf("After Delete Second\nBookings: %#v\nError: %#v\n\n", bookings, err)
}

func examplesLaunchpadRepo(launchpadRepo repository.LaunchpadRepository) {
	for _, id := range []string{
		"5e9e4501f5090910d4566f83-BAD", // bad
		"5e9e4501f5090910d4566f83",     // retired
		"5e9e4502f509092b78566f87",     // active - booked for "2021-12-31"
		"5e9e4501f509094ba4566f84",     // active - available for "2021-12-31"
	} {
		res, err := launchpadRepo.IsExists(id)
		fmt.Printf("Launchpad ID: %#v IsExists: %#v Error: %#v\n", id, res, err)

		res, err = launchpadRepo.IsActive(id)
		fmt.Printf("Launchpad ID: %#v IsActive: %#v Error: %#v\n", id, res, err)

		d := "2021-12-31"
		t, _ := time.Parse("2006-01-02", d)
		res, err = launchpadRepo.IsDateAvailableForLaunch(id, t)
		fmt.Printf("Launchpad ID: %#v Date: %#v IsDateAvailableForLaunch: %#v Error: %#v\n\n", id, d, res, err)
	}
}

func examplesLandpadRepo(landpadRepo repository.LandpadRepository) {
	for _, id := range []string{
		"5e9e3032383ecb267a34e7c7-BAD", // bad
		"5e9e3032383ecb267a34e7c7",     // retired
		"5e9e3032383ecb761634e7cb",     // active
	} {
		res, err := landpadRepo.IsExists(id)
		fmt.Printf("Landpad ID: %#v IsExists: %#v Error: %#v\n", id, res, err)

		res, err = landpadRepo.IsActive(id)
		fmt.Printf("Landpad ID: %#v IsActive: %#v Error: %#v\n\n", id, res, err)
	}
}

func examplesBookingUsecase(bookingUsecase *usecase.BookingUsecase) {
	p := gofakeit.Person()

	launchDate, _ := time.Parse("2006-01-02", "2021-12-31")
	firstBooking, err := bookingUsecase.CreateBooking(
		p.FirstName, p.LastName, p.Gender, gofakeit.DateRange(time.Now().Add(-900*time.Hour), time.Now().Add(-300*time.Hour)),
		"5e9e4501f509094ba4566f84", "5e9e3032383ecb761634e7cb", launchDate,
	)
	fmt.Printf("CreateBooking\nBooking: %#v\nError: %#v\n\n", firstBooking, err)

	limit, offset := 20, 0
	bookings, err := bookingUsecase.GetBookings(&limit, &offset)
	fmt.Printf("GetBookings\nBookings: %#v\nError: %#v\n\n", bookings, err)

	bookingUsecase.VerifyFirstAvailableBooking()

	bookings, err = bookingUsecase.GetBookings(&limit, &offset)
	fmt.Printf("GetBookings after VerifyFirstAvailableBooking\nBookings:\n")
	for _, b := range bookings {
		fmt.Printf(" %#v\n", b)
	}
	fmt.Printf("Error: %#v\n\n", err)

	err = bookingUsecase.DeleteBooking(firstBooking.ID())
	fmt.Printf("DeleteBooking\nError: %#v\n\n", err)

	bookings, err = bookingUsecase.GetBookings(&limit, &offset)
	fmt.Printf("GetBookings\nBookings: %#v\nError: %#v\n\n", bookings, err)

	badBooking, err := bookingUsecase.CreateBooking(
		p.FirstName, p.LastName, p.Gender, gofakeit.DateRange(time.Now().Add(-900*time.Hour), time.Now().Add(-300*time.Hour)),
		"5e9e4501f5090910d4566f83", "5e9e3032383ecb267a34e7c7", launchDate,
	)
	fmt.Printf("CreateBooking\nBooking: %#v\nError: %#v\n\n", badBooking, err)
	bookingUsecase.VerifyFirstAvailableBooking()

	bookings, err = bookingUsecase.GetBookings(&limit, &offset)
	fmt.Printf("GetBookings after VerifyFirstAvailableBooking\nBookings:\n")
	for _, b := range bookings {
		fmt.Printf(" %#v\n", b)
	}
	fmt.Printf("Error: %#v\n\n", err)

	badBooking, err = bookingUsecase.CreateBooking(
		p.FirstName, p.LastName, p.Gender, gofakeit.DateRange(time.Now().Add(-900*time.Hour), time.Now().Add(-300*time.Hour)),
		"5e9e4501f5090910d4566f83-BAD", "5e9e3032383ecb267a34e7c7-BAD", launchDate,
	)
	fmt.Printf("CreateBooking\nBooking: %#v\nError: %#v\n\n", badBooking, err)
	bookingUsecase.VerifyFirstAvailableBooking()

	bookings, err = bookingUsecase.GetBookings(&limit, &offset)
	fmt.Printf("GetBookings after VerifyFirstAvailableBooking\nBookings:\n")
	for _, b := range bookings {
		fmt.Printf(" %#v\n", b)
	}
	fmt.Printf("Error: %#v\n\n", err)
}
