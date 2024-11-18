package main

import (
	"fmt"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/configs"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/infra/web/handlers"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/pgk/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	temperatureHandler := handlers.NewTemperatureHandler(
		usecases.NewTemperatureUseCase(
			services.NewZipCodeApi(),
			services.NewTemperatureApi(),
		),
	)

	r.Get("/{code}", temperatureHandler.GetTemperature)
	fmt.Println("Server running on port", cfg.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
