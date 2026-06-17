package handler

import (
	"context"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
)

type Service interface {
	CreateMatch(ctx context.Context, matchParam dtos.MatchParam) error
	FindMatchHighlight(ctx context.Context, matchID int64) (model.MatchHighlight, error)
	AddRecordGoal(ctx context.Context, matchID int64, goalParam dtos.GoalParam) error
	Finish(ctx context.Context, matchID int64) error

	CreateVenue(ctx context.Context, venueParam dtos.VenueParam) error
	ModifyVenue(ctx context.Context, venueID int64, venueParam dtos.VenueParam) error

	SubstitutePlayer(ctx context.Context, matchID int64, param dtos.SubstitutePlayerParam) error
	AssignMatchPlayerLineup(ctx context.Context, matchID int64, param dtos.MatchPlayerParam) error
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	h := new(Handler)
	h.svc = svc

	return h
}
