package usecases

import (
	"context"
	"github.com/DiegoOpenheimer/go/rate-limiter/config"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/usecases/ports"
	"strconv"
	"time"
)

type RateLimiterUseCase struct {
	Storage     ports.Storage
	LimitIP     int
	LimitToken  int
	BlockedTime time.Duration
}

func NewRateLimiterUseCase(storage ports.Storage) *RateLimiterUseCase {
	cfg := config.GetConfig()
	return &RateLimiterUseCase{
		Storage:     storage,
		LimitIP:     cfg.RateLimitIP,
		LimitToken:  cfg.RateLimitToken,
		BlockedTime: cfg.BlockedTime,
	}
}

type RateLimiterUseCaseInput struct {
	Key   string
	Limit int
}

func (rl *RateLimiterUseCase) Execute(ctx context.Context, input RateLimiterUseCaseInput) bool {
	valueStorage, err := rl.Storage.GetWithContext(ctx, input.Key)
	if err != nil {
		return false
	}
	value := rl.getValueInStorage(valueStorage)
	if value >= input.Limit {
		_ = rl.Storage.IncrementWithContext(ctx, input.Key, &ports.StorageOptions{
			Expire: rl.BlockedTime,
		})
		return false
	}
	_ = rl.Storage.IncrementWithContext(ctx, input.Key, &ports.StorageOptions{
		Expire: time.Second,
	})
	return true
}

func (rl *RateLimiterUseCase) getValueInStorage(value string) int {
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}
