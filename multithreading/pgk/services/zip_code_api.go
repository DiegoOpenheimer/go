package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ZipCode struct {
	ZipCode      string `json:"cep"`
	Street       string `json:"logradouro"`
	Complement   string `json:"complemento"`
	Unit         string `json:"unidade"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	AreaCode     string `json:"ddd"`
	Siafi        string `json:"siafi"`
}

type ZipCodeApi struct{}

func (b *BrazilZipCode) JSON() string {
	result, _ := json.MarshalIndent(&b, "", "  ")
	return string(result)
}

func (z ZipCodeApi) GetZipCodeWithContext(ctx context.Context, zipCode string) (*ZipCode, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipCode), nil)
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
		var zipCodeResponse ZipCode
		if err = json.Unmarshal(body, &zipCodeResponse); err != nil {
			return nil, err
		}
		return &zipCodeResponse, nil
	}
}
