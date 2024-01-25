package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/api/resource", handleRequest)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}
