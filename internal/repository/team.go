package repository

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	FindByMemberId(ctx context.Context, userID *uint64)
	AddMember(ctx context.Context, member *model.TeamMember) error
}
