package flusher

import (
	"context"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	"github.com/ozoncp/ocp-team-api/internal/utils"
)

// Flusher is the interface for flushing teams into repo.
type Flusher interface {
	Flush(ctx context.Context, teams []models.Team) []models.Team
}

// flusher is the struct that implements Flusher interface.
type flusher struct {
	chunkSize int
	teamRepo  repo.Repo
}

// NewFlusher is the constructor method for flusher struct.
func NewFlusher(
	chunkSize int,
	teamRepo repo.Repo,
) *flusher {
	return &flusher{
		chunkSize: chunkSize,
		teamRepo:  teamRepo,
	}
}

// Flush is the method that creates new teams using repo.Repo batch-by-batch.
func (f *flusher) Flush(ctx context.Context, teams []models.Team) []models.Team {
	batches := utils.SplitToBulks(teams, f.chunkSize)

	failed := make([]models.Team, 0)

	for _, chunk := range batches {
		if _, err := f.teamRepo.CreateTeams(ctx, chunk); err != nil {
			failed = append(failed, chunk...)
		}
	}

	return failed
}
