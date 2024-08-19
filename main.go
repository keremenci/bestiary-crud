package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keremenci/bestiary-crud/api"
	"github.com/keremenci/bestiary-crud/config"
)

func main() {

	// Load config
	config.GetAppConfig()
	// Init db connection
	api.InitializeDB(config.GetAppConfig().DatabaseUrl)

	router := gin.Default()

	router.GET("/", api.HealthCheck)
	router.GET("/beasts", api.ListItems)
	router.GET("/beasts/:key", api.GetItem)
	router.POST("/beasts", api.PutItem)
	router.PUT("/beasts/:key", api.UpdateItem)
	router.DELETE("/beasts/:key", api.DeleteItem)

	router.Run(":8080")
}
