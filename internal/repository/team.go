package repository

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	FindByMemberID(ctx context.Context, userID uint64) ([]*model.Team, error)
	ListExtended(ctx context.Context) ([]*model.TeamExtended, error)
	FindByID(ctx context.Context, id uint64) (*model.Team, error)
	AddMember(ctx context.Context, member *model.TeamMember) error
	GetMember(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error)
}
