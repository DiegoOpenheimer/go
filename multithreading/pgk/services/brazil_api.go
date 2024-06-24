package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DiegoOpenheimer/go/multithreading/pgk/models"
	"io"
	"net/http"
)

type BrazilZipCode struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type BrazilApi struct{}

func (b *BrazilZipCode) JSON() string {
	result, _ := json.MarshalIndent(&b, "", "  ")
	return string(result)
}

func (b BrazilApi) GetZipCodeWithContext(ctx context.Context, zipCode string) (models.ZipCodeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", zipCode), nil)
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
	if body, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		var zipCodeResponse BrazilZipCode
		if err = json.Unmarshal(body, &zipCodeResponse); err != nil {
			return nil, err
		}
		return &zipCodeResponse, nil
	}
}
