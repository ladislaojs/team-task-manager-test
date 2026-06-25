package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (team_id, assignee_id, created_by, title, description, status, due_date)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query,
		task.TeamID, task.AssigneeID, task.CreatedBy,
		task.Title, task.Description,
		task.Status, task.DueDate,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = uint64(id)
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uint64) (*model.Task, error) {
	query := `
		SELECT id, team_id, assignee_id, created_by, title, description, status, due_date, created_at
		FROM tasks
		WHERE id = ?`

	t := &model.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.TeamID, &t.AssigneeID, &t.CreatedBy,
		&t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) List(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	where := []string{"team_id = ?"}
	args := []any{filter.TeamID}

	if filter.Status != nil {
		where = append(where, "status = ?")
		args = append(args, *filter.Status)
	}
	if filter.AssigneeID != nil {
		where = append(where, "assignee_id = ?")
		args = append(args, *filter.AssigneeID)
	}

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT id, team_id, assignee_id, created_by, title, description, status, due_date, created_at
		FROM tasks
		WHERE %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, strings.Join(where, " AND "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(
			&t.ID, &t.TeamID, &t.AssigneeID, &t.CreatedBy,
			&t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, assignee_id = ?, due_date = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Status,
		task.AssigneeID, task.DueDate, task.ID,
	)
	return err
}

func (r *TaskRepository) AddHistory(ctx context.Context, entry *model.TaskHistory) error {
	query := `
		INSERT INTO task_history (task_id, changed_by, field, old_value, new_value)
		VALUES (?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query,
		entry.TaskID, entry.ChangedBy, entry.Field, entry.OldValue, entry.NewValue,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	entry.ID = uint64(id)
	return nil
}

func (r *TaskRepository) GetHistory(ctx context.Context, taskID uint64) ([]*model.TaskHistory, error) {
	query := `
		SELECT id, task_id, changed_by, field, old_value, new_value, changed_at
		FROM task_history
		WHERE task_id = ?
		ORDER BY changed_at ASC`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*model.TaskHistory
	for rows.Next() {
		h := &model.TaskHistory{}
		if err := rows.Scan(&h.ID, &h.TaskID, &h.ChangedBy, &h.Field, &h.OldValue, &h.NewValue, &h.ChangedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

func (r *TaskRepository) OtherTeamAssigneeTasks(ctx context.Context) ([]*model.OtherTeamAssigneeTask, error) {
	query := `
		SELECT ta.id, ta.team_id, ta.assignee_id, ta.title
		FROM tasks ta
		WHERE ta.assignee_id IS NOT NULL
		  AND NOT EXISTS (
			  SELECT 1
			  FROM team_members tm
			  WHERE tm.team_id = ta.team_id
			    AND tm.user_id = ta.assignee_id
		  )`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.OtherTeamAssigneeTask
	for rows.Next() {
		t := &model.OtherTeamAssigneeTask{}
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.AssigneeID, &t.Title); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, rows.Err()
}
