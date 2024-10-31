package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/gavv/httpexpect/v2"
	"testing"
)

type config struct {
	ServerURL string `env:"SERVER_URL" default:"http://localhost:8080"`
}

func getConfig() *config {
	cfg, err := env.ParseAsWithOptions[config](env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic("configuration error: " + err.Error())
	}
	return &cfg
}

func getExpect(t *testing.T) *httpexpect.Expect {
	t.Parallel()
	cfg := getConfig()
	return httpexpect.Default(t, cfg.ServerURL)
}
