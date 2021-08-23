package api

import (
	"context"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	log "github.com/rs/zerolog/log"
)

type api struct {
	desc.UnimplementedOcpTeamApiServer
}

func NewOcpTeamApi() desc.OcpTeamApiServer {
	return &api{}
}

func (a *api) CreateTeamV1(
	ctx context.Context,
	req *desc.CreateTeamV1Request) (*desc.CreateTeamV1Response, error) {
	log.Printf("Create team (name=%s, description=%s)", req.Name, req.Description)
	return &desc.CreateTeamV1Response{}, nil
}

func (a *api) GetTeamV1(
	ctx context.Context,
	req *desc.GetTeamV1Request) (*desc.GetTeamV1Response, error) {
	log.Printf("Get team (id=%d)", req.Id)
	return &desc.GetTeamV1Response{}, nil
}

func (a *api) ListTeamsV1(
	ctx context.Context,
	req *desc.ListTeamsV1Request) (*desc.ListTeamsV1Response, error) {
	log.Printf("List teams (limit=%d, offset=%d)", req.Limit, req.Offset)
	return &desc.ListTeamsV1Response{}, nil
}

func (a *api) RemoveTeamV1(
	ctx context.Context,
	req *desc.RemoveTeamV1Request) (*desc.RemoveTeamV1Response, error) {
	log.Printf("Remove team (id=%d)", req.Id)
	return &desc.RemoveTeamV1Response{}, nil
}