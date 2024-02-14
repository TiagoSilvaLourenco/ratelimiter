package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	ipConfig := RateLimiterConfig{
		Limit:     getEnvInt("IP_LIMIT", 10),
		BlockTime: time.Second,
		UseIP:     true,
		UseToken:  false,
	}
	r.Use(rateLimiterMiddleware(ipConfig))

	tokenConfig := RateLimiterConfig{
		Limit:     getEnvInt("TOKEN_LIMIT", 100),
		BlockTime: 2 * time.Second,
		UseIP:     false,
		UseToken:  true,
	}
	r.Use(rateLimiterMiddleware(tokenConfig))

	r.GET("/api/resource", handleRequest)
	r.Run(":8080")
}

func handleRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}

func rateLimiterMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if config.UseIP {
			key = "ip:" + c.ClientIP()
		} else if config.UseToken {
			key = "token:" + c.GetHeader("API_KEY")
		}

		limit, err := getLimit(key, config.Limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		count, _ := client.Incr(ctx, "count:"+key).Result()
		client.Expire(ctx, "count:"+key, config.BlockTime)

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
