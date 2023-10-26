package main

import (
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
	"time"
)

func Test_Bookings_POST__ValidationError(t *testing.T) {
	t.Parallel()
	cfg := getConfig()
	e := httpexpect.Default(t, cfg.ServerURL)

	e.POST("/v1/bookings").
		WithJSON(map[string]string{}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Value("error").NotNull()
}

func Test_Bookings_POST__Success(t *testing.T) {
	t.Parallel()
	cfg := getConfig()
	e := httpexpect.Default(t, cfg.ServerURL)

	req := map[string]string{
		"firstName":     faker.FirstName(),
		"lastName":      faker.LastName(),
		"gender":        faker.Gender(),
		"birthday":      "1982-10-27",
		"launchpadId":   "5e9e4501f509094ba4566f84",
		"destinationId": "5e9e3032383ecb761634e7cb",
		"launchDate":    time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}
	resp := e.POST("/v1/bookings").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	for _, k := range []string{"id", "status", "statusReason"} {
		resp.ContainsKey(k)
	}
	for k, v := range req {
		resp.Value(k).IsEqual(v)
	}
}

func Test_Bookings_Full_Flow__Success(t *testing.T) {
	t.Parallel()
	cfg := getConfig()
	e := httpexpect.Default(t, cfg.ServerURL)

	// Create booking
	req := map[string]string{
		"firstName":     faker.FirstName(),
		"lastName":      faker.LastName(),
		"gender":        faker.Gender(),
		"birthday":      "1982-10-27",
		"launchpadId":   "5e9e4501f509094ba4566f84",
		"destinationId": "5e9e3032383ecb761634e7cb",
		"launchDate":    time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	}
	booking := e.POST("/v1/bookings").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	for _, k := range []string{"id", "status", "statusReason"} {
		booking.ContainsKey(k)
	}
	for k, v := range req {
		booking.Value(k).IsEqual(v)
	}

	// Get bookings
	e.GET("/v1/bookings").
		WithQuery("limit", 1000).
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		Find(func(index int, value *httpexpect.Value) bool {
			return value.Object().Value("id").Raw() == booking.Value("id").Raw()
		}).Object().ContainsSubset(booking.Raw())

	// Delete booking
	e.DELETE("/v1/bookings/{id}").WithPath("id", booking.Value("id").Raw()).
		Expect().
		Status(http.StatusNoContent)

	// Verify that the booking has been deleted successfully
	e.GET("/v1/bookings").
		WithQuery("limit", 1000).
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		NotFind(func(index int, value *httpexpect.Value) bool {
			return value.Object().Value("id").Raw() == booking.Value("id").Raw()
		})
}
