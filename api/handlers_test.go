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

// TestHealthCheck tests the GET / endpoint
func TestHealthCheck(t *testing.T) {
	router := gin.Default()
	router.GET("/", HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

// TestListItems tests the GET /beasts endpoint
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

// TestGetItem tests the GET /beasts/:key endpoint
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

// TestPutItem tests the PUT /beasts endpoint
func TestPutItem(t *testing.T) {
	// Setup mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	SetDBPool(mock)

	// Use ExpectExec for INSERT queries
	queryRegex := regexp.QuoteMeta("INSERT INTO beasts (beast_name, type, cr, attributes, description) VALUES ($1, $2, $3, $4, $5)")
	mock.ExpectExec(queryRegex).
		WithArgs("TestBeast", "TestType", "1", map[string]string{"STR": "10"}, "Test description").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Setup Router
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

	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Add more assertions based on the expected JSON response
	assert.Equal(t, "TestBeast", response["BeastName"])

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateItem(t *testing.T) {
	// Setup mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	SetDBPool(mock)

	// Define the expected query and arguments for the UPDATE operation
	queryRegex := regexp.QuoteMeta("UPDATE beasts SET type=$1, cr=$2, attributes=$3, description=$4 WHERE beast_name=$5")
	mock.ExpectExec(queryRegex).
		WithArgs("UpdatedType", "2", map[string]string{"STR": "12"}, "Updated description", "TestBeast").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Setup router
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

	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Add assertions based on the expected JSON response
	assert.Equal(t, "Beast updated successfully", response["message"])

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteItem(t *testing.T) {
	// Setup mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	SetDBPool(mock)

	// Define the expected query for the DELETE operation
	queryRegex := regexp.QuoteMeta("DELETE FROM beasts WHERE beast_name=$1")
	mock.ExpectExec(queryRegex).
		WithArgs("TestBeast").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Setup router
	router := gin.Default()
	router.DELETE("/beasts/:key", DeleteItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/beasts/TestBeast", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Add assertions based on the expected JSON response
	assert.Equal(t, "Beast deleted successfully", response["message"])

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
