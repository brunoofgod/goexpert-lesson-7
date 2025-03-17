package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/brunoofgod/goexpert-lesson-7/internal/ratelimiter"
	"github.com/brunoofgod/goexpert-lesson-7/internal/ratelimiter/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	limitPerSecond, err := strconv.Atoi(os.Getenv("LIMIT_PER_SECOND"))
	if err != nil {
		limitPerSecond = 10
	}

	configuredDuration, err := strconv.Atoi(os.Getenv("BLOCK_DURATION_PER_SECOND"))
	blockDurationSeconds := 5 * time.Minute

	if err == nil {
		blockDurationSeconds = time.Duration(configuredDuration) * time.Minute
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	router := chi.NewRouter()
	limiter := middlewares.NewRateLimiterMiddleware(ratelimiter.NewRedisRateLimiterStorage(rdb), limitPerSecond, blockDurationSeconds)

	router.Use(limiter.Middleware)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
