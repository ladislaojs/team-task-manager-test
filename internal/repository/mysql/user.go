package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error)
func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error)
