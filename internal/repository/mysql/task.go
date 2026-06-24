package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error
func (r *TaskRepository) List(ctx context.Context) ([]*model.Task, error) // TODO: add filters
func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error
func (r *TaskRepository) GetHistory(ctx context.Context, taskID uint64) error // TODO: add history model
