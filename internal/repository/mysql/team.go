package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *model.Team) error {
	query := `
		INSERT INTO teams (name, created_by)
		VALUES (?, ?)`

	res, err := r.db.ExecContext(ctx, query, team.Name, team.CreatedBy)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	team.ID = uint64(id)
	return nil
}

func (r *TeamRepository) FindByMemberID(ctx context.Context, userID uint64) ([]*model.Team, error) {
	query := `
		SELECT t.id, t.name, t.created_by
		FROM teams t
		JOIN team_members tm ON tm.team_id = t.id
		WHERE tm.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*model.Team
	for rows.Next() {
		t := &model.Team{}
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedBy); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	return teams, rows.Err()
}

func (r *TeamRepository) ListExtended(ctx context.Context) ([]*model.TeamExtended, error) {
	query := `
		SELECT
			t.id,
			t.name,
			COUNT(DISTINCT tm.user_id)                                AS member_count,
			COUNT(DISTINCT CASE
				WHEN ta.status = 'done'
				 AND ta.updated_at >= NOW() - INTERVAL 7 DAY
				THEN ta.id
			END)                                                      AS last_week_done_task_count
		FROM teams t
		LEFT JOIN team_members tm ON tm.team_id = t.id
		LEFT JOIN tasks        ta ON ta.team_id = t.id
		GROUP BY t.id, t.name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.TeamExtended
	for rows.Next() {
		s := &model.TeamExtended{}
		if err := rows.Scan(&s.ID, &s.Name, &s.MemberCount, &s.LastWeekDoneTaskCount); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, rows.Err()
}

func (r *TeamRepository) FindByID(ctx context.Context, id uint64) (*model.Team, error) {
	query := `
		SELECT id, name
		FROM teams
		WHERE id = ?`

	t := &model.Team{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, member *model.TeamMember) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, member.TeamID, member.UserID, member.Role)
	return err
}

func (r *TeamRepository) GetMember(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error) {
	query := `
		SELECT team_id, user_id, role
		FROM team_members
		WHERE team_id = ? AND user_id = ?`

	m := &model.TeamMember{}
	err := r.db.QueryRowContext(ctx, query, teamID, userID).Scan(
		&m.TeamID, &m.UserID, &m.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}
