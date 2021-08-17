package flusher

import (
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	"github.com/ozoncp/ocp-team-api/internal/utils"
)

type Flusher interface {
	Flush(teams []models.Team) []models.Team
}

type flusher struct {
	chunkSize int
	teamRepo repo.Repo
}

func NewFlusher(
	chunkSize int,
	teamRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		teamRepo:  teamRepo,
	}
}

func (f *flusher) Flush(teams []models.Team) []models.Team {
	batches := utils.SplitToBulks(teams, f.chunkSize)

	failed := make([]models.Team, 0)

	for _, chunk := range batches {
		if err := f.teamRepo.AddTeams(chunk); err != nil {
			failed = append(failed, chunk...)
		}
	}

	return failed
}