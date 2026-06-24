package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ladislaojs/team-task-manager-test/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, password, name)
		VALUES (?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, name
		FROM users
		WHERE email = ?`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	query := `
		SELECT id, email, password, name
		FROM users
		WHERE id = ?`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) TopTaskCreatorsPerTeam(ctx context.Context) ([]*model.TaskCreator, error) {
	query := `
		WITH task_counts AS (
			SELECT
				ta.team_id,
				ta.created_by                                         AS user_id,
				u.name                                                AS user_name,
				COUNT(*)                                              AS task_count,
				ROW_NUMBER() OVER (
					PARTITION BY ta.team_id
					ORDER BY COUNT(*) DESC
				)                                                     AS rank
			FROM tasks ta
			JOIN users u ON u.id = ta.created_by
			WHERE ta.created_at >= DATE_FORMAT(NOW(), '%Y-%m-01')
			GROUP BY ta.team_id, ta.created_by, u.name
		)
		SELECT team_id, user_id, user_name, task_count, rank
		FROM task_counts
		WHERE rank <= 3
		ORDER BY team_id, rank`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.TaskCreator
	for rows.Next() {
		c := &model.TaskCreator{}
		if err := rows.Scan(&c.TeamID, &c.UserID, &c.UserName, &c.TaskCount, &c.Rank); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, rows.Err()
}
