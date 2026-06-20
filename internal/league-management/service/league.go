package service

import (
	"context"
	"errors"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"gorm.io/gorm"
)

func (l *LeagueManagement) SubstitutePlayer(ctx context.Context, matchID int64, req dtos.SubstitutePlayerParam) error {

	matchPlayers, err := gorm.G[model.MatchPlayerLineup](l.db).
		Where("player_id IN ? ", []int64{
			req.PlayerID,
			req.SubstitutedForPlayerID,
		}).
		Where("team_id = ? ", req.TeamID).
		Where("match_id = ? ", matchID).
		Find(ctx)
	if err != nil {
		return err
	}

	matchPlayerID := make(map[int64]model.MatchPlayerLineup)
	for _, m := range matchPlayers {
		matchPlayerID[m.PlayerID] = m
	}

	if len(matchPlayerID) != 2 {
		return errors.New("invalid player")
	}

	if m, ok := matchPlayerID[req.SubstitutedForPlayerID]; ok && m.IsStarter {
		return errors.New("invalid subtitute player")
	}

	return l.db.Transaction(func(tx *gorm.DB) error {

		// TODO
		// 1. update current match player id to is_starter: false
		// 2. then change new match player to is_starter: true

		// out update
		newMatchPlayer := matchPlayerID[req.SubstitutedForPlayerID]

		_, err := gorm.G[model.MatchPlayerLineup](tx).
			Where("team_id = ? ", req.TeamID).
			Where("match_id = ? ", matchID).
			Where("player_id = ? ", req.PlayerID).
			Updates(ctx, model.MatchPlayerLineup{
				IsStarter:              false,
				MinuteOut:              req.Minute,
				Reason:                 req.Reason,
				SubstitutedForPlayerID: &req.SubstitutedForPlayerID,
			})

		if err != nil {
			logger.Info("error update")
			return err
		}

		positionSlot := req.PositionSlot
		if positionSlot == "" {
			positionSlot = newMatchPlayer.PositionSlot
		}

		// in
		_, err = gorm.G[model.MatchPlayerLineup](tx).
			Where("team_id = ? ", req.TeamID).
			Where("match_id = ? ", matchID).
			Where("player_id = ? ", req.SubstitutedForPlayerID).
			Updates(ctx, model.MatchPlayerLineup{
				IsStarter:    true,
				PositionSlot: positionSlot,
				MinuteIn:     req.Minute,
				Reason:       req.Reason,
			})

		return err

	})

}
