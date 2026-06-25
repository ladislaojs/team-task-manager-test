package mocks

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type UserRepository struct {
	CreateFn                 func(ctx context.Context, user *model.User) error
	FindByEmailFn            func(ctx context.Context, email string) (*model.User, error)
	FindByIDFn               func(ctx context.Context, id uint64) (*model.User, error)
	TopTaskCreatorsPerTeamFn func(ctx context.Context) ([]*model.TaskCreator, error)
}

func (m *UserRepository) Create(ctx context.Context, user *model.User) error {
	return m.CreateFn(ctx, user)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.FindByEmailFn(ctx, email)
}

func (m *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	return m.FindByIDFn(ctx, id)
}

func (m *UserRepository) TopTaskCreatorsPerTeam(ctx context.Context) ([]*model.TaskCreator, error) {
	return m.TopTaskCreatorsPerTeamFn(ctx)
}
