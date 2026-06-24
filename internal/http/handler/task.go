package handler

import "github.com/ladislaojs/team-task-manager-test/internal/service"

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}
