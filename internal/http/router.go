package http

import (
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/internal/http/handler"
	"github.com/ladislaojs/team-task-manager-test/internal/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	jwtSecret string,
	redisClient *redis.Client,
	maxRequestsPerMinute int,
	userHandler *handler.UserHandler,
	teamHandler *handler.TeamHandler,
	taskHandler *handler.TaskHandler,
) http.Handler {
	mux := http.NewServeMux()
	auth := middleware.JWT(jwtSecret)

	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("POST /api/v1/register", userHandler.Register)
	mux.HandleFunc("POST /api/v1/login", userHandler.Login)

	mux.Handle("POST /api/v1/teams", auth(http.HandlerFunc(teamHandler.Create)))
	mux.Handle("GET /api/v1/teams", auth(http.HandlerFunc(teamHandler.ListForUser)))
	mux.Handle("POST /api/v1/teams/{id}/invite", auth(http.HandlerFunc(teamHandler.Invite)))
	mux.Handle("GET /api/v1/teams/extended", auth(http.HandlerFunc(teamHandler.ListExtended)))
	mux.Handle("GET /api/v1/teams/top-creators", auth(http.HandlerFunc(teamHandler.TopTaskCreators)))

	mux.Handle("POST /api/v1/tasks", auth(http.HandlerFunc(taskHandler.Create)))
	mux.Handle("GET /api/v1/tasks", auth(http.HandlerFunc(taskHandler.List)))
	mux.Handle("PUT /api/v1/tasks/{id}", auth(http.HandlerFunc(taskHandler.Update)))
	mux.Handle("GET /api/v1/tasks/{id}/history", auth(http.HandlerFunc(taskHandler.GetHistory)))
	mux.Handle("GET /api/v1/tasks/other-team-assignee", auth(http.HandlerFunc(taskHandler.OtherTeamAssigneeTasks)))

	rateLimit := middleware.RateLimit(redisClient, maxRequestsPerMinute)
	return middleware.Metrics(rateLimit(mux))
}
