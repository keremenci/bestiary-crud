package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keremenci/bestiary-crud/api"
	"github.com/keremenci/bestiary-crud/config"
)

func main() {

	// Load config
	cfg := config.GetAppConfig("config/config.yml")
	// Init db connection
	api.InitializeDB(cfg.DatabaseUrl)

	router := gin.Default()

	// Configure CORS to allow all origins TODO: tie this to an env variable
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
	}))

	router.GET("/", api.HealthCheck)
	router.GET("/beasts", api.ListItems)
	router.GET("/beasts/:key", api.GetItem)
	router.POST("/beasts", api.PutItem)
	router.PUT("/beasts/:key", api.UpdateItem)
	router.DELETE("/beasts/:key", api.DeleteItem)

	// Use port from config, default to 8080 if not set
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
