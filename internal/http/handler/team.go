package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ladislaojs/team-task-manager-test/internal/http/dto"
	"github.com/ladislaojs/team-task-manager-test/internal/http/middleware"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
	"github.com/ladislaojs/team-task-manager-test/pkg/response"
)

type TeamHandler struct {
	teams *service.TeamService
	users *service.UserService
}

func NewTeamHandler(teams *service.TeamService, users *service.UserService) *TeamHandler {
	return &TeamHandler{teams: teams, users: users}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	team, err := h.teams.Create(r.Context(), claims.UserID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, dto.TeamResponse{
		ID:   team.ID,
		Name: team.Name,
	})
}

func (h *TeamHandler) ListForUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	teams, err := h.teams.ListForUser(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	teamDTOs := make([]dto.TeamResponse, len(teams))
	for i, t := range teams {
		teamDTOs[i] = dto.TeamResponse{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	response.JSON(w, http.StatusOK, teamDTOs)
}

func (h *TeamHandler) Invite(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	teamID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid team id")
		return
	}

	var req dto.InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == 0 {
		response.Error(w, http.StatusBadRequest, "user_id is required")
		return
	}

	err = h.teams.Invite(r.Context(), claims.UserID, teamID, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotTeamMember):
			response.Error(w, http.StatusForbidden, "not a team member")
		case errors.Is(err, service.ErrForbidden):
			response.Error(w, http.StatusForbidden, "only owner or admin can invite")
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "user not found")
		case errors.Is(err, service.ErrAlreadyMember):
			response.Error(w, http.StatusConflict, "user is already a member")
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ListExtended(w http.ResponseWriter, r *http.Request) {
	stats, err := h.teams.ListExtended(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	teamExtDTOs := make([]dto.TeamExtendedResponse, len(stats))
	for i, s := range stats {
		teamExtDTOs[i] = dto.TeamExtendedResponse{
			ID:                    s.ID,
			Name:                  s.Name,
			MemberCount:           s.MemberCount,
			LastWeekDoneTaskCount: s.LastWeekDoneTaskCount,
		}
	}
	response.JSON(w, http.StatusOK, teamExtDTOs)
}

func (h *TeamHandler) TopTaskCreators(w http.ResponseWriter, r *http.Request) {
	creators, err := h.users.TopTaskCreatorsPerTeam(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	creatorDTOs := make([]dto.TopTaskCreatorResponse, len(creators))
	for i, c := range creators {
		creatorDTOs[i] = dto.TopTaskCreatorResponse{
			TeamID:    c.TeamID,
			UserID:    c.UserID,
			UserName:  c.UserName,
			TaskCount: c.TaskCount,
			Rank:      c.Rank,
		}
	}
	response.JSON(w, http.StatusOK, creatorDTOs)
}
