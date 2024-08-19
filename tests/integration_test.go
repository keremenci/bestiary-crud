package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keremenci/bestiary-crud/api"
	"github.com/keremenci/bestiary-crud/config"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	// Initialize the database connection
	api.InitializeDB(config.GetAppConfig("../config/config.yml").DatabaseUrl)

	router := gin.Default()
	router.GET("/", api.HealthCheck)
	router.GET("/beasts", api.ListItems)
	router.GET("/beasts/:key", api.GetItem)
	router.POST("/beasts", api.PutItem)
	router.PUT("/beasts/:key", api.UpdateItem)
	router.DELETE("/beasts/:key", api.DeleteItem)

	return router
}

func TestIntegration(t *testing.T) {
	router := setupRouter()

	t.Run("HealthCheck", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
	})

	t.Run("CreateBeast", func(t *testing.T) {
		beast := api.Beast{
			BeastName:   "IntegrationTestBeast",
			Type:        "TestType",
			CR:          "1",
			Attributes:  map[string]string{"STR": "10"},
			Description: "Integration test description",
		}
		jsonValue, _ := json.Marshal(beast)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/beasts", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		assert.Equal(t, "IntegrationTestBeast", response["BeastName"])
	})

	t.Run("GetBeast", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/beasts/IntegrationTestBeast", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.Beast
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		assert.Equal(t, "IntegrationTestBeast", response.BeastName)
	})

	t.Run("UpdateBeast", func(t *testing.T) {
		beast := api.Beast{
			Type:        "UpdatedType",
			CR:          "2",
			Attributes:  map[string]string{"STR": "12"},
			Description: "Updated description",
		}
		jsonValue, _ := json.Marshal(beast)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/beasts/IntegrationTestBeast", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		assert.Equal(t, "Beast updated successfully", response["message"])
	})

	t.Run("DeleteBeast", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/beasts/IntegrationTestBeast", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		assert.Equal(t, "Beast deleted successfully", response["message"])
	})
}
