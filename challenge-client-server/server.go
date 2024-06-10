package main

import (
	"challenge-client-server/database"
	"challenge-client-server/entities"
	"challenge-client-server/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valyala/fastjson"
)

func init() {
	database.InitDB()
}

func main() {
	http.Handle("GET /cotacao", middlewareError(http.HandlerFunc(handlerQuotation)))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func middlewareError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(error); ok && errors.Is(r.(error), context.DeadlineExceeded) {
					fmt.Println("Timeout exceeded", r)
					http.Error(w, "Timeout exceeded", http.StatusRequestTimeout)
					return
				}
				fmt.Println("Recovered from panic", r)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func handlerQuotation(w http.ResponseWriter, r *http.Request) {
	ctxApiRequest, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()
	request, err := http.NewRequestWithContext(ctxApiRequest, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	utils.HandlerError(err)
	resp, err := http.DefaultClient.Do(request)
	utils.HandlerError(err)
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyInBytes, err := io.ReadAll(resp.Body)
	utils.HandlerError(err)

	var p fastjson.Parser
	value, err := p.ParseBytes(bodyInBytes)
	utils.HandlerError(err)
	var quotation entities.Quotation
	err = json.Unmarshal(value.GetObject("USDBRL").MarshalTo(nil), &quotation)
	utils.HandlerError(err)

	ctxDB, cancel := context.WithTimeout(r.Context(), time.Millisecond*10)
	defer cancel()
	result := database.Db.WithContext(ctxDB).Create(&quotation)
	utils.HandlerError(result.Error)

	err = json.NewEncoder(w).Encode(quotation)
	utils.HandlerError(err)
	fmt.Println("Current quote", quotation.Bid)
}
