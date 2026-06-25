package mysqlrepo_test

import (
	"context"
	"testing"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	mysqlrepo "github.com/ladislaojs/team-task-manager-test/internal/repository/mysql"
)

func TestTeamRepository_Create_And_FindByMemberID(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)

	user := &model.User{Email: "owner@test.com", Password: "h", Name: "Owner"}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	team := &model.Team{Name: "Alpha", CreatedBy: user.ID}
	if err := teamRepo.Create(ctx, team); err != nil {
		t.Fatalf("Create team: %v", err)
	}
	if team.ID == 0 {
		t.Fatal("expected team ID to be set")
	}

	if err := teamRepo.AddMember(ctx, &model.TeamMember{
		TeamID: team.ID, UserID: user.ID, Role: model.TeamRoleOwner,
	}); err != nil {
		t.Fatalf("AddMember: %v", err)
	}

	teams, err := teamRepo.FindByMemberID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByMemberID: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "Alpha" {
		t.Errorf("expected name Alpha, got %s", teams[0].Name)
	}
}

func TestTeamRepository_GetMember_NotFound(t *testing.T) {
	db := newTestDB(t)
	teamRepo := mysqlrepo.NewTeamRepository(db)

	m, err := teamRepo.GetMember(context.Background(), 99, 99)
	if err != nil {
		t.Fatalf("GetMember: %v", err)
	}
	if m != nil {
		t.Error("expected nil for non-existent member")
	}
}

func TestTeamRepository_GetMember_ReturnsRole(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)

	user := &model.User{Email: "admin@test.com", Password: "h", Name: "Admin"}
	userRepo.Create(ctx, user)

	team := &model.Team{Name: "Beta", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)

	teamRepo.AddMember(ctx, &model.TeamMember{
		TeamID: team.ID, UserID: user.ID, Role: model.TeamRoleAdmin,
	})

	m, err := teamRepo.GetMember(ctx, team.ID, user.ID)
	if err != nil {
		t.Fatalf("GetMember: %v", err)
	}
	if m == nil {
		t.Fatal("expected member, got nil")
	}
	if m.Role != model.TeamRoleAdmin {
		t.Errorf("expected role admin, got %s", m.Role)
	}
}

func TestTeamRepository_Stats(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	userRepo := mysqlrepo.NewUserRepository(db)
	teamRepo := mysqlrepo.NewTeamRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	user := &model.User{Email: "u@test.com", Password: "h", Name: "U"}
	userRepo.Create(ctx, user)

	team := &model.Team{Name: "Extended Team", CreatedBy: user.ID}
	teamRepo.Create(ctx, team)
	teamRepo.AddMember(ctx, &model.TeamMember{TeamID: team.ID, UserID: user.ID, Role: model.TeamRoleOwner})

	task := &model.Task{
		TeamID: team.ID, CreatedBy: user.ID,
		Title: "Done task", Status: model.StatusDone,
	}
	taskRepo.Create(ctx, task)
	db.ExecContext(ctx, "UPDATE tasks SET updated_at = NOW() WHERE id = ?", task.ID)

	stats, err := teamRepo.ListExtended(ctx)
	if err != nil {
		t.Fatalf("Stats: %v", err)
	}

	var found *model.TeamExtended
	for _, s := range stats {
		if s.ID == team.ID {
			found = s
			break
		}
	}
	if found == nil {
		t.Fatal("expected extended data for created team")
	}
	if found.MemberCount != 1 {
		t.Errorf("expected 1 member, got %d", found.MemberCount)
	}
	if found.LastWeekDoneTaskCount != 1 {
		t.Errorf("expected 1 done task, got %d", found.LastWeekDoneTaskCount)
	}
}
