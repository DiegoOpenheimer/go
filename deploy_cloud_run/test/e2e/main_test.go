package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/configs"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/infra/web/handlers"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/internal/usecases"
	"github.com/DiegoOpenheimer/go/deploy_cloud_run/pgk/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APITestSuite struct {
	suite.Suite
	server *httptest.Server
	router *chi.Mux
}

func (suite *APITestSuite) SetupTest() {
	_, err := configs.LoadConfig("../../")
	require.NoError(suite.T(), err, "Failed to load config")

	suite.router = chi.NewRouter()
	suite.router.Use(middleware.Logger)

	temperatureHandler := handlers.NewTemperatureHandler(
		usecases.NewTemperatureUseCase(
			services.NewZipCodeApi(),
			services.NewTemperatureApi(),
		),
	)

	suite.router.Get("/{code}", temperatureHandler.GetTemperature)
	suite.server = httptest.NewServer(suite.router)
}

func (suite *APITestSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *APITestSuite) TestServerRunningOnPort() {
	resp, err := http.Get(fmt.Sprintf("%s/37561190", suite.server.URL))
	require.NoError(suite.T(), err, "Failed to make GET request")

	suite.Equal(http.StatusOK, resp.StatusCode, "Expected status code 200")
}

func (suite *APITestSuite) TestGetTemperatureInvalidZipCode() {
	resp, err := http.Get(fmt.Sprintf("%s/invalid_zip_code", suite.server.URL))
	require.NoError(suite.T(), err, "Failed to make GET request")

	suite.Equal(http.StatusUnprocessableEntity, resp.StatusCode, "Expected status code 422")
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(suite.T(), err, "Failed to decode response")
	suite.Equal("invalid zipcode", result["message"], "Expected message to be invalid zipcode")
}

func (suite *APITestSuite) TestGetTemperatureNotFound() {
	resp, err := http.Get(fmt.Sprintf("%s/01001007", suite.server.URL))
	require.NoError(suite.T(), err, "Failed to make GET request")

	suite.Equal(http.StatusNotFound, resp.StatusCode, "Expected status code 404")
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(suite.T(), err, "Failed to decode response")
	suite.Equal("can not find zipcode", result["message"], "Expected message to be can not find zipcode")
}

func (suite *APITestSuite) TestGetTemperatureValidResponse() {
	resp, err := http.Get(fmt.Sprintf("%s/37561190", suite.server.URL))
	require.NoError(suite.T(), err, "Failed to make GET request")

	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(suite.T(), err, "Failed to decode response")

	_, ok := result["temp_C"]
	suite.True(ok, "Expected temp_C in response")
	_, ok = result["temp_F"]
	suite.True(ok, "Expected temp_F in response")
	_, ok = result["temp_K"]
	suite.True(ok, "Expected temp_K in response")

	tempK := result["temp_K"].(float64)
	suite.Greater(tempK, 0.0, "Expected temp_K to be greater than 0")
	tempC := result["temp_C"].(float64)
	suite.Greater(tempC, 0.0, "Expected temp_C to be greater than 0")
	suite.True(tempC == tempK-273, "Expected temp_C to be equal to temp_K - 273")
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
