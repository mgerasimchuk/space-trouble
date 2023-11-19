package main

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/kelseyhightower/envconfig"
	"testing"
)

type config struct {
	ServerURL string `envconfig:"SERVER_URL" required:"true" default:"http://localhost:8080"`
}

func getConfig() *config {
	cfg := &config{}
	if err := envconfig.Process("", cfg); err != nil {
		panic("configuration error: " + err.Error())
	}
	return cfg
}

func getExpect(t *testing.T) *httpexpect.Expect {
	t.Parallel()
	cfg := getConfig()
	return httpexpect.Default(t, cfg.ServerURL)
}
