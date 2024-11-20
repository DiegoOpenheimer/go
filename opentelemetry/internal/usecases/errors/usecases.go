package errors

import "net/http"

type UseCaseError interface {
	error
	GetCode() int
}

type UnknownError struct{}

func NewUnknownError() error {
	return UnknownError{}
}

func (u UnknownError) Error() string {
	return "Unknown error"
}

func (u UnknownError) GetCode() int {
	return http.StatusInternalServerError
}

type ZipCodeNotFoundError struct{}

func NewZipCodeNotFoundError() error {
	return ZipCodeNotFoundError{}
}

func (z ZipCodeNotFoundError) Error() string {
	return "can not find zipcode"
}

func (z ZipCodeNotFoundError) GetCode() int {
	return http.StatusNotFound
}

type InvalidZipCodeError struct{}

func NewInvalidZipCodeError() error {
	return InvalidZipCodeError{}
}

func (i InvalidZipCodeError) Error() string {
	return "invalid zipcode"
}

func (i InvalidZipCodeError) GetCode() int {
	return http.StatusUnprocessableEntity
}
