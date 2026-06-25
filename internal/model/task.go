package model

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uint64
	TeamID      uint64
	Title       string
	Description string
	AssigneeID  *uint64
	DueDate     *time.Time
	Status      TaskStatus
	CreatedBy   uint64
	CreatedAt   time.Time
}

type TaskFilter struct {
	TeamID     uint64
	Status     *TaskStatus
	AssigneeID *uint64
	Page       int
	PageSize   int
}

type TaskHistory struct {
	ID        uint64
	TaskID    uint64
	Field     string
	OldValue  *string
	NewValue  *string
	ChangedAt time.Time
	ChangedBy uint64
}

type OtherTeamAssigneeTask struct {
	TaskID     uint64
	TeamID     uint64
	AssigneeID uint64
	Title      string
}
