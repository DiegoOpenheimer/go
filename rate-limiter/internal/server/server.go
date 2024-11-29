package server

import (
	"github.com/DiegoOpenheimer/go/rate-limiter/config"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/infra/storage"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
)

func StartServer() {
	cfg := config.GetConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	rateLimiterUseCase := usecases.NewRateLimiterUseCase(storage.NewRedisStorage(rdb))

	r := gin.Default()
	r.Use(RateLimitMiddleware(rateLimiterUseCase))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	log.Fatal(r.Run(":8080"))
}
