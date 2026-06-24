package http

import (
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/internal/http/handler"
)

func NewRouter(
	userHandler *handler.UserHandler,
	teamHandler *handler.TeamHandler,
	taskHandler *handler.TaskHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/v1/users/login", userHandler.Login)

	return mux
}
