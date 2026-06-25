package mocks

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TeamRepository struct {
	CreateFn         func(ctx context.Context, team *model.Team) error
	FindByMemberIDFn func(ctx context.Context, userID uint64) ([]*model.Team, error)
	FindByIDFn       func(ctx context.Context, id uint64) (*model.Team, error)
	GetMemberFn      func(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error)
	AddMemberFn      func(ctx context.Context, member *model.TeamMember) error
	ListExtendedFn   func(ctx context.Context) ([]*model.TeamExtended, error)
}

func (m *TeamRepository) Create(ctx context.Context, team *model.Team) error {
	return m.CreateFn(ctx, team)
}

func (m *TeamRepository) FindByMemberID(ctx context.Context, userID uint64) ([]*model.Team, error) {
	return m.FindByMemberIDFn(ctx, userID)
}

func (m *TeamRepository) GetMember(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error) {
	return m.GetMemberFn(ctx, teamID, userID)
}

func (m *TeamRepository) AddMember(ctx context.Context, member *model.TeamMember) error {
	return m.AddMemberFn(ctx, member)
}

func (m *TeamRepository) ListExtended(ctx context.Context) ([]*model.TeamExtended, error) {
	return m.ListExtendedFn(ctx)
}

func (m *TeamRepository) FindByID(ctx context.Context, id uint64) (*model.Team, error) {
	return m.FindByIDFn(ctx, id)
}
