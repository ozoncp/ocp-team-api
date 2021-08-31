package converter

import (
	"github.com/ozoncp/ocp-team-api/internal/models"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
)

// TeamToDTO is the method for converting
// inner team model (models.Team) into
// protobuf-generated data transport object.
func TeamToDTO(team *models.Team) *desc.Team {
	return &desc.Team{
		Id:          team.Id,
		Name:        team.Name,
		Description: team.Description,
	}
}

// TeamFromDTO is the method for converting
// protobuf-generated data transport object
// into inner team model (models.Team).
func TeamFromDTO(dto *desc.Team) *models.Team {
	return &models.Team{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: dto.Description,
	}
}
