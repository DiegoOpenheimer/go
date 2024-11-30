package infra

import (
	"github.com/DiegoOpenheimer/go/stress-test/internal/usecase/ports"
	"io"
	"net/http"
)

type HttpWebService struct {
}

func NewHttpWebService() *HttpWebService {
	return &HttpWebService{}
}

func (h *HttpWebService) Request(url string) (*ports.WebServiceResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return &ports.WebServiceResponse{
		Status: resp.StatusCode,
	}, nil
}
