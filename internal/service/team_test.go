package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository/mocks"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
)

// noopMailer satisfies email.Service without importing the package.
type noopMailer struct{ err error }

func (m *noopMailer) SendInvitation(_ context.Context, _, _ string) error { return m.err }

func newTeamService(teamRepo *mocks.TeamRepository, userRepo *mocks.UserRepository) *service.TeamService {
	return service.NewTeamService(teamRepo, userRepo, &noopMailer{})
}

func TestTeamService_Create_Success(t *testing.T) {
	var addedMember *model.TeamMember

	teamRepo := &mocks.TeamRepository{
		CreateFn: func(_ context.Context, team *model.Team) error {
			team.ID = 10
			return nil
		},
		AddMemberFn: func(_ context.Context, m *model.TeamMember) error {
			addedMember = m
			return nil
		},
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	team, err := svc.Create(context.Background(), 1, "DevOps")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if team.ID != 10 {
		t.Errorf("expected team ID 10, got %d", team.ID)
	}
	if addedMember == nil {
		t.Fatal("expected creator to be added as member")
	}
	if addedMember.Role != model.TeamRoleOwner {
		t.Errorf("expected owner role, got %s", addedMember.Role)
	}
	if addedMember.UserID != 1 {
		t.Errorf("expected member userID 1, got %d", addedMember.UserID)
	}
}

func TestTeamService_Create_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	teamRepo := &mocks.TeamRepository{
		CreateFn: func(_ context.Context, _ *model.Team) error { return repoErr },
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	_, err := svc.Create(context.Background(), 1, "")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestTeamService_Create_AddMemberError(t *testing.T) {
	addErr := errors.New("add member failed")
	teamRepo := &mocks.TeamRepository{
		CreateFn:    func(_ context.Context, team *model.Team) error { team.ID = 1; return nil },
		AddMemberFn: func(_ context.Context, _ *model.TeamMember) error { return addErr },
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	_, err := svc.Create(context.Background(), 1, "DevOps")

	if !errors.Is(err, addErr) {
		t.Errorf("expected add member error, got %v", err)
	}
}

func TestTeamService_ListForUser(t *testing.T) {
	expected := []*model.Team{{ID: 1, Name: "T"}}
	teamRepo := &mocks.TeamRepository{
		FindByMemberIDFn: func(_ context.Context, _ uint64) ([]*model.Team, error) {
			return expected, nil
		},
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	teams, err := svc.ListForUser(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(teams) != 1 || teams[0].Name != "T" {
		t.Errorf("unexpected result: %v", teams)
	}
}

func TestTeamService_Stats(t *testing.T) {
	expected := []*model.TeamExtended{{ID: 1, MemberCount: 3}}
	teamRepo := &mocks.TeamRepository{
		ListExtendedFn: func(_ context.Context) ([]*model.TeamExtended, error) {
			return expected, nil
		},
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	stats, err := svc.ListExtended(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(stats) != 1 || stats[0].MemberCount != 3 {
		t.Errorf("unexpected stats: %v", stats)
	}
}

func TestTeamService_Invite_Success(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, userID uint64) (*model.TeamMember, error) {
			if userID == 1 {
				return &model.TeamMember{Role: model.TeamRoleOwner}, nil
			}
			return nil, nil
		},
		AddMemberFn:      func(_ context.Context, _ *model.TeamMember) error { return nil },
		FindByMemberIDFn: func(_ context.Context, _ uint64) ([]*model.Team, error) { return nil, nil },
		FindByIDFn:       func(ctx context.Context, id uint64) (*model.Team, error) { return &model.Team{Name: "Test"}, nil },
	}
	userRepo := &mocks.UserRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.User, error) {
			return &model.User{ID: 2, Email: "invited@test.com"}, nil
		},
	}

	svc := newTeamService(teamRepo, userRepo)
	if err := svc.Invite(context.Background(), 1, 10, 2); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTeamService_Invite_RequesterNotMember(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) { return nil, nil },
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	if err := svc.Invite(context.Background(), 99, 10, 2); !errors.Is(err, service.ErrNotTeamMember) {
		t.Errorf("expected ErrNotTeamMember, got %v", err)
	}
}

func TestTeamService_Invite_RequesterNotPrivileged(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleMember}, nil
		},
	}

	svc := newTeamService(teamRepo, &mocks.UserRepository{})
	if err := svc.Invite(context.Background(), 1, 10, 2); !errors.Is(err, service.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestTeamService_Invite_InviteeNotFound(t *testing.T) {
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	userRepo := &mocks.UserRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.User, error) { return nil, nil },
	}

	svc := newTeamService(teamRepo, userRepo)
	if err := svc.Invite(context.Background(), 1, 10, 99); !errors.Is(err, service.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestTeamService_Invite_AlreadyMember(t *testing.T) {
	callCount := 0
	teamRepo := &mocks.TeamRepository{
		GetMemberFn: func(_ context.Context, _, _ uint64) (*model.TeamMember, error) {
			callCount++
			return &model.TeamMember{Role: model.TeamRoleOwner}, nil
		},
	}
	userRepo := &mocks.UserRepository{
		FindByIDFn: func(_ context.Context, _ uint64) (*model.User, error) {
			return &model.User{ID: 2}, nil
		},
	}

	svc := newTeamService(teamRepo, userRepo)
	if err := svc.Invite(context.Background(), 1, 10, 2); !errors.Is(err, service.ErrAlreadyMember) {
		t.Errorf("expected ErrAlreadyMember, got %v", err)
	}
}
