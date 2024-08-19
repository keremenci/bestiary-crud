package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	router := gin.Default()
	router.GET("/", HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestShouldListItems(t *testing.T) {
	// Setup mock database

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	SetDBPool(mock)

	// Setup rows
	rows := mock.NewRows([]string{"beast_name", "type", "cr", "attributes", "description"}).
		AddRow("TestBeast", "TestType", "1", map[string]string{"STR": "10"}, "Test description")
	mock.ExpectQuery("SELECT beast_name, type, cr, attributes, description FROM beasts").WillReturnRows(rows)

	// Setup router

	router := gin.Default()
	router.GET("/beasts", ListItems)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/beasts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetItem(t *testing.T) {
	// Setup mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	SetDBPool(mock)

	rows := mock.NewRows([]string{"beast_name", "type", "cr", "attributes", "description"}).
		AddRow("TestBeast", "TestType", "1", map[string]string{"STR": "10"}, "Test description")

	// Use a regular expression to match the query with parameter placeholder
	queryRegex := regexp.QuoteMeta("SELECT beast_name, type, cr, attributes, description FROM beasts WHERE beast_name=$1")
	mock.ExpectQuery(queryRegex).WithArgs("TestBeast").WillReturnRows(rows)

	// Setup router
	router := gin.Default()
	router.GET("/beasts/:key", GetItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/beasts/TestBeast", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	log.Printf("Response: %v", response)

	// Add more assertions based on the expected JSON response
	assert.Equal(t, "TestBeast", response["BeastName"])
	assert.Equal(t, "TestType", response["Type"])
	assert.Equal(t, "1", response["CR"])
	assert.Equal(t, "Test description", response["Description"])

	// Check attributes
	attributesResponse, ok := response["Attributes"].(map[string]interface{})
	if assert.True(t, ok, "Attributes should be a map") {
		assert.Equal(t, "10", attributesResponse["STR"])
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPutItem(t *testing.T) {
	router := gin.Default()
	router.POST("/beasts", PutItem)

	beast := Beast{
		BeastName:   "TestBeast",
		Type:        "TestType",
		CR:          "1",
		Attributes:  map[string]string{"STR": "10"},
		Description: "Test description",
	}
	jsonValue, _ := json.Marshal(beast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/beasts", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	// Add more assertions based on the expected JSON response
}

func TestUpdateItem(t *testing.T) {
	router := gin.Default()
	router.PUT("/beasts/:key", UpdateItem)

	beast := Beast{
		BeastName:   "TestBeast",
		Type:        "UpdatedType",
		CR:          "2",
		Attributes:  map[string]string{"STR": "12"},
		Description: "Updated description",
	}
	jsonValue, _ := json.Marshal(beast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/beasts/TestBeast", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions based on the expected JSON response
}

func TestDeleteItem(t *testing.T) {
	router := gin.Default()
	router.DELETE("/beasts/:key", DeleteItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/beasts/TestBeast", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions based on the expected JSON response
}
