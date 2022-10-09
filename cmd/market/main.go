package main

import (
	"log"

	"github.com/OksidGen/market-couriers_salary/config"
	"github.com/OksidGen/market-couriers_salary/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
