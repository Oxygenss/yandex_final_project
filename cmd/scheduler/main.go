package main

import (
	"log"
	"net/http"

	"github.com/Oxygenss/yandex_final_project/internal/config"
	"github.com/Oxygenss/yandex_final_project/internal/handler"
	"github.com/Oxygenss/yandex_final_project/internal/service"
	repository "github.com/Oxygenss/yandex_final_project/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	repository, err := repository.New(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	service := service.New(repository)
	handler := handler.New(*service, *cfg)

	router := handler.InitRoutes()

	serve := cfg.Server.Host + ":" + cfg.Server.Port
	err = http.ListenAndServe(serve, router)
	if err != nil {
		log.Fatal(err)
	}

}
