package converter

import (
	"github.com/ozoncp/ocp-team-api/internal/models"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
)

func TeamToDTO(team *models.Team) *desc.Team {
	return &desc.Team{
		Id:          team.Id,
		Name:        team.Name,
		Description: team.Description,
	}
}

func TeamFromDTO(dto *desc.Team) *models.Team {
	return &models.Team{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: dto.Description,
	}
}
