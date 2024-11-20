package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases"
	errorsUseCase "github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/errors"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type TemperatureHandler struct {
	useCase usecases.GetTemperatureUseCase
}

func NewTemperatureHandler(temperatureUseCase usecases.GetTemperatureUseCase) *TemperatureHandler {
	return &TemperatureHandler{
		temperatureUseCase,
	}
}

func (t *TemperatureHandler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	input := chi.URLParam(r, "code")
	resp, err := t.useCase.Execute(usecases.TemperatureUseCaseInput{Code: input, Context: ctx})
	var useCaseError errorsUseCase.UseCaseError
	if errors.As(err, &useCaseError) {
		w.WriteHeader(useCaseError.GetCode())
		response := struct {
			Message string `json:"message"`
		}{
			Message: useCaseError.Error(),
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
