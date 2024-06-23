package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DiegoOpenheimer/go/multithreading/pgk/services"
	"github.com/DiegoOpenheimer/go/multithreading/pgk/utils"
	"os"
	"time"
)

type DataResponse interface {
	JSON() string
}

type Response struct {
	Api  string
	Data DataResponse
}

func main() {
	zipCode := os.Getenv("ZIP_CODE")
	if zipCode == "" {
		zipCode = "01153000"
	}
	ctx := context.Background()
	ch := make(chan Response)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	go func() {
		var zipCodeApi services.ZipCodeApi
		result, err := zipCodeApi.GetZipCodeWithContext(ctxWithTimeout, zipCode)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		utils.HandleError(err)
		ch <- Response{Api: "Via Cep", Data: result}
	}()
	go func() {
		var brazilApi services.BrazilApi
		result, err := brazilApi.GetZipCodeWithContext(ctxWithTimeout, zipCode)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		utils.HandleError(err)
		ch <- Response{Api: "Brasil", Data: result}
	}()

	select {
	case result := <-ch:
		fmt.Printf("Resultado da Api %s: %v\n", result.Api, result.Data.JSON())
	case <-ctxWithTimeout.Done():
		fmt.Println("Timeout: ")
		fmt.Println(ctxWithTimeout.Err())
	}
}
