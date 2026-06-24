package dto

import "time"

type CreateTaskRequest struct {
	TeamID      uint64 `json:"team_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  uint64 `json:"assignee_id"`
	DueDate     string `json:"due_date"`
}

type TaskResponse struct {
	ID          uint64    `json:"id"`
	TeamID      uint64    `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	AssigneeID  uint64    `json:"assignee_id"`
	CreatedBy   uint64    `json:"created_by"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssignedTo  uint64 `json:"assigned_to"`
	DueDate     string `json:"due_date"`
}

type TaskHistoryResponse struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"task_id"`
	ChangedBy uint64    `json:"changed_by"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

type OtherTeamAssigneeTaskResponse struct {
	TaskResponse
}
