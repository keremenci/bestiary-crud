package main

import (
	"github.com/keremenci/bestiary-crud/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", api.HealthCheck)
	router.GET("/beasts", api.ListItems)
	router.POST("/beasts", api.PutItem)
	router.PUT("/beasts", api.UpdateItem)

	router.Run(":8080")
}
