package mysqlrepo_test

import (
	"context"
	"testing"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	mysqlrepo "github.com/ladislaojs/team-task-manager-test/internal/repository/mysql"
)

func TestTaskRepository_Create_And_FindByID(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	user := &model.User{Email: "u@test.com", Password: "h", Name: "U"}
	userRepo.Create(ctx, user)
	team := &model.Team{Name: "T", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)

	task := &model.Task{
		TeamID:    team.ID,
		CreatedBy: user.ID,
		Title:     "Fix login",
		Status:    model.StatusTodo,
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if task.ID == 0 {
		t.Fatal("expected task ID to be set")
	}

	found, err := taskRepo.FindByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found == nil {
		t.Fatal("expected task, got nil")
	}
	if found.Title != "Fix login" {
		t.Errorf("expected title 'Fix login', got %s", found.Title)
	}
}

func TestTaskRepository_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	found, err := mysqlrepo.NewTaskRepository(db).FindByID(context.Background(), 99999)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found != nil {
		t.Error("expected nil for non-existent task")
	}
}

func TestTaskRepository_List_Filter(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	user := &model.User{Email: "u@test.com", Password: "h", Name: "U"}
	userRepo.Create(ctx, user)
	team := &model.Team{Name: "T", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)

	// 2 todo (first assigned to user), 1 done
	taskRepo.Create(ctx, &model.Task{
		TeamID: team.ID, CreatedBy: user.ID, AssigneeID: &user.ID,
		Title: "T1", Status: model.StatusTodo,
	})
	taskRepo.Create(ctx, &model.Task{
		TeamID: team.ID, CreatedBy: user.ID,
		Title: "T2", Status: model.StatusTodo,
	})
	taskRepo.Create(ctx, &model.Task{
		TeamID: team.ID, CreatedBy: user.ID,
		Title: "T3", Status: model.StatusDone,
	})

	todoStatus := model.StatusTodo
	tasks, err := taskRepo.List(ctx, model.TaskFilter{
		TeamID: team.ID, Status: &todoStatus, Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("List by status: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 todo tasks, got %d", len(tasks))
	}

	tasks, err = taskRepo.List(ctx, model.TaskFilter{
		TeamID: team.ID, AssigneeID: &user.ID, Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("List by assignee: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 assigned task, got %d", len(tasks))
	}
}

func TestTaskRepository_Update(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	user := &model.User{Email: "u@test.com", Password: "h", Name: "U"}
	userRepo.Create(ctx, user)
	team := &model.Team{Name: "T", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)

	task := &model.Task{
		TeamID: team.ID, CreatedBy: user.ID,
		Title: "Old", Status: model.StatusTodo,
	}
	taskRepo.Create(ctx, task)

	task.Title = "New"
	task.Status = model.StatusDone
	if err := taskRepo.Update(ctx, task); err != nil {
		t.Fatalf("Update: %v", err)
	}

	found, _ := taskRepo.FindByID(ctx, task.ID)
	if found.Title != "New" {
		t.Errorf("expected title 'New', got %s", found.Title)
	}
	if found.Status != model.StatusDone {
		t.Errorf("expected status done, got %s", found.Status)
	}
}

func TestTaskRepository_History(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	user := &model.User{Email: "u@test.com", Password: "h", Name: "U"}
	userRepo.Create(ctx, user)
	team := &model.Team{Name: "T", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)
	task := &model.Task{
		TeamID: team.ID, CreatedBy: user.ID,
		Title: "T", Status: model.StatusTodo,
	}
	taskRepo.Create(ctx, task)

	oldVal, newVal := "todo", "done"
	entry := &model.TaskHistory{
		TaskID:    task.ID,
		ChangedBy: user.ID,
		Field:     "status",
		OldValue:  &oldVal,
		NewValue:  &newVal,
	}
	if err := taskRepo.AddHistory(ctx, entry); err != nil {
		t.Fatalf("AddHistory: %v", err)
	}
	if entry.ID == 0 {
		t.Fatal("expected history entry ID to be set")
	}

	history, err := taskRepo.GetHistory(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}
	if history[0].Field != "status" {
		t.Errorf("expected field 'status', got %s", history[0].Field)
	}
	if *history[0].OldValue != "todo" {
		t.Errorf("expected old value 'todo', got %s", *history[0].OldValue)
	}
}

func TestTaskRepository_OtherTeamAssigneeTasks(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	owner := &model.User{Email: "owner@test.com", Password: "h", Name: "Owner"}
	alien := &model.User{Email: "alien@test.com", Password: "h", Name: "Alien"}
	userRepo.Create(ctx, owner)
	userRepo.Create(ctx, alien)

	team := &model.Team{Name: "T", CreatedBy: owner.ID}
	teamRepo.Create(ctx, team)
	teamRepo.AddMember(ctx, &model.TeamMember{TeamID: team.ID, UserID: owner.ID, Role: model.TeamRoleOwner})

	taskRepo.Create(ctx, &model.Task{
		TeamID: team.ID, CreatedBy: owner.ID, AssigneeID: &alien.ID,
		Title: "Alien", Status: model.StatusTodo,
	})
	taskRepo.Create(ctx, &model.Task{
		TeamID: team.ID, CreatedBy: owner.ID, AssigneeID: &owner.ID,
		Title: "Owner", Status: model.StatusTodo,
	})

	otherTeamTasks, err := taskRepo.OtherTeamAssigneeTasks(ctx)
	if err != nil {
		t.Fatalf("OtherTeamAssigneeTasks: %v", err)
	}
	if len(otherTeamTasks) != 1 {
		t.Fatalf("expected 1 other team assignee task, got %d", len(otherTeamTasks))
	}
	if otherTeamTasks[0].Title != "Alien" {
		t.Errorf("expected title 'Alien', got %s", otherTeamTasks[0].Title)
	}
}
