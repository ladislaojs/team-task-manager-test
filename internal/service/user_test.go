package service_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository/mocks"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
)

const testJWTSecret = "test-secret"

func TestUserService_Register_Success(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, nil },
		CreateFn:      func(_ context.Context, u *model.User) error { u.ID = 1; return nil },
	}

	svc := service.NewUserService(repo, testJWTSecret)
	user, err := svc.Register(context.Background(), "ladislaojs@test.com", "password123", "Ladislao")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Password == "password123" {
		t.Error("password must be hashed")
	}
}

func TestUserService_Register_EmailTaken(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: 1}, nil
		},
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Register(context.Background(), "ladislaojs@test.com", "password123", "Ladislao")

	if !errors.Is(err, service.ErrEmailTaken) {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_Register_FindEmailError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, repoErr },
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Register(context.Background(), "ladislaojs@test.com", "password123", "Ladislao")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestUserService_Register_CreateError(t *testing.T) {
	createErr := errors.New("insert failed")
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, nil },
		CreateFn:      func(_ context.Context, _ *model.User) error { return createErr },
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Register(context.Background(), "ladislaojs@test.com", "password123", "Ladislao")

	if !errors.Is(err, createErr) {
		t.Errorf("expected create error, got %v", err)
	}
}

func TestUserService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: 1, Password: string(hash)}, nil
		},
	}

	svc := service.NewUserService(repo, testJWTSecret)
	tokens, err := svc.Login(context.Background(), "ladislaojs@test.com", "password123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Error("expected non-empty tokens")
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, nil },
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Login(context.Background(), "notfound@test.com", "password123")

	if !errors.Is(err, service.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_Login_FindEmailError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) { return nil, repoErr },
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Login(context.Background(), "ladislaojs@test.com", "pass")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: 1, Password: string(hash)}, nil
		},
	}

	svc := service.NewUserService(repo, testJWTSecret)
	_, err := svc.Login(context.Background(), "ladislaojs@test.com", "wrong")

	if !errors.Is(err, service.ErrIncorrectPassword) {
		t.Errorf("expected ErrIncorrectPassword, got %v", err)
	}
}

func TestUserService_TopTaskCreatorsPerTeam(t *testing.T) {
	expected := []*model.TaskCreator{{TeamID: 1, UserID: 2, TaskCount: 5, Rank: 1}}
	repo := &mocks.UserRepository{
		TopTaskCreatorsPerTeamFn: func(_ context.Context) ([]*model.TaskCreator, error) {
			return expected, nil
		},
	}

	svc := service.NewUserService(repo, testJWTSecret)
	creators, err := svc.TopTaskCreatorsPerTeam(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(creators) != 1 || creators[0].TaskCount != 5 {
		t.Errorf("unexpected result: %v", creators)
	}
}
