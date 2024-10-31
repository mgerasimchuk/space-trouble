package main

import (
	"github.com/mgerasimchuk/space-trouble/internal/infrastructure/app"
	"github.com/mgerasimchuk/space-trouble/internal/infrastructure/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	cfg := config.GetRootConfig()
	app.StartVerifierApp(cfg)
}
