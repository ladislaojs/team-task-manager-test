package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/internal/http/dto"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
	"github.com/ladislaojs/team-task-manager-test/pkg/response"
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
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	if req.Password == "" {
		response.Error(w, http.StatusBadRequest, "password is required")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	user, err := h.users.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			response.Error(w, http.StatusConflict, "email is already taken")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, dto.RegisterResponse{
		ID:    uint64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	if req.Password == "" {
		response.Error(w, http.StatusBadRequest, "password is required")
		return
	}

	tokens, err := h.users.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrIncorrectPassword) {
			response.Error(w, http.StatusUnauthorized, "incorrect email and/or password")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
