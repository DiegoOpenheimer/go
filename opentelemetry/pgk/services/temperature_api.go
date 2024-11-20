package services

import (
	"context"
	"encoding/json"
	"github.com/DiegoOpenheimer/go/opentelemetry/configs"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/errors"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/ports"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
)

type TemperatureApi struct{}

func NewTemperatureApi() *TemperatureApi {
	return &TemperatureApi{}
}

func (t *TemperatureApi) GetByCityWithContext(ctx context.Context, city string) (*ports.TemperatureResponse, error) {
	ctx, span := pgk.Tracer.Start(ctx, "request_weather_api")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weatherapi.com/v1/current.json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("KEY", configs.GetConfig().WeatherApiKey)
	q.Add("q", city)
	req.URL.RawQuery = q.Encode()
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewUnknownError()
	}
	if body, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		var temperatureResponse ports.TemperatureResponse
		if err = json.Unmarshal(body, &temperatureResponse); err != nil {
			return nil, err
		}
		return &temperatureResponse, nil
	}
}
