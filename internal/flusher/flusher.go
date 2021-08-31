package flusher

import (
	"context"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	"github.com/ozoncp/ocp-team-api/internal/utils"
)

// IFlusher is the interface for flushing teams into repo.
type IFlusher interface {
	Flush(ctx context.Context, teams []models.Team) []models.Team
}

// Flusher is the struct that implements IFlusher interface.
type Flusher struct {
	chunkSize int
	teamRepo  repo.IRepo
}

// NewFlusher is the constructor method for Flusher struct.
func NewFlusher(
	chunkSize int,
	teamRepo repo.IRepo,
) *Flusher {
	return &Flusher{
		chunkSize: chunkSize,
		teamRepo:  teamRepo,
	}
}

// Flush is the method that creates new teams using repo.IRepo batch-by-batch.
func (f *Flusher) Flush(ctx context.Context, teams []models.Team) []models.Team {
	batches := utils.SplitToBulks(teams, f.chunkSize)

	failed := make([]models.Team, 0)

	for _, chunk := range batches {
		if _, err := f.teamRepo.CreateTeams(ctx, chunk); err != nil {
			failed = append(failed, chunk...)
		}
	}

	return failed
}
