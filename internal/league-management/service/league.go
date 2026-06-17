package service

import (
	"context"
	"errors"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
	teamManagement "github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"gorm.io/gorm"
)

func (l *LeagueManagement) SubstitutePlayer(ctx context.Context, matchID int64, req dtos.SubstitutePlayerParam) error {

	players, err := gorm.G[teamManagement.Player](l.db).
		Preload("Team", nil).
		Where("id IN ? ", []int64{
			req.PlayerID,
			req.SubstitutedForPlayerID,
		}).
		Find(ctx)

	if err != nil {
		return err
	}

	for _, player := range players {
		if !player.IsMemberOf(req.TeamID) {

			logger.Info("invalid player data. player: %v", players)
			return errors.New("invalid player")
		}
	}

	return l.db.Transaction(func(tx *gorm.DB) error {

		// current player lineup
		matchPlayer, err := gorm.G[model.MatchPlayerLineup](l.db).
			Where("player_id = ? ", req.PlayerID).
			Where("team_id = ? ", req.TeamID).
			Where("match_id = ? ", matchID).
			First(ctx)
		if err != nil {
			return err
		}

		// in
		err = gorm.G[model.MatchPlayerLineup](l.db).
			Create(ctx, &model.MatchPlayerLineup{
				MatchID:   matchID,
				TeamID:    req.TeamID,
				IsStarter: true,

				PositionSlot: matchPlayer.PositionSlot,

				MinuteIn: req.Minute,
				PlayerID: req.SubstitutedForPlayerID,
				Reason:   req.Reason,
			})
		if err != nil {
			return err
		}

		// out update
		matchPlayer.IsStarter = false
		matchPlayer.MinuteOut = req.Minute
		matchPlayer.Reason = req.Reason
		matchPlayer.SubstitutedForPlayerID = req.PlayerID

		err = l.db.
			Save(&matchPlayer).
			Error

		return err
	})

}
