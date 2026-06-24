package repository

import (
	"context"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	List(ctx context.Context) ([]*model.Task, error) // TODO: add filters
	Update(ctx context.Context, task *model.Task) error
	GetHistory(ctx context.Context, taskID uint64) error // TODO: add history model
}
