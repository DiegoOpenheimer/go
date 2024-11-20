package usecases

import (
	"context"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/usecases/ports"
)

const (
	ValueToConvertCelsiusToKelvin = 273
)

type GetTemperature struct {
	zipCodeService     ports.ZipCodeService
	temperatureService ports.TemperatureService
}

func NewTemperatureUseCase(zipCodeService ports.ZipCodeService, temperatureService ports.TemperatureService) *GetTemperature {
	return &GetTemperature{zipCodeService: zipCodeService, temperatureService: temperatureService}
}

type TemperatureUseCaseInput struct {
	Code    string
	Context context.Context
}

type TemperatureUseCaseOutput struct {
	Tempc float64 `json:"temp_C"`
	Tempf float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
	City  string  `json:"city"`
}

type GetTemperatureUseCase = UseCase[TemperatureUseCaseInput, *TemperatureUseCaseOutput]

func (t *GetTemperature) Execute(input TemperatureUseCaseInput) (*TemperatureUseCaseOutput, error) {
	zipCode, err := t.zipCodeService.GetWithContext(input.Context, input.Code)
	if err != nil {
		return nil, err
	}
	temperatureResponse, err := t.temperatureService.GetByCityWithContext(input.Context, zipCode.City)
	if err != nil {
		return nil, err
	}
	return &TemperatureUseCaseOutput{
		Tempc: temperatureResponse.Current.TempC,
		Tempf: temperatureResponse.Current.TempF,
		TempK: temperatureResponse.Current.TempC + ValueToConvertCelsiusToKelvin,
		City:  zipCode.City,
	}, nil
}
