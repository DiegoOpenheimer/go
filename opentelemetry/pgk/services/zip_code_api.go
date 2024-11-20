package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/errors"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/ports"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
)

type ZipCodeApi struct{}

func (z *ZipCodeApi) GetWithContext(ctx context.Context, code string) (*ports.ZipCodeResponse, error) {
	ctx, span := pgk.Tracer.Start(ctx, "request_cep")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", code), nil)
	if err != nil {
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, errors.NewInvalidZipCodeError()
	case http.StatusNotFound:
		return nil, errors.NewZipCodeNotFoundError()
	case http.StatusOK:
		break
	default:
		return nil, errors.NewUnknownError()
	}
	if body, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		var zipCodeResponse ports.ZipCodeResponse
		if err = json.Unmarshal(body, &zipCodeResponse); err != nil {
			return nil, err
		}
		if zipCodeResponse.Error == "true" {
			return nil, errors.NewZipCodeNotFoundError()
		}
		return &zipCodeResponse, nil
	}
}

func NewZipCodeApi() *ZipCodeApi {
	return &ZipCodeApi{}
}
