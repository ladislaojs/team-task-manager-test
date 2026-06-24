package model

import "time"

type TaskID uint64
type TaskHistoryId uint64

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          TaskID
	TeamID      TeamID
	Title       string
	Description string
	AssigneeID  UserID
	DueDate     time.Time
	Status      TaskStatus
	CreatedBy   UserID
}

type TaskHistory struct {
	ID        TaskHistoryId
	TaskID    TaskID
	Field     string
	OldValue  string
	NewValue  string
	ChangedAt time.Time
	ChangedBy UserID
}
