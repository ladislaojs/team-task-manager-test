package service

import "github.com/ladislaojs/team-task-manager-test/internal/repository"

type TeamService struct {
	teams repository.TeamRepository
}

func NewTeamService(teams repository.TeamRepository) *TeamService {
	return &TeamService{teams: teams}
}
