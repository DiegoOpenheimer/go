package main

import (
	"challenge-client-server/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/cotacao", port), nil)
	utils.HandlerError(err)
	resp, err := http.DefaultClient.Do(req)
	utils.HandlerError(err)
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error on request")
		body, err := io.ReadAll(resp.Body)
		utils.HandlerError(err)
		fmt.Printf("Status %d, Error: %s\n", resp.StatusCode, string(body))
		return
	}

	body, err := io.ReadAll(resp.Body)
	utils.HandlerError(err)

	err = os.WriteFile("cotacao.txt", []byte("Dólar: "+string(body)), 0644)
	utils.HandlerError(err)
	fmt.Println("Quotation saved in cotacao.txt")
}
