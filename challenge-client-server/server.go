package main

import (
	"challenge-client-server/Entities"
	"challenge-client-server/config"
	"challenge-client-server/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/valyala/fastjson"
)

func init() {
	config.InitDB()
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
				if errors.Is(r.(error), context.DeadlineExceeded) {
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
	rawBid := value.GetStringBytes("USDBRL", "bid")
	bid, err := strconv.ParseFloat(string(rawBid), 64)
	utils.HandlerError(err)

	ctxDB, cancel := context.WithTimeout(r.Context(), time.Millisecond*10)
	defer cancel()
	config.Db.WithContext(ctxDB).Create(&Entities.Quotation{Value: bid})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(strconv.FormatFloat(bid, 'f', -1, 64)))
	fmt.Println("Current quote", bid)
}
