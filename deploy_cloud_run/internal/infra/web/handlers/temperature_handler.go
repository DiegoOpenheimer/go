package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases"
	errorsUseCase "github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/errors"
	"github.com/go-chi/chi/v5"
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
	input := chi.URLParam(r, "code")
	resp, err := t.useCase.Execute(usecases.TemperatureUseCaseInput{Code: input})
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
