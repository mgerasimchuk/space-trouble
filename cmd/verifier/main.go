package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/api"
	"github.com/mgerasimchuk/space-trouble/internal/infrastructure/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"

	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/pg"
	"github.com/mgerasimchuk/space-trouble/internal/entity/service"
	"github.com/mgerasimchuk/space-trouble/internal/usecase"
)

const logDateTimeLayout = "2006-01-02 15:04:05"

func main() {
	cfg := config.GetRootConfig()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.Level(cfg.Log.Level))

	logger.Info("Application starting...")

	dbConnectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name, cfg.DB.Password,
	)
	db, err := gorm.Open("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}

	bookingRepo := pg.NewBookingRepository(db)
	launchpadRepo := api.NewLaunchpadRepository(cfg.Launchpad.APIBaseURL)
	landpadRepo := api.NewLandpadRepository(cfg.Landpad.APIBaseURL)
	bookingService := service.NewBookingVerifierService()
	bookingUsecase := usecase.NewBookingUsecase(bookingService, bookingRepo, launchpadRepo, landpadRepo)

	defer func() {
		if r := recover(); r != nil {
			logger.Panicf("App crashed & recovered with: %#v", r)
		}
	}()

	ticker := time.NewTicker(time.Duration(cfg.Verifier.RunWorkersEveryMilliseconds * int(time.Millisecond)))
	tickerDone := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerDone:
				return
			case t := <-ticker.C:
				logger.Debugf("[%s] Ticker triggered", t.Format(logDateTimeLayout))

				wg := sync.WaitGroup{}
				for i := 1; i < cfg.Verifier.WorkersCount; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err = bookingUsecase.VerifyFirstAvailableBooking()
						if err != nil {
							logger.Error(err)
						}
					}()
				}
				wg.Wait()
			}
		}
	}()

	logger.Infof("Application has been started")

	// Wait for interrupting signal to gracefully shutdown the server with a 5 seconds timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ticker.Stop()
	tickerDone <- true

	logger.Info("Ticker stopped")

	if err := db.Close(); err != nil {
		logger.Fatal("Error during db connection close:", err)
	}

	logger.Info("Service stopped")
}
