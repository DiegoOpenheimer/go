package webservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DiegoOpenheimer/go/opentelemetry/configs"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk/utils_http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"log"
	"net/http"
)

func StartServiceA() {
	cfg := configs.GetConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/", handlerServiceA)
	fmt.Println("Server A running on port", cfg.ServiceAPort)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServiceAPort), r))
}

type ZipCodeRequest struct {
	ZipCode string `json:"cep" validate:"required,len=8"`
}

func handlerServiceA(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := pgk.Tracer.Start(ctx, "request_service_b")
	defer span.End()

	cfg := configs.GetConfig()
	var request ZipCodeRequest
	if err := utils_http.ParseBody(r, &request); err != nil {
		handlerError(w, errors.New("invalid zipcode"), http.StatusUnprocessableEntity)
		return
	}
	if err := validator.New().Struct(request); err != nil {
		handlerError(w, errors.New("invalid zipcode"), http.StatusUnprocessableEntity)
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", cfg.ServiceBUrl, request.ZipCode), nil)
	if err != nil {
		handlerError(w, err, http.StatusInternalServerError)
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		handlerError(w, err, http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		handlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func handlerError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	response := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
