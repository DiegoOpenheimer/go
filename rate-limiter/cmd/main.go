package main

import (
	"github.com/DiegoOpenheimer/go/rate-limiter/config"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/server"
)

func main() {
	config.LoadConfig(".")
	server.StartServer()
}
