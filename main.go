package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func init() {

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use default DB
	})
}

func main() {
	r := gin.Default()
	r.Use(rateLimiterMiddleware)
	r.GET("/api/resource", handleRequest)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}

func rateLimiterMiddleware(c *gin.Context) {

	// get ip address of client
	ip := c.ClientIP()

	// get token of access from header
	token := c.GetHeader("API_KEY")

	println("IP: ", ip)
	println("Token: ", token)

}
