package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(redisClient *redis.Client, maxRequests int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maxRequests <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			key := fmt.Sprintf("ratelimit:ip:%s", r.RemoteAddr)
			ctx := context.Background()

			count, err := redisClient.Incr(ctx, key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if count == 1 {
				redisClient.Expire(ctx, key, time.Minute)
			}

			if count > int64(maxRequests) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
