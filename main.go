package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client *redis.Client

type RateLimiterConfig struct {
	Limit     int
	BlockTime time.Duration
	UseIP     bool
	UseToken  bool
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://resttesttest.com"}
	config.AllowHeaders = []string{"Content-Type", "API_KEY"}
	config.AllowMethods = []string{"GET"}
	config.AllowCredentials = true

	r.Use(cors.New(config))
	r.Use(rateLimiterMiddleware())

	r.GET("/api/resource", handleRequest)
	r.Run(":8080")
}

func handleRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}

func rateLimiterMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		var key string
		var limit int
		blockTime := time.Second

		if c.GetHeader("API_KEY") != "" {
			key = "token:" + c.GetHeader("API_KEY")
			limit = getEnvInt("TOKEN_LIMIT", 100)

		} else {
			key = "ip:" + c.ClientIP()
			limit = getEnvInt("IP_LIMIT", 10)
			blockTime = blockTime * 2
		}

		limit, err := getLimit(key, limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		count, _ := client.Incr(ctx, "count:"+key).Result()
		client.Expire(ctx, "count:"+key, blockTime)

		log.Println("count: ", count)
		log.Println("key: ", key)
		log.Println("limit: ", limit)

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "You have reached the maximum number of requests or actions allowed within a certain time frame"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getLimit(key string, defaultLimit int) (int, error) {
	limitStr, err := client.Get(ctx, key).Result()
	// log.Println("Limit: ", limitStr)
	if err == redis.Nil {
		client.Set(ctx, key, strconv.Itoa(defaultLimit), 0)
		return defaultLimit, nil
	} else if err != nil {
		return 0, err
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, err
		}
		if limit == defaultLimit {

			return limit, nil
		}
	}

	err = client.Set(ctx, key, strconv.Itoa(defaultLimit), 0).Err()
	if err != nil {
		return 0, err
	}
	return defaultLimit, nil
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default value %d\n", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}
