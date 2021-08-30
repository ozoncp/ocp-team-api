package api

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
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

const (
	batchSize = 2
)

type api struct {
	desc.UnimplementedOcpTeamApiServer
	repo     repo.Repo
	producer kafka.Producer
}

func NewOcpTeamApi(repo repo.Repo, producer kafka.Producer) desc.OcpTeamApiServer {
	return &api{
		repo:     repo,
		producer: producer,
	}
}

func (a *api) CreateTeamV1(
	ctx context.Context,
	req *desc.CreateTeamV1Request) (*desc.CreateTeamV1Response, error) {
	log.Printf("Create team (name=%s, description=%s)", req.Name, req.Description)
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("CreateTeamV1")
	defer span.Finish()

	team := models.Team{Name: req.Name, Description: req.Description}

	err := a.repo.CreateTeam(ctx, &team)

	if err != nil {
		log.Info().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncCreateSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(team.Id, kafka.Create))
	if err != nil {
		log.Info().Err(err)
	}

	log.Info().Msgf("new team was created successfully with id=%d", team.Id)

	return &desc.CreateTeamV1Response{Id: team.Id}, nil
}

func (a *api) MultiCreateTeamV1(
	ctx context.Context,
	req *desc.MultiCreateTeamV1Request) (*desc.MultiCreateTeamV1Response, error) {
	log.Printf("Multi create team")
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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

	batches := utils.SplitToBulks(teams, batchSize)

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

func (a *api) GetTeamV1(
	ctx context.Context,
	req *desc.GetTeamV1Request) (*desc.GetTeamV1Response, error) {
	log.Printf("Get team (id=%d)", req.Id)
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("GetTeamV1")
	defer span.Finish()

	team, err := a.repo.GetTeam(ctx, req.Id)

	if err != nil {
		log.Info().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &desc.GetTeamV1Response{
		Team: &desc.Team{
			Id:          team.Id,
			Name:        team.Name,
			Description: team.Description,
		},
	}

	return response, nil
}

func (a *api) ListTeamsV1(
	ctx context.Context,
	req *desc.ListTeamsV1Request) (*desc.ListTeamsV1Response, error) {
	log.Printf("List teams (limit=%d, offset=%d)", req.Limit, req.Offset)
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("ListTeamsV1")
	defer span.Finish()

	teams, total, err := a.repo.ListTeams(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Info().Err(err)
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

	return &desc.ListTeamsV1Response{Total: total, Teams: responseTeams}, nil
}

func (a *api) RemoveTeamV1(
	ctx context.Context,
	req *desc.RemoveTeamV1Request) (*desc.RemoveTeamV1Response, error) {
	log.Printf("Remove team (id=%d)", req.Id)
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("RemoveTeamV1")
	defer span.Finish()

	err := a.repo.RemoveTeam(ctx, req.Id)
	if err != nil {
		log.Info().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncDeleteSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(req.Id, kafka.Delete))
	if err != nil {
		log.Info().Err(err)
	}

	return &desc.RemoveTeamV1Response{}, nil
}

func (a *api) UpdateTeamV1(
	ctx context.Context,
	req *desc.UpdateTeamV1Request) (*desc.UpdateTeamV1Response, error) {
	log.Printf("Update team (id=%d)", req.Team.Id)
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("UpdateTeamV1")
	defer span.Finish()

	team := models.Team{
		Id:          req.Team.Id,
		Name:        req.Team.Name,
		Description: req.Team.Description,
	}

	err := a.repo.UpdateTeam(ctx, team)

	if err != nil {
		log.Info().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics.IncUpdateSuccessCounter()
	err = a.producer.Send(kafka.NewMessage(team.Id, kafka.Update))
	if err != nil {
		log.Info().Err(err)
	}

	return &desc.UpdateTeamV1Response{}, nil
}
