package config

import "github.com/caarlos0/env/v11"

type (
	LogConfig struct {
		Level int `env:"LEVEL"`
	}
	DBConfig struct {
		Host     string `env:"HOST"`
		Port     int    `env:"PORT"`
		Name     string `env:"NAME"`
		User     string `env:"USER"`
		Password string `env:"PASSWORD"`
	}
	APIClientConfig struct {
		APIBaseURL string `env:"API_BASE_URL"`
	}
	HTTPServerConfig struct {
		Port int    `env:"PORT"`
		Mode string `env:"MODE"`
	}
	Verifier struct {
		WorkersCount                int `env:"WORKERS_COUNT"`
		RunWorkersEveryMilliseconds int `env:"RUN_WORKERS_EVERY_MILLISECONDS"`
	}

	RootConfig struct {
		Log        LogConfig        `envPrefix:"LOG_"`
		DB         DBConfig         `envPrefix:"DB_"`
		Launchpad  APIClientConfig  `envPrefix:"LAUNCHPAD_"`
		Landpad    APIClientConfig  `envPrefix:"LANDPAD_"`
		HTTPServer HTTPServerConfig `envPrefix:"HTTP_SERVER_"`
		Verifier   Verifier         `envPrefix:"VERIFIER_"`
	}
)

func GetRootConfig() *RootConfig {
	cfg, err := env.ParseAsWithOptions[RootConfig](env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic("Configuration error: " + err.Error())
	}
	return &cfg
}
