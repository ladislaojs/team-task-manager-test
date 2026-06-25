package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type UpdatedTask struct {
	Title       *string
	Description *string
	Status      *model.TaskStatus
	AssigneeID  *uint64
	DueDate     *time.Time
}

type TaskService struct {
	tasks repository.TaskRepository
	teams repository.TeamRepository
}

func NewTaskService(tasks repository.TaskRepository, teams repository.TeamRepository) *TaskService {
	return &TaskService{tasks: tasks, teams: teams}
}

func (s *TaskService) Create(ctx context.Context, creatorID, teamID uint64, assigneeID *uint64, title, description string, dueDate *time.Time) (*model.Task, error) {
	member, err := s.teams.GetMember(ctx, teamID, creatorID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	task := &model.Task{
		TeamID:      teamID,
		CreatedBy:   creatorID,
		AssigneeID:  assigneeID,
		Title:       title,
		Description: description,
		Status:      model.StatusTodo,
		DueDate:     dueDate,
	}

	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) List(ctx context.Context, requesterID uint64, filter model.TaskFilter) ([]*model.Task, error) {
	member, err := s.teams.GetMember(ctx, filter.TeamID, requesterID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	if filter.PageSize == 0 {
		filter.PageSize = 20
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	return s.tasks.List(ctx, filter)
}

func (s *TaskService) Update(ctx context.Context, updaterID, taskID uint64, updatedTask UpdatedTask) (*model.Task, error) {
	task, err := s.tasks.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	member, err := s.teams.GetMember(ctx, task.TeamID, updaterID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	canUpdate := member.Role.CanUpdate()
	isInvolved := task.CreatedBy == updaterID || (task.AssigneeID != nil && *task.AssigneeID == updaterID)
	if !canUpdate && !isInvolved {
		return nil, ErrForbidden
	}

	var entries []*model.TaskHistory

	if updatedTask.Title != nil && *updatedTask.Title != task.Title {
		entries = append(entries, newHistoryEntry(taskID, updaterID, "title", task.Title, *updatedTask.Title))
		task.Title = *updatedTask.Title
	}

	if updatedTask.Description != nil && *updatedTask.Description != task.Description {
		entries = append(entries, newHistoryEntry(taskID, updaterID, "description", task.Description, *updatedTask.Description))
		task.Description = *updatedTask.Description
	}

	if updatedTask.Status != nil && *updatedTask.Status != task.Status {
		entries = append(entries, newHistoryEntry(taskID, updaterID, "status", string(task.Status), string(*updatedTask.Status)))
		task.Status = model.TaskStatus(*updatedTask.Status)
	}

	if updatedTask.AssigneeID != nil {
		oldAssigneeID := ptrToString(task.AssigneeID)
		newAssigneeID := ptrToString(updatedTask.AssigneeID)
		if oldAssigneeID != newAssigneeID {
			entries = append(entries, newHistoryEntry(taskID, updaterID, "assignee_id", oldAssigneeID, newAssigneeID))
			task.AssigneeID = updatedTask.AssigneeID
		}
	}

	if updatedTask.DueDate != nil {
		oldDueDate := fmt.Sprint(task.DueDate)
		newDueDate := fmt.Sprint(updatedTask.DueDate)
		if oldDueDate != newDueDate {
			entries = append(entries, newHistoryEntry(taskID, updaterID, "due_date", oldDueDate, newDueDate))
			task.DueDate = updatedTask.DueDate
		}
	}

	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, err
	}

	for _, e := range entries {
		_ = s.tasks.AddHistory(ctx, e)
	}

	return task, nil
}

func (s *TaskService) GetHistory(ctx context.Context, requesterID, taskID uint64) ([]*model.TaskHistory, error) {
	task, err := s.tasks.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	member, err := s.teams.GetMember(ctx, task.TeamID, requesterID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	return s.tasks.GetHistory(ctx, taskID)
}

func (s *TaskService) OtherTeamAssigneeTasks(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error) {
	return s.tasks.OtherTeamAssigneeTasks(ctx)
}

func newHistoryEntry(taskID, changedBy uint64, field, oldValue, newValue string) *model.TaskHistory {
	return &model.TaskHistory{
		TaskID:    taskID,
		ChangedBy: changedBy,
		Field:     field,
		OldValue:  &oldValue,
		NewValue:  &newValue,
	}
}

func ptrToString(v *uint64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}
