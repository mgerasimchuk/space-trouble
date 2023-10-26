package main

import (
	"github.com/kelseyhightower/envconfig"
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
