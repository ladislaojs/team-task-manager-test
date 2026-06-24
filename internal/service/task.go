package service

import "github.com/ladislaojs/team-task-manager-test/internal/repository"

type TaskService struct {
	tasks repository.TaskRepository
}

func NewTaskService(tasks repository.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks}
}
