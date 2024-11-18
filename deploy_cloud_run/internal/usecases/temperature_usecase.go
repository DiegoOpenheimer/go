package usecases

import "github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/ports"

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
	Code string
}

type TemperatureUseCaseOutput struct {
	Tempc float64 `json:"temp_C"`
	Tempf float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTemperatureUseCase = UseCase[TemperatureUseCaseInput, *TemperatureUseCaseOutput]

func (t *GetTemperature) Execute(input TemperatureUseCaseInput) (*TemperatureUseCaseOutput, error) {
	zipCode, err := t.zipCodeService.Get(input.Code)
	if err != nil {
		return nil, err
	}
	temperatureResponse, err := t.temperatureService.GetByCity(zipCode.City)
	if err != nil {
		return nil, err
	}
	return &TemperatureUseCaseOutput{
		Tempc: temperatureResponse.Current.TempC,
		Tempf: temperatureResponse.Current.TempF,
		TempK: temperatureResponse.Current.TempC + ValueToConvertCelsiusToKelvin,
	}, nil
}
