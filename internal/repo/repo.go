package repo

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-team-api/internal/models"
)

const (
	tableName = "team"
)

type Repo interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	CreateTeams(ctx context.Context, teams []models.Team) ([]uint64, error)
	GetTeam(ctx context.Context, teamId uint64) (*models.Team, error)
	ListTeams(ctx context.Context, limit, offset uint64) ([]models.Team, uint64, error)
	RemoveTeam(ctx context.Context, teamId uint64) error
	UpdateTeam(ctx context.Context, team models.Team) error
	SearchTeams(ctx context.Context, query string, searchType uint8) ([]models.Team, error)
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreateTeam(ctx context.Context, team *models.Team) error {
	query := sq.Insert(tableName).
		Columns("name", "description").
		Values(team.Name, team.Description).
		Suffix("RETURNING id").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRowContext(ctx).Scan(&team.Id)

	return err
}

func (r *repo) CreateTeams(ctx context.Context, teams []models.Team) ([]uint64, error) {
	query := sq.Insert(tableName).
		Columns("name", "description").
		Suffix("RETURNING id").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	for _, team := range teams {
		query = query.Values(team.Name, team.Description)
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(teams))
	for rows.Next() {
		var id uint64

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *repo) GetTeam(ctx context.Context, teamId uint64) (*models.Team, error) {
	query := sq.Select("id", "name", "description").
		From(tableName).
		Where(sq.And{
			sq.Eq{"id": teamId},
			sq.Eq{"is_deleted": false},
		}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	var team models.Team
	if err := query.QueryRowContext(ctx).Scan(&team.Id, &team.Name, &team.Description); err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *repo) ListTeams(ctx context.Context, limit, offset uint64) ([]models.Team, uint64, error) {
	query := sq.Select("id", "name", "description").
		From(tableName).
		Where(sq.Eq{"is_deleted": false}).
		RunWith(r.db).
		OrderBy("id").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Description); err != nil {
			return nil, 0, err
		}

		teams = append(teams, team)
	}

	var total uint64
	query = sq.Select("COUNT(*)").
		From(tableName).
		Where(sq.Eq{"is_deleted": false}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)
	err = query.QueryRowContext(ctx).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *repo) RemoveTeam(ctx context.Context, teamId uint64) error {
	query := sq.Update(tableName).
		Set("is_deleted", true).
		Where(sq.Eq{"id": teamId}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	_, err := query.ExecContext(ctx)
	return err
}

func (r *repo) UpdateTeam(ctx context.Context, team models.Team) error {
	query := sq.Update(tableName).
		Set("name", team.Name).
		Set("description", team.Description).
		Where(sq.Eq{"id": team.Id}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	_, err := query.ExecContext(ctx)

	return err
}

func (r *repo) SearchTeams(ctx context.Context, query string, searchType uint8) ([]models.Team, error) {
	plainTextQuery := `SELECT id, ts_headline(name, q), ts_headline(description, q) FROM team, 
			plainto_tsquery($1) AS q WHERE is_deleted = FALSE AND tsv @@ q ORDER BY ts_rank(tsv, q) DESC`

	phraseTextQuery := `SELECT id, ts_headline(name, q), ts_headline(description, q) FROM team, 
			phraseto_tsquery($1) AS q WHERE is_deleted = FALSE AND tsv @@ q ORDER BY ts_rank(tsv, q) DESC`

	var rows *sql.Rows
	var err error
	if searchType == uint8(0) {
		rows, err = r.db.QueryContext(ctx, plainTextQuery, query)
	} else {
		rows, err = r.db.QueryContext(ctx, phraseTextQuery, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err = rows.Scan(&team.Id, &team.Name, &team.Description); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}
