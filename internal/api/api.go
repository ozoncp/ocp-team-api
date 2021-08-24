package api

import (
	"context"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type api struct {
	desc.UnimplementedOcpTeamApiServer
	repo repo.Repo
}

func NewOcpTeamApi(repo repo.Repo) desc.OcpTeamApiServer {
	return &api{
		repo: repo,
	}
}

func (a *api) CreateTeamV1(
	ctx context.Context,
	req *desc.CreateTeamV1Request) (*desc.CreateTeamV1Response, error) {
	log.Printf("Create team (name=%s, description=%s)", req.Name, req.Description)

	id, err := a.repo.AddTeam(ctx, models.Team{Name: req.Name, Description: req.Description})

	if err != nil {
		log.Warn().Msgf("failed to create team %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msgf("new team was created successfully with id=%d", id)

	return &desc.CreateTeamV1Response{Id: strconv.FormatUint(id, 10)}, nil
}

func (a *api) GetTeamV1(
	ctx context.Context,
	req *desc.GetTeamV1Request) (*desc.GetTeamV1Response, error) {
	log.Printf("Get team (id=%d)", req.Id)

	team, err := a.repo.GetTeam(ctx, req.Id)

	if err != nil {
		log.Warn().Msgf("failed to fetch team with id=%d: %v", req.Id, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	response := &desc.GetTeamV1Response{
		Team: &desc.Team{
			Id: team.Id,
			Name: team.Name,
			Description: team.Description,
		},
	}

	return response, nil
}

func (a *api) ListTeamsV1(
	ctx context.Context,
	req *desc.ListTeamsV1Request) (*desc.ListTeamsV1Response, error) {
	log.Printf("List teams (limit=%d, offset=%d)", req.Limit, req.Offset)

	teams, err := a.repo.ListTeams(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Warn().Msgf("failed to fetch teams: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseTeams := make([]*desc.Team, len(teams), cap(teams))
	for i, team := range teams {
		responseTeams[i] = &desc.Team{
			Id: team.Id,
			Name: team.Name,
			Description: team.Description,
		}
	}

	return &desc.ListTeamsV1Response{Teams: responseTeams}, nil
}

func (a *api) RemoveTeamV1(
	ctx context.Context,
	req *desc.RemoveTeamV1Request) (*desc.RemoveTeamV1Response, error) {
	log.Printf("Remove team (id=%d)", req.Id)

	err := a.repo.RemoveTeam(ctx, req.Id)
	if err != nil {
		log.Warn().Msgf("failed to delete team with id=%d", req.Id)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.RemoveTeamV1Response{}, nil
}