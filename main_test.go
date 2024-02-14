package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(rateLimiterMiddleware(RateLimiterConfig{
		Limit:     2,
		BlockTime: 1,
		UseIP:     true,
		UseToken:  false,
	}))

	r.GET("/api/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
	})

	req, err := http.NewRequest("GET", "/api/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Real-IP", "192.168.0.1")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperava um código 200, mas obteve %d", w.Code)
	}

}
