package middleware

import (
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/pkg/cache"
)

func RateLimit(redisClient *cache.RedisClient, maxRequests int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
}
