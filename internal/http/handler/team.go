package handler

import (
	"net/http"

	"github.com/ladislaojs/team-task-manager-test/internal/service"
)

type TeamHandler struct {
	teams *service.TeamService
}

func NewTeamHandler(teams *service.TeamService) *TeamHandler {
	return &TeamHandler{teams: teams}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
}

func (h *TeamHandler) ListForUser(w http.ResponseWriter, r *http.Request) {

}

func (h *TeamHandler) Invite(w http.ResponseWriter, r *http.Request) {

}

func (h *TeamHandler) ListExtended(w http.ResponseWriter, r *http.Request) {

}

func (h *TeamHandler) TopTaskCreators(w http.ResponseWriter, r *http.Request) {

}
