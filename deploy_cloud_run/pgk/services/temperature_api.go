package services

import (
	"encoding/json"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/configs"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/errors"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/ports"
	"io"
	"net/http"
)

type TemperatureApi struct{}

func NewTemperatureApi() *TemperatureApi {
	return &TemperatureApi{}
}

func (t *TemperatureApi) GetByCity(city string) (*ports.TemperatureResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.weatherapi.com/v1/current.json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("KEY", configs.GetConfig().WeatherApiKey)
	q.Add("q", city)
	req.URL.RawQuery = q.Encode()
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
