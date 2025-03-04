package main

import (
	"log"
	"net/http"

	"github.com/Oxygenss/yandex_final_project/internal/config"
	"github.com/Oxygenss/yandex_final_project/internal/handler"
	"github.com/Oxygenss/yandex_final_project/internal/repository"
	"github.com/Oxygenss/yandex_final_project/internal/service"
)

func main() {
	cfg := config.MustLoad()

	repository, err := repository.New(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	service := service.NewService(repository)
	handler := handler.NewHandler(*service, *cfg)

	router := handler.InitRoutes(*cfg)

	serve := cfg.Server.Host + ":" + cfg.Server.Port
	err = http.ListenAndServe(serve, router)
	if err != nil {
		log.Fatal(err)
	}

}
