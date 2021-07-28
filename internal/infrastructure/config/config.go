package config

import (
	"os"

	"github.com/spf13/viper"
)

type (
	LogConfig struct {
		Level int
	}
	DBConfig struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
	}
	LaunchpadAPIConfig struct {
		ApiBaseUri string
	}
	LandpadAPIConfig struct {
		ApiBaseUri string
	}
	HTTPServerConfig struct {
		Port int
		Mode string
	}
	Verifier struct {
		WorkersCount                int
		RunWorkersEveryMilliseconds int
	}

	RootConfig struct {
		Log        LogConfig
		DB         DBConfig
		Launchpad  LaunchpadAPIConfig
		Landpad    LandpadAPIConfig
		HTTPServer HTTPServerConfig
		Verifier   Verifier
	}
)

func GetRootConfig() RootConfig {
	viper.SetConfigFile(os.Getenv("CONFIG_FILE"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := RootConfig{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
