package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	LogConfig struct {
		Level int `envconfig:"LEVEL" required:"true"`
	}
	DBConfig struct {
		Host     string `envconfig:"HOST" required:"true"`
		Port     int    `envconfig:"PORT" required:"true"`
		Name     string `envconfig:"NAME" required:"true"`
		User     string `envconfig:"USER" required:"true"`
		Password string `envconfig:"PASSWORD" required:"true"`
	}
	APIClientConfig struct {
		ApiBaseUri string `envconfig:"API_BASE_URI" required:"true"`
	}
	HTTPServerConfig struct {
		Port int    `envconfig:"PORT" required:"true"`
		Mode string `envconfig:"MODE" required:"true"`
	}
	Verifier struct {
		WorkersCount                int `envconfig:"WORKERS_COUNT" required:"true"`
		RunWorkersEveryMilliseconds int `envconfig:"RUN_WORKERS_EVERY_MILLISECONDS" required:"true"`
	}

	RootConfig struct {
		Log        LogConfig        `envconfig:"LOG" required:"true"`
		DB         DBConfig         `envconfig:"DB" required:"true"`
		Launchpad  APIClientConfig  `envconfig:"LAUNCHPAD" required:"true"`
		Landpad    APIClientConfig  `envconfig:"LANDPAD" required:"true"`
		HTTPServer HTTPServerConfig `envconfig:"HTTP_SERVER" required:"true"`
		Verifier   Verifier         `envconfig:"VERIFIER" required:"true"`
	}
)

func GetRootConfig() RootConfig {
	cfg := &RootConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		panic("Configuration error: " + err.Error())
	}
	return *cfg
}
