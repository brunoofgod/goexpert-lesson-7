package middlewares

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brunoofgod/goexpert-lesson-7/internal/ratelimiter"
)

type RateLimiter struct {
	storage       ratelimiter.RateLimiterStorage
	limit         int
	blockDuration time.Duration
}

func NewRateLimiterMiddleware(storage ratelimiter.RateLimiterStorage, limit int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		storage:       storage,
		limit:         limit,
		blockDuration: blockDuration,
	}
}
func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := getRateLimitKey(req)
		headerToken := req.Header.Get("API_KEY")
		tokenFor20Requests := os.Getenv("CUSTOM_TOKEN_REQUESTS")

		if headerToken != "" && headerToken == tokenFor20Requests {
			tokenValueFor20Requests, err := strconv.Atoi(os.Getenv("CUSTOM_TOKEN_REQUESTS_VALUE"))
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			headerToken = tokenFor20Requests
			r.limit = tokenValueFor20Requests
			r.blockDuration = 5 * time.Minute
			key = "rate_limit:token:" + headerToken
		}

		blocked, err := r.storage.IsBlocked(key)
		if err == nil && blocked {
			http.Error(w, "You have reached the maximum number of requests", http.StatusTooManyRequests)
			return
		}

		count, err := r.storage.IncrementRequestCount(key)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			r.storage.SetExpiration(key, r.blockDuration)
		}

		if count > r.limit {
			r.storage.BlockKey(key, r.blockDuration)
			http.Error(w, "You have reached the maximum number of requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getRateLimitKey(req *http.Request) string {
	apiKey := req.Header.Get("API_KEY")
	if strings.TrimSpace(apiKey) != "" {
		return "token:" + apiKey
	}
	return "ip:" + req.RemoteAddr
}
