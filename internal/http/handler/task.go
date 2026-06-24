package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ladislaojs/team-task-manager-test/internal/http/dto"
	"github.com/ladislaojs/team-task-manager-test/internal/http/middleware"
	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
	"github.com/ladislaojs/team-task-manager-test/pkg/response"
)

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TeamID == 0 || req.Title == "" {
		response.Error(w, http.StatusBadRequest, "team_id and title are required")
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "due_date must be RFC3339")
			return
		}
		dueDate = &t
	}

	task, err := h.tasks.Create(r.Context(), claims.UserID, req.TeamID, req.AssigneeID, req.Title, req.Description, dueDate)
	if err != nil {
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Error(w, http.StatusForbidden, "not a team member")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, taskToDTO(task))
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())
	q := r.URL.Query()

	teamID, err := strconv.ParseUint(q.Get("team_id"), 10, 64)
	if err != nil || teamID == 0 {
		response.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}

	filter := model.TaskFilter{TeamID: teamID}

	if s := q.Get("status"); s != "" {
		status := model.TaskStatus(s)
		filter.Status = &status
	}
	if a := q.Get("assignee_id"); a != "" {
		id, err := strconv.ParseUint(a, 10, 64)
		if err == nil {
			filter.AssigneeID = &id
		}
	}
	if p := q.Get("page"); p != "" {
		page, _ := strconv.Atoi(p)
		filter.Page = page
	}
	if ps := q.Get("page_size"); ps != "" {
		pageSize, _ := strconv.Atoi(ps)
		filter.PageSize = pageSize
	}

	tasks, err := h.tasks.List(r.Context(), claims.UserID, filter)
	if err != nil {
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Error(w, http.StatusForbidden, "not a team member")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	taskDTOs := make([]dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = taskToDTO(t)
	}
	response.JSON(w, http.StatusOK, taskDTOs)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	taskID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// todo: history entries

	task, err := h.tasks.Update(r.Context(), claims.UserID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			response.Error(w, http.StatusNotFound, "task not found")
		case errors.Is(err, service.ErrNotTeamMember):
			response.Error(w, http.StatusForbidden, "not a team member")
		case errors.Is(err, service.ErrForbidden):
			response.Error(w, http.StatusForbidden, "forbidden")
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, taskToDTO(task))
}

func (h *TaskHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())

	taskID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	history, err := h.tasks.GetHistory(r.Context(), claims.UserID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			response.Error(w, http.StatusNotFound, "task not found")
		case errors.Is(err, service.ErrNotTeamMember):
			response.Error(w, http.StatusForbidden, "not a team member")
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	historyDTOs := make([]dto.TaskHistoryResponse, len(history))
	for i, h := range history {
		historyDTOs[i] = dto.TaskHistoryResponse{
			ID:        h.ID,
			TaskID:    h.TaskID,
			ChangedBy: h.ChangedBy,
			Field:     h.Field,
			OldValue:  h.OldValue,
			NewValue:  h.NewValue,
			ChangedAt: h.ChangedAt,
		}
	}
	response.JSON(w, http.StatusOK, historyDTOs)
}

func (h *TaskHandler) OtherTeamAssigneeTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.OtherTeamAssigneeTasks(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	taskDTOs := make([]dto.OtherTeamAssigneeTaskResponse, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = dto.OtherTeamAssigneeTaskResponse{
			TaskID:     t.TaskID,
			TeamID:     t.TeamID,
			AssigneeID: t.AssigneeID,
			Title:      t.Title,
		}
	}
	response.JSON(w, http.StatusOK, taskDTOs)
}

func taskToDTO(t *model.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          t.ID,
		TeamID:      t.TeamID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID,
		CreatedBy:   t.CreatedBy,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
	}
}
