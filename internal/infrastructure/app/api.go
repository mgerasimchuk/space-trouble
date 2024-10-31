package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/controller"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/api"
	"github.com/mgerasimchuk/space-trouble/internal/adapter/repository/pg"
	"github.com/mgerasimchuk/space-trouble/internal/entity/service"
	"github.com/mgerasimchuk/space-trouble/internal/infrastructure/config"
	"github.com/mgerasimchuk/space-trouble/internal/usecase"
	"github.com/mgerasimchuk/space-trouble/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func StartAPIApp(cfg *config.RootConfig) {
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
	bookingController := controller.NewBookingController(bookingUsecase, logger)

	gin.SetMode(cfg.HTTPServer.Mode)
	router := gin.Default()
	router.Use(ginlogrus.Logger(logger), gin.Recovery())
	router.Use(util.RequestLogger(logger))

	router.POST("/v1/bookings", bookingController.CreateBooking)
	router.GET("/v1/bookings", bookingController.GetBookings)
	router.DELETE("/v1/bookings/:id", bookingController.DeleteBooking)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTPServer.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	logger.Infof("Application has been started")

	// Wait for interrupting signal to gracefully shutdown the server with  a 5 seconds timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	if err := db.Close(); err != nil {
		logger.Fatal("Error during db connection close:", err)
	}

	logger.Info("Service stopped")
}
