package e2e

import (
	"github.com/DiegoOpenheimer/go/rate-limiter/config"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/infra/storage"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/server"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/usecases"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type RateLimiterTestSuite struct {
	suite.Suite
	server *httptest.Server
	cfg    *config.Config
	redis  *redis.Client
	mr     *miniredis.Miniredis
}

func (suite *RateLimiterTestSuite) SetupTest() {
	_ = os.Setenv("BLOCKED_TIME", "3s")

	// Load configuration
	cfg := config.LoadConfig("../../")

	mr, err := miniredis.Run()
	if err != nil {
		suite.T().Fatalf("não foi possível inicializar o miniredis: %v", err)
	}

	// Configura o cliente Redis para usar o servidor miniredis
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Initialize rate limiter use case
	rateLimiterUseCase := usecases.NewRateLimiterUseCase(storage.NewRedisStorage(client))

	// Create a new Gin router
	router := gin.Default()

	// Apply the rate limit middleware
	router.Use(server.RateLimitMiddleware(rateLimiterUseCase))

	// Define a simple handler for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	// Create a test server
	suite.server = httptest.NewServer(router)
	suite.cfg = cfg
	suite.redis = client
	suite.mr = mr
}

func (suite *RateLimiterTestSuite) TearDownTest() {
	suite.server.Close()
	_ = suite.redis.Close()
}

func (suite *RateLimiterTestSuite) TestRateLimiterForIP() {
	cfg := suite.cfg
	client := &http.Client{}
	var runnable = func() {
		// Test exceeding rate limit for IP
		for i := range cfg.RateLimitIP + 1 {
			suite.assertRateLimiter(client, cfg.RateLimitIP, i, true)
		}
	}
	runnable()
	suite.mr.FastForward(cfg.BlockedTime)
	runnable()
}

func (suite *RateLimiterTestSuite) TestRateLimiterForToken() {
	cfg := suite.cfg
	client := &http.Client{}

	var runnable = func() {
		// Test exceeding rate limit for token
		for i := range cfg.RateLimitToken + 1 {
			suite.assertRateLimiter(client, cfg.RateLimitToken, i, false)
		}
	}
	runnable()
	suite.mr.FastForward(cfg.BlockedTime)
	runnable()
	// Test exceeding rate limit for token after blocked time
}

func (suite *RateLimiterTestSuite) assertRateLimiter(client *http.Client, limit int, i int, isTestIp bool) {
	req, err := http.NewRequest("GET", suite.server.URL, nil)
	assert.NoError(suite.T(), err)
	if !isTestIp {
		req.Header.Set("API_KEY", "test-token")
	}

	resp, err := client.Do(req)
	assert.NoError(suite.T(), err)

	if i < limit {
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(suite.T(), "Hello World!", body["message"])
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, resp.StatusCode)
		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(suite.T(), "you have reached the maximum number of requests or actions allowed within a certain time frame", body["message"])
	}
	_ = resp.Body.Close()
}

func TestRateLimiter(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}
