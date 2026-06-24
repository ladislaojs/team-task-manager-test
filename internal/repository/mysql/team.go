package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *model.Team) error
func (r *TeamRepository) FindByMemberId(ctx context.Context, userID *uint64)
func (r *TeamRepository) AddMember(ctx context.Context, member *model.TeamMember) error
