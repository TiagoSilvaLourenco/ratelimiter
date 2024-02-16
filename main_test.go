package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisPersistence(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r := &RedisPersistence{client: client}

	t.Run("Test GetLimit", func(t *testing.T) {
		limit, err := r.GetLimit("test", 10)
		assert.NoError(t, err)
		assert.Equal(t, 10, limit)
	})

	t.Run("Test Incr", func(t *testing.T) {
		count, err := r.Incr("test")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Test Expire", func(t *testing.T) {
		err := r.Expire(context.Background(), "test", 1*time.Second)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second)
		val, err := r.client.Get(context.Background(), "test").Result()
		assert.Equal(t, redis.Nil, err)
		assert.Equal(t, "", val)
	})
}
