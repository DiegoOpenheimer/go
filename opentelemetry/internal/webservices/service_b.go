package webservices

import (
	"fmt"
	"github.com/DiegoOpenheimer/go/opentelemetry/configs"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/infra/web/handlers"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func StartServiceB() {
	cfg := configs.GetConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	temperatureHandler := handlers.NewTemperatureHandler(
		usecases.NewTemperatureUseCase(
			services.NewZipCodeApi(),
			services.NewTemperatureApi(),
		),
	)

	r.Get("/{code}", temperatureHandler.GetTemperature)
	fmt.Println("Server B running on port", cfg.ServiceBPort)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServiceBPort), r))
}
