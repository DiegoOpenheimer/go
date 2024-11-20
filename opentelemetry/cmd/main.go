package main

import (
	"context"
	"github.com/DiegoOpenheimer/go/opentelemetry/configs"
	"github.com/DiegoOpenheimer/go/opentelemetry/internal/webservices"
	"github.com/DiegoOpenheimer/go/opentelemetry/pgk"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	shutdown, err := pgk.InitProvider(cfg.ServiceName, cfg.OtelExporterOtlpEndpoint)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	if cfg.ServiceAPort != 0 {
		go func() {
			webservices.StartServiceA()
		}()
	}
	if cfg.ServiceBPort != 0 {
		go func() {
			webservices.StartServiceB()
		}()
	}

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}
	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

}
