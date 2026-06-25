package mysqlrepo_test

import (
	"context"
	"testing"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	mysqlrepo "github.com/ladislaojs/team-task-manager-test/internal/repository/mysql"
)

func TestUserRepository_Create_And_FindByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := mysqlrepo.NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Email:    "alice@test.com",
		Password: "hashed",
		Name:     "Alice",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected ID to be set after Create")
	}

	found, err := repo.FindByEmail(ctx, "alice@test.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, found.ID)
	}
	if found.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", found.Name)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := mysqlrepo.NewUserRepository(db)

	found, err := repo.FindByEmail(context.Background(), "ghost@test.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if found != nil {
		t.Error("expected nil for non-existent email")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := mysqlrepo.NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{Email: "bob@test.com", Password: "hashed", Name: "Bob"}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found == nil || found.Email != "bob@test.com" {
		t.Errorf("expected bob@test.com, got %v", found)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := mysqlrepo.NewUserRepository(db)

	found, err := repo.FindByID(context.Background(), 99999)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found != nil {
		t.Error("expected nil for non-existent ID")
	}
}

func TestUserRepository_TopTaskCreatorsPerTeam(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	// seed: 1 user, 1 team, 2 tasks created this month
	user := &model.User{Email: "creator@test.com", Password: "h", Name: "Creator"}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	team := &model.Team{Name: "Team A", CreatedBy: user.ID}
	if err := teamRepo.Create(ctx, team); err != nil {
		t.Fatalf("create team: %v", err)
	}
	teamRepo.AddMember(ctx, &model.TeamMember{TeamID: team.ID, UserID: user.ID, Role: model.TeamRoleOwner})

	for i := 0; i < 2; i++ {
		task := &model.Task{
			TeamID: team.ID, CreatedBy: user.ID,
			Title: "Task", Status: model.StatusTodo,
		}
		if err := taskRepo.Create(ctx, task); err != nil {
			t.Fatalf("create task: %v", err)
		}
	}

	creators, err := userRepo.TopTaskCreatorsPerTeam(ctx)
	if err != nil {
		t.Fatalf("TopTaskCreatorsPerTeam: %v", err)
	}
	if len(creators) == 0 {
		t.Fatal("expected at least one result")
	}
	if creators[0].TaskCount != 2 {
		t.Errorf("expected task count 2, got %d", creators[0].TaskCount)
	}
	if creators[0].Rank != 1 {
		t.Errorf("expected rank 1, got %d", creators[0].Rank)
	}
}
