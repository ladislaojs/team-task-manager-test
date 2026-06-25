package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository"
)

const taskListTTL = 5 * time.Minute

type TaskRepository struct {
	repo  repository.TaskRepository
	redis *redis.Client
}

func NewTaskRepository(repo repository.TaskRepository, redis *redis.Client) *TaskRepository {
	return &TaskRepository{repo: repo, redis: redis}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	if err := r.repo.Create(ctx, task); err != nil {
		return err
	}
	r.invalidateTeam(ctx, task.TeamID)
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uint64) (*model.Task, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *TaskRepository) List(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	key := listCacheKey(filter)

	cached, err := r.redis.Get(ctx, key).Bytes()
	if err == nil {
		var tasks []*model.Task
		if json.Unmarshal(cached, &tasks) == nil {
			return tasks, nil
		}
	}

	tasks, err := r.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(tasks); err == nil {
		r.redis.Set(ctx, key, data, taskListTTL)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	if err := r.repo.Update(ctx, task); err != nil {
		return err
	}
	r.invalidateTeam(ctx, task.TeamID)
	return nil
}

func (r *TaskRepository) AddHistory(ctx context.Context, entry *model.TaskHistory) error {
	return r.repo.AddHistory(ctx, entry)
}

func (r *TaskRepository) GetHistory(ctx context.Context, taskID uint64) ([]*model.TaskHistory, error) {
	return r.repo.GetHistory(ctx, taskID)
}

func (r *TaskRepository) OtherTeamAssigneeTasks(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error) {
	return r.repo.OtherTeamAssigneeTasks(ctx)
}

func (r *TaskRepository) invalidateTeam(ctx context.Context, teamID uint64) {
	pattern := fmt.Sprintf("tasks:team:%d:*", teamID)
	iter := r.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		r.redis.Del(ctx, iter.Val())
	}
}

func listCacheKey(f model.TaskFilter) string {
	status := ""
	if f.Status != nil {
		status = string(*f.Status)
	}
	assignee := uint64(0)
	if f.AssigneeID != nil {
		assignee = *f.AssigneeID
	}
	return fmt.Sprintf("tasks:team:%d:status:%s:assignee:%d:page:%d:size:%d",
		f.TeamID, status, assignee, f.Page, f.PageSize)
}
