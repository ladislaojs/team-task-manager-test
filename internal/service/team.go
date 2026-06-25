package service

import (
	"context"
	"errors"
	"log"

	"github.com/ladislaojs/team-task-manager-test/internal/http/email"
	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository"
)

var (
	ErrTeamNotFound  = errors.New("team not found")
	ErrNotTeamMember = errors.New("not a team member")
	ErrForbidden     = errors.New("forbidden")
	ErrAlreadyMember = errors.New("user is already a member")
)

type TeamService struct {
	teams  repository.TeamRepository
	users  repository.UserRepository
	mailer email.Mailer
}

func NewTeamService(teams repository.TeamRepository, users repository.UserRepository, mailer email.Mailer) *TeamService {
	return &TeamService{teams: teams, users: users, mailer: mailer}
}

func (s *TeamService) Create(ctx context.Context, creatorID uint64, name string) (*model.Team, error) {
	team := &model.Team{
		Name:      name,
		CreatedBy: creatorID,
	}

	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}

	if err := s.teams.AddMember(ctx, &model.TeamMember{
		TeamID: team.ID,
		UserID: creatorID,
		Role:   model.TeamRoleOwner,
	}); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) ListForUser(ctx context.Context, userID uint64) ([]*model.Team, error) {
	return s.teams.FindByMemberID(ctx, userID)
}

func (s *TeamService) ListExtended(ctx context.Context) ([]*model.TeamExtended, error) {
	return s.teams.ListExtended(ctx)
}

func (s *TeamService) Invite(ctx context.Context, inviterID, teamID, inviteeID uint64) error {
	requesterMember, err := s.teams.GetMember(ctx, teamID, inviterID)
	if err != nil {
		return err
	}
	if requesterMember == nil {
		return ErrNotTeamMember
	}
	if !requesterMember.Role.CanInvite() {
		return ErrForbidden
	}

	invitee, err := s.users.FindByID(ctx, inviteeID)
	if err != nil {
		return err
	}
	if invitee == nil {
		return ErrUserNotFound
	}

	existing, err := s.teams.GetMember(ctx, teamID, inviteeID)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrAlreadyMember
	}

	if err := s.teams.AddMember(ctx, &model.TeamMember{
		TeamID: teamID,
		UserID: inviteeID,
		Role:   model.TeamRoleMember,
	}); err != nil {
		return err
	}

	team, err := s.teams.FindByID(ctx, teamID)
	if emailErr := s.mailer.SendInvitation(ctx, invitee.Email, team.Name); emailErr != nil {
		log.Printf("[team-service] invite email failed (user %d): %v", inviteeID, emailErr)
	}

	return nil
}
