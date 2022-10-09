package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OksidGen/market-couriers_salary/config"
	"github.com/OksidGen/market-couriers_salary/internal/db"
	handler "github.com/OksidGen/market-couriers_salary/internal/handler/http"
	"github.com/OksidGen/market-couriers_salary/internal/infrastructure/repository"
	"github.com/OksidGen/market-couriers_salary/internal/infrastructure/webapi"
	"github.com/OksidGen/market-couriers_salary/internal/server"
	"github.com/OksidGen/market-couriers_salary/internal/usecase"
)

func Run(cfg *config.Config) {
	dbPostgres, err := db.NewPostgresClient(cfg.PG.URL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(dbPostgres); err != nil {
		log.Fatal(err)
	}

	repos := repository.New(dbPostgres)
	webapi := webapi.New(webapi.Config{
		TG: webapi.TG{
			Webhook:  cfg.TG.Webhook,
			Endpoint: cfg.TG.Endpoint,
			Token:    cfg.TG.Token,
		},
		Ya: webapi.Ya{
			ClientID:     cfg.YaApp.ClienID,
			ClientSecret: cfg.YaApp.ClienSecret,
			Host:         cfg.YaApp.Host,
		},
	})

	useCases := usecase.New(repos, webapi)

	handlers := handler.New(useCases)

	srv := server.NewServerHTTP(handlers.Init(), cfg.Server.Port)

	srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-quit:
		log.Println("app - Run - signal: " + s.String())
	case err = <-srv.Notify():
		log.Println("app - Run - httpServer.Notify: ", err)
	}

	err = srv.Shutdown()
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: ", err)
	}
}
