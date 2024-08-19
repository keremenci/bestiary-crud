package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database connection pool
var dbPool *pgxpool.Pool

// InitializeDB initializes the database connection pool
func InitializeDB(connString string) {
	var err error
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to the database successfully.")
}

// HealthCheck endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListItems retrieves all items from the database
func ListItems(c *gin.Context) {
	rows, err := dbPool.Query(context.Background(), "SELECT beast_name, type, cr, attributes, description FROM beasts")
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var beasts []Beast
	for rows.Next() {
		var beast Beast
		err = rows.Scan(&beast.BeastName, &beast.Type, &beast.CR, &beast.Attributes, &beast.Description)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		beasts = append(beasts, beast)
	}

	c.JSON(http.StatusOK, beasts)
}

// GetItem retrieves a single item by key from the database
func GetItem(c *gin.Context) {
	key := c.Param("key")
	var beast Beast

	err := dbPool.QueryRow(context.Background(), "SELECT beast_name, type, cr, attributes, description FROM beasts WHERE beast_name=$1", key).
		Scan(&beast.BeastName, &beast.Type, &beast.CR, &beast.Attributes, &beast.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Beast not found"})
		} else {
			log.Printf("Error querying database: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, beast)
}

// PutItem creates a new item in the database
func PutItem(c *gin.Context) {
	var beast Beast
	if err := c.ShouldBindJSON(&beast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := dbPool.Exec(context.Background(), "INSERT INTO beasts (beast_name, type, cr, attributes, description) VALUES ($1, $2, $3, $4, $5)",
		beast.BeastName, beast.Type, beast.CR, beast.Attributes, beast.Description)
	if err != nil {
		log.Printf("Error inserting into database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"BeastName": beast.BeastName})
}

// UpdateItem updates an existing item in the database
func UpdateItem(c *gin.Context) {
	key := c.Param("key")
	var beast Beast
	if err := c.ShouldBindJSON(&beast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := dbPool.Exec(context.Background(), "UPDATE beasts SET type=$1, cr=$2, attributes=$3, description=$4 WHERE beast_name=$5",
		beast.Type, beast.CR, beast.Attributes, beast.Description, key)
	if err != nil {
		log.Printf("Error updating database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beast updated successfully"})
}

// DeleteItem deletes an item from the database
func DeleteItem(c *gin.Context) {
	key := c.Param("key")

	_, err := dbPool.Exec(context.Background(), "DELETE FROM beasts WHERE beast_name=$1", key)
	if err != nil {
		log.Printf("Error deleting from database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beast deleted successfully"})
}
