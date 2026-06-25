package repository

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id uint64) (*model.Task, error)
	List(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	AddHistory(ctx context.Context, entry *model.TaskHistory) error
	GetHistory(ctx context.Context, taskID uint64) ([]*model.TaskHistory, error)
	OtherTeamAssigneeTasks(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error)
}
