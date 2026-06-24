package service

import "github.com/ladislaojs/team-task-manager-test/internal/repository"

type UserService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Register() {

}

func (s *UserService) Login() {

}
