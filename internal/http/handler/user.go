package handler

import (
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{
		users: users,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
}
