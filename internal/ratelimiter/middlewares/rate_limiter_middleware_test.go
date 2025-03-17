package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brunoofgod/goexpert-lesson-7/internal/ratelimiter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(t *testing.T) (*redis.Client, func()) {
	os.Setenv("CUSTOM_TOKEN_REQUESTS", "5ae67c14-c5f2-4ced-b96b-8c1630e1c5e6")
	os.Setenv("CUSTOM_TOKEN_REQUESTS_VALUE", "20")
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ip, err := container.Host(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	port, err := container.MappedPort(context.Background(), "6379")
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: ip + ":" + port.Port(),
	})

	return rdb, func() { container.Terminate(context.Background()) }
}
func TestRateLimiterWithdoutToken(t *testing.T) {
	rdb, cleanup := setupRedisContainer(t)
	defer cleanup()

	rateLimiter := NewRateLimiterMiddleware(ratelimiter.NewRedisRateLimiterStorage(rdb), 10, 5*time.Second)
	ts := httptest.NewServer(rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer ts.Close()

	client := &http.Client{}

	for i := 0; i < 11; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		if i < 10 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "A requisição %d deveria ser permitida", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "A requisição %d deveria ser bloqueada", i+1)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestRateLimiterWithTokenFor20Requests(t *testing.T) {
	rdb, cleanup := setupRedisContainer(t)
	defer cleanup()

	rateLimiter := NewRateLimiterMiddleware(ratelimiter.NewRedisRateLimiterStorage(rdb), 10, 5*time.Second)
	ts := httptest.NewServer(rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer ts.Close()

	client := &http.Client{}
	token := os.Getenv("CUSTOM_TOKEN_REQUESTS")

	for i := 0; i < 21; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/", nil)
		require.NoError(t, err)

		req.Header.Set("API_KEY", token)

		resp, err := client.Do(req)
		require.NoError(t, err)

		if i < 20 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "A requisição %d deveria ser permitida", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "A requisição %d deveria ser bloqueada", i+1)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
