package mocks

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TaskRepository struct {
	CreateFn                 func(ctx context.Context, task *model.Task) error
	FindByIDFn               func(ctx context.Context, id uint64) (*model.Task, error)
	ListFn                   func(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error)
	UpdateFn                 func(ctx context.Context, task *model.Task) error
	AddHistoryFn             func(ctx context.Context, entry *model.TaskHistory) error
	GetHistoryFn             func(ctx context.Context, taskID uint64) ([]*model.TaskHistory, error)
	OtherTeamAssigneeTasksFn func(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error)
}

func (m *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	return m.CreateFn(ctx, task)
}

func (m *TaskRepository) FindByID(ctx context.Context, id uint64) (*model.Task, error) {
	return m.FindByIDFn(ctx, id)
}

func (m *TaskRepository) List(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	return m.ListFn(ctx, filter)
}

func (m *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	return m.UpdateFn(ctx, task)
}

func (m *TaskRepository) AddHistory(ctx context.Context, entry *model.TaskHistory) error {
	return m.AddHistoryFn(ctx, entry)
}

func (m *TaskRepository) GetHistory(ctx context.Context, taskID uint64) ([]*model.TaskHistory, error) {
	return m.GetHistoryFn(ctx, taskID)
}

func (m *TaskRepository) OtherTeamAssigneeTasks(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error) {
	return m.OtherTeamAssigneeTasksFn(ctx)
}
