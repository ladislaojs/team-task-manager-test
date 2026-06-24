package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ladislaojs/team-task-manager-test/internal/model"
	"github.com/ladislaojs/team-task-manager-test/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

var (
	ErrEmailTaken        = errors.New("email is already taken")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserNotFound      = errors.New("user does not exist")
)

type UserService struct {
	users     repository.UserRepository
	jwtSecret []byte
}

func NewUserService(users repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{users: users, jwtSecret: []byte(jwtSecret)}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (*model.User, error) {
	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.Tokens, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrIncorrectPassword
	}

	accessToken, err := s.generateToken(user.ID, accessTokenTTL, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, refreshTokenTTL, "refresh")
	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) TopTaskCreatorsPerTeam(ctx context.Context) ([]*model.TaskCreator, error) {
	return s.users.TopTaskCreatorsPerTeam(ctx)
}

func (s *UserService) generateToken(userID uint64, ttl time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    tokenType,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
