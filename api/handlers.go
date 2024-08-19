package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListItems(c *gin.Context) {
	return
}

func GetItem(c *gin.Context) {
	return
}

func PutItem(c *gin.Context) {
	return
}

func DeleteItem(c *gin.Context) {
	return
}

func UpdateItem(c *gin.Context) {
	return
}
