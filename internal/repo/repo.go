package repo

import (
	"context"
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
	ListTeams(ctx context.Context, limit, offset uint64) ([]models.Team, error)
	RemoveTeam(ctx context.Context, teamId uint64) error
	UpdateTeam(ctx context.Context, team models.Team) error
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
	query := sq.Select("*").
		From(tableName).
		Where(sq.Eq{"id": teamId}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	var team models.Team
	if err := query.QueryRowContext(ctx).Scan(&team.Id, &team.Name, &team.Description); err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *repo) ListTeams(ctx context.Context, limit, offset uint64) ([]models.Team, error) {
	query := sq.Select("*").
		From(tableName).
		RunWith(r.db).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Description); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (r *repo) RemoveTeam(ctx context.Context, teamId uint64) error {
	query := sq.Delete(tableName).
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
