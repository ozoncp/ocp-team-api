package api

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-team-api/internal/config"
	"github.com/ozoncp/ocp-team-api/internal/converter"
	"github.com/ozoncp/ocp-team-api/internal/kafka"
	"github.com/ozoncp/ocp-team-api/internal/metrics"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	"github.com/ozoncp/ocp-team-api/internal/utils"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// api is the struct that implements protobuf-interface.
type api struct {
	desc.UnimplementedOcpTeamApiServer
	repo     repo.Repo
	producer kafka.Producer
}

// NewOcpTeamApi is the constructor method for api struct.
func NewOcpTeamApi(repo repo.Repo, producer kafka.Producer) *api {
	return &api{
		repo:     repo,
		producer: producer,
	}
}

// CreateTeamV1 is the method that handles creating new team.
func (a *api) CreateTeamV1(
	ctx context.Context,
	req *desc.CreateTeamV1Request) (*desc.CreateTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("CreateTeamV1() was called (name=%s, description=%s)", req.Name, req.Description)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("CreateTeamV1")
	defer span.Finish()

	team := models.Team{Name: req.Name, Description: req.Description}

	err := a.repo.CreateTeam(ctx, &team)

	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncCreateSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(team.Id, kafka.Create))
	if err != nil {
		log.Error().Err(err)
	}

	log.Debug().Msgf("new team was created successfully with id=%d", team.Id)

	return &desc.CreateTeamV1Response{Id: team.Id}, nil
}

// MultiCreateTeamV1 is the method that handles creating multiple teams.
func (a *api) MultiCreateTeamV1(
	ctx context.Context,
	req *desc.MultiCreateTeamV1Request) (*desc.MultiCreateTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("MultiCreateTeamV1() was called with len=%d", len(req.Teams))

	tracer := opentracing.GlobalTracer()
	parentSpan := tracer.StartSpan("MultiCreateTeamV1")
	defer parentSpan.Finish()

	teams := make([]models.Team, 0, len(req.Teams))
	for _, team := range req.Teams {
		teams = append(teams, models.Team{
			Name:        team.Name,
			Description: team.Description,
		})
	}

	batches := utils.SplitToBulks(teams, config.GetInstance().Common.BatchSize)

	var teamsIds []uint64

	for i, batch := range batches {
		ids, err := a.repo.CreateTeams(ctx, batch)

		if err != nil {
			return &desc.MultiCreateTeamV1Response{Ids: teamsIds}, status.Error(codes.Internal, err.Error())
		}

		childSpan := tracer.StartSpan(
			fmt.Sprintf("batch_index=%d, batch_size=%d", i, len(batch)),
			opentracing.ChildOf(parentSpan.Context()),
		)
		childSpan.Finish()

		teamsIds = append(teamsIds, ids...)
	}

	return &desc.MultiCreateTeamV1Response{
		Ids: teamsIds,
	}, nil
}

// GetTeamV1 is the method that handles fetching requested team.
func (a *api) GetTeamV1(
	ctx context.Context,
	req *desc.GetTeamV1Request) (*desc.GetTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("GetTeamV1() was called (id=%d)", req.Id)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("GetTeamV1")
	defer span.Finish()

	team, err := a.repo.GetTeam(ctx, req.Id)

	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &desc.GetTeamV1Response{Team: converter.TeamToDTO(team)}

	return response, nil
}

// ListTeamsV1 is the method that handles fetching multiple teams using pagination settings.
func (a *api) ListTeamsV1(
	ctx context.Context,
	req *desc.ListTeamsV1Request) (*desc.ListTeamsV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("ListTeamsV1() was called (limit=%d, offset=%d)", req.Limit, req.Offset)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("ListTeamsV1")
	defer span.Finish()

	teams, total, err := a.repo.ListTeams(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseTeams := make([]*desc.Team, 0, len(teams))
	for _, team := range teams {
		responseTeams = append(responseTeams, converter.TeamToDTO(&team))
	}

	return &desc.ListTeamsV1Response{Total: total, Teams: responseTeams}, nil
}

// RemoveTeamV1 is the method that handles removing team by id if exists.
func (a *api) RemoveTeamV1(
	ctx context.Context,
	req *desc.RemoveTeamV1Request) (*desc.RemoveTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("RemoveTeamV1() was called (id=%d)", req.Id)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("RemoveTeamV1")
	defer span.Finish()

	err := a.repo.RemoveTeam(ctx, req.Id)
	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncDeleteSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(req.Id, kafka.Delete))
	if err != nil {
		log.Error().Err(err)
	}

	return &desc.RemoveTeamV1Response{}, nil
}

// UpdateTeamV1 is the method that handles updating corresponding team.
func (a *api) UpdateTeamV1(
	ctx context.Context,
	req *desc.UpdateTeamV1Request) (*desc.UpdateTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msgf("UpdateTeamV1() was called (id=%d)", req.Team.Id)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("UpdateTeamV1")
	defer span.Finish()

	team := converter.TeamFromDTO(req.Team)

	err := a.repo.UpdateTeam(ctx, team)

	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncUpdateSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(team.Id, kafka.Update))
	if err != nil {
		log.Error().Err(err)
	}

	return &desc.UpdateTeamV1Response{}, nil
}

// SearchTeamsV1 is the method that handles teams searching.
func (a *api) SearchTeamsV1(
	ctx context.Context,
	req *desc.SearchTeamV1Request) (*desc.SearchTeamV1Response, error) {
	metrics.IncTotalRequestsCounter()
	if err := req.Validate(); err != nil {
		metrics.IncInvalidRequestsCounter()
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Debug().Msg("SearchTeamsV1() was called")

	teams, err := a.repo.SearchTeams(ctx, req.Query, utils.SearchType(req.Type))
	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseTeams := make([]*desc.Team, 0, len(teams))
	for _, team := range teams {
		responseTeams = append(responseTeams, &desc.Team{
			Id:          team.Id,
			Name:        team.Name,
			Description: team.Description,
		})
	}

	return &desc.SearchTeamV1Response{Teams: responseTeams}, nil
}
