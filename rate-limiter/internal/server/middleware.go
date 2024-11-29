package server

import (
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(rateLimiterUseCase *usecases.RateLimiterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")
		key := ip
		limit := rateLimiterUseCase.LimitIP
		if token != "" {
			key = token
			limit = rateLimiterUseCase.LimitToken
		}
		allowRequest := rateLimiterUseCase.Execute(c.Request.Context(), usecases.RateLimiterUseCaseInput{
			Key:   key,
			Limit: limit,
		})
		if allowRequest {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
		})
	}
}
