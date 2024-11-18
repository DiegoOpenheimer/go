package services

import (
	"encoding/json"
	"fmt"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/errors"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases/ports"
	"io"
	"net/http"
)

type ZipCodeApi struct{}

func (z *ZipCodeApi) Get(code string) (*ports.ZipCodeResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", code), nil)
	if err != nil {
		return nil, err
	}
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
