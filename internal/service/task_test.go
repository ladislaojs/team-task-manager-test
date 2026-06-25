package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository/mocks"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
)

func TestTaskService_Create_Success(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		CreateFn: func(_ context.Context, task *model.Task) error { task.ID = 5; return nil },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	task, err := svc.Create(context.Background(), 1, 10, nil, "Fix bug", "", nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.ID != 5 {
		t.Errorf("expected task ID 5, got %d", task.ID)
	}
	if task.Status != model.StatusTodo {
		t.Errorf("expected status todo, got %s", task.Status)
	}
}

func TestTaskService_Create_NotMember(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) { return nil, nil },
	}

	svc := service.NewTaskService(&mocks.TaskRepository{}, teamRepo)
	_, err := svc.Create(context.Background(), 1, 10, nil, "Fix bug", "", nil)

	if !errors.Is(err, service.ErrNotTeamMember) {
		t.Errorf("expected ErrNotTeamMember, got %v", err)
	}
}

func TestTaskService_Create_RepoError(t *testing.T) {
	createErr := errors.New("insert failed")
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		CreateFn: func(_ context.Context, _ *model.Task) error { return createErr },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.Create(context.Background(), 1, 10, nil, "T", "", nil)

	if !errors.Is(err, createErr) {
		t.Errorf("expected create error, got %v", err)
	}
}

func TestTaskService_List_DefaultsPagination(t *testing.T) {
	var capturedFilter model.TaskFilter

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		ListFn: func(_ context.Context, f model.TaskFilter) ([]*model.Task, error) {
			capturedFilter = f
			return nil, nil
		},
	}

	service.NewTaskService(taskRepo, teamRepo).List(context.Background(), 1, model.TaskFilter{TeamID: 10})

	if capturedFilter.Page != 1 {
		t.Errorf("expected default page 1, got %d", capturedFilter.Page)
	}
	if capturedFilter.PageSize != 20 {
		t.Errorf("expected default page size 20, got %d", capturedFilter.PageSize)
	}
}

func TestTaskService_List_NotMember(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) { return nil, nil },
	}

	svc := service.NewTaskService(&mocks.TaskRepository{}, teamRepo)
	_, err := svc.List(context.Background(), 1, model.TaskFilter{TeamID: 10})

	if !errors.Is(err, service.ErrNotTeamMember) {
		t.Errorf("expected ErrNotTeamMember, got %v", err)
	}
}

func TestTaskService_List_RepoError(t *testing.T) {
	listErr := errors.New("db error")
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		ListFn: func(_ context.Context, _ model.TaskFilter) ([]*model.Task, error) { return nil, listErr },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.List(context.Background(), 1, model.TaskFilter{TeamID: 10})

	if !errors.Is(err, listErr) {
		t.Errorf("expected list error, got %v", err)
	}
}

func TestTaskService_Update_Success_RecordsHistory(t *testing.T) {
	assigneeID := uint64(3)
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 1,
		Title: "Old title", Status: model.StatusTodo,
		AssigneeID: &assigneeID,
	}
	var historyEntries []*model.TaskHistory

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		UpdateFn:   func(_ context.Context, _ *model.Task) error { return nil },
		AddHistoryFn: func(_ context.Context, e *model.TaskHistory) error {
			historyEntries = append(historyEntries, e)
			return nil
		},
	}

	newTitle := "New title"
	newStatus := model.StatusInProgress

	svc := service.NewTaskService(taskRepo, teamRepo)
	task, err := svc.Update(context.Background(), 1, 5, service.UpdatedTask{
		Title:      &newTitle,
		Status:     &newStatus,
		AssigneeID: &assigneeID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Title != "New title" {
		t.Errorf("expected 'New title', got %s", task.Title)
	}
	if len(historyEntries) != 2 {
		t.Errorf("expected 2 history entries, got %d", len(historyEntries))
	}
}

func TestTaskService_Update_AssignedTo_RecordsHistory(t *testing.T) {
	oldAssignee := uint64(3)
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 1,
		Title: "T", Status: model.StatusTodo,
		AssigneeID: &oldAssignee,
	}
	var historyEntries []*model.TaskHistory

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		UpdateFn:   func(_ context.Context, _ *model.Task) error { return nil },
		AddHistoryFn: func(_ context.Context, e *model.TaskHistory) error {
			historyEntries = append(historyEntries, e)
			return nil
		},
	}

	newAssignee := uint64(7)
	svc := service.NewTaskService(taskRepo, teamRepo)
	svc.Update(context.Background(), 1, 5, service.UpdatedTask{AssigneeID: &newAssignee})

	if len(historyEntries) != 1 || historyEntries[0].Field != "assignee_id" {
		t.Errorf("expected 1 assignee_id history entry, got %d entries", len(historyEntries))
	}
}

func TestTaskService_Update_DueDate_RecordsHistory(t *testing.T) {
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 1,
		Title: "T", Status: model.StatusTodo,
	}
	var historyEntries []*model.TaskHistory

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		UpdateFn:   func(_ context.Context, _ *model.Task) error { return nil },
		AddHistoryFn: func(_ context.Context, e *model.TaskHistory) error {
			historyEntries = append(historyEntries, e)
			return nil
		},
	}

	due := time.Now().Add(24 * time.Hour)
	svc := service.NewTaskService(taskRepo, teamRepo)
	svc.Update(context.Background(), 1, 5, service.UpdatedTask{DueDate: &due})

	if len(historyEntries) != 1 || historyEntries[0].Field != "due_date" {
		t.Errorf("expected 1 due_date history entry, got %d", len(historyEntries))
	}
}

func TestTaskService_Update_TaskNotFound(t *testing.T) {
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return nil, nil },
	}

	svc := service.NewTaskService(taskRepo, &mocks.TeamRepository{})
	_, err := svc.Update(context.Background(), 1, 99, service.UpdatedTask{})

	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskService_Update_Forbidden(t *testing.T) {
	otherUser := uint64(99)
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 2, AssigneeID: &otherUser,
		Title: "T", Status: model.StatusTodo,
	}

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.Update(context.Background(), 1, 5, service.UpdatedTask{})

	if !errors.Is(err, service.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestTaskService_Update_NoChanges_NoHistory(t *testing.T) {
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 1,
		Title: "Same", Status: model.StatusTodo,
	}
	var historyCount int

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn:   func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		UpdateFn:     func(_ context.Context, _ *model.Task) error { return nil },
		AddHistoryFn: func(_ context.Context, _ *model.TaskHistory) error { historyCount++; return nil },
	}

	sameTitle := "Same"
	service.NewTaskService(taskRepo, teamRepo).Update(context.Background(), 1, 5, service.UpdatedTask{Title: &sameTitle})

	if historyCount != 0 {
		t.Errorf("expected no history for unchanged fields, got %d", historyCount)
	}
}

func TestTaskService_Update_RepoUpdateError(t *testing.T) {
	updateErr := errors.New("update failed")
	existing := &model.Task{
		ID: 5, TeamID: 10, CreatedBy: 1,
		Title: "T", Status: model.StatusTodo,
	}

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		UpdateFn:   func(_ context.Context, _ *model.Task) error { return updateErr },
	}

	newTitle := "New"
	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.Update(context.Background(), 1, 5, service.UpdatedTask{Title: &newTitle})

	if !errors.Is(err, updateErr) {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestTaskService_OtherTeamAssigneeTasks(t *testing.T) {
	expected := []*model.OtherTeamAssigneeTask{{TaskID: 1, TeamID: 2, AssigneeID: 3, Title: "T"}}
	taskRepo := &mocks.TaskRepository{
		OtherTeamAssigneeTasksFn: func(_ context.Context) ([]*model.OtherTeamAssigneeTask, error) { return expected, nil },
	}

	svc := service.NewTaskService(taskRepo, &mocks.TeamRepository{})
	result, err := svc.OtherTeamAssigneeTasks(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}

func TestTaskService_GetHistory_Success(t *testing.T) {
	existing := &model.Task{ID: 5, TeamID: 10}
	history := []*model.TaskHistory{{ID: 1, Field: "status"}}

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn:   func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
		GetHistoryFn: func(_ context.Context, _ uint64) ([]*model.TaskHistory, error) { return history, nil },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	result, err := svc.GetHistory(context.Background(), 1, 5)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(result))
	}
}

func TestTaskService_GetHistory_NotMember(t *testing.T) {
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) {
			return &model.Task{ID: 5, TeamID: 10}, nil
		},
	}
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) { return nil, nil },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.GetHistory(context.Background(), 1, 5)

	if !errors.Is(err, service.ErrNotTeamMember) {
		t.Errorf("expected ErrNotTeamMember, got %v", err)
	}
}

func TestTaskService_Update_GetMemberError(t *testing.T) {
	memberErr := errors.New("db error")
	existing := &model.Task{ID: 5, TeamID: 10, CreatedBy: 1, Title: "T", Status: model.StatusTodo}

	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) { return nil, memberErr },
	}
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return existing, nil },
	}

	svc := service.NewTaskService(taskRepo, teamRepo)
	_, err := svc.Update(context.Background(), 1, 5, service.UpdatedTask{})

	if !errors.Is(err, memberErr) {
		t.Errorf("expected member repo error, got %v", err)
	}
}

func TestTaskService_GetHistory_TaskNotFound(t *testing.T) {
	taskRepo := &mocks.TaskRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.Task, error) { return nil, nil },
	}

	svc := service.NewTaskService(taskRepo, &mocks.TeamRepository{})
	_, err := svc.GetHistory(context.Background(), 1, 99)

	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
