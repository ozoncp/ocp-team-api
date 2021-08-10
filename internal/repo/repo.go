package repo

import "github.com/ozoncp/ocp-team-api/internal/models"

type Repo interface {
	AddTeams(teams []models.Team) error
	ListTeams(limit, offset uint64) ([]models.Team, error)
	DescribeTeam(teamId uint64) (*models.Team, error)
}