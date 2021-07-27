package utils

import (
	"errors"
	"fmt"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"math"
)

func SplitToBulks(teams []models.Team, batchSize uint)[][]models.Team {
	if len(teams) == 0 || batchSize == 0 {
		return [][]models.Team{}
	}

	if int(batchSize) >= len(teams) {
		return [][]models.Team{teams}
	}

	batches := make([][]models.Team, int(math.Ceil(float64(len(teams)) / float64(batchSize))))

	for i := 0; i < cap(batches); i++ {
		if start, end := i*int(batchSize), (i+1)*int(batchSize); end < len(teams) {
			batches[i] = teams[start:end]
		} else {
			batches[i] = teams[start:]
		}
	}

	return batches
}

func TeamsToMap(teams []models.Team) (map[uint64]models.Team, error) {
	teamsMap := make(map[uint64]models.Team, 0)

	for _, team := range teams {
		if _, exists := teamsMap[team.Id]; exists {
			return nil, errors.New(fmt.Sprintf("duplicate ids: team with id=%d already exists", team.Id))
		}

		teamsMap[team.Id] = team
	}

	return teamsMap, nil
}