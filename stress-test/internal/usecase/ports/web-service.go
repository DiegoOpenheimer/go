package ports

type WebServiceResponse struct {
	Status int
}

type WebService interface {
	Request(url string) (*WebServiceResponse, error)
}
