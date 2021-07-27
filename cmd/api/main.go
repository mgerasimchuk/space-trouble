package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgerasimchuk/space-trouble/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel) // TODO move to config

	gin.SetMode(gin.DebugMode) // TODO move to config
	router := gin.Default()
	router.Use(ginlogrus.Logger(logger), gin.Recovery())
	router.Use(utils.RequestLogger(logger))

	router.POST("/bookings", func(ctx *gin.Context) { ctx.Status(http.StatusNotImplemented) })
	router.GET("/bookings", func(ctx *gin.Context) { ctx.Status(http.StatusNotImplemented) })
	router.DELETE("/bookings/:id", func(ctx *gin.Context) { ctx.Status(http.StatusNotImplemented) })

	srv := &http.Server{
		Addr:    ":8080", // TODO move to config
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupting signal to gracefully shutdown the server with 5 seconds timeout
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Service stopped")
}
