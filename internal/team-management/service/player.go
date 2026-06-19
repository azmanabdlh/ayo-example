package service

import (
	"context"
	"errors"
	"time"

	"github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"gorm.io/gorm"
)

func (t *TeamManagement) AssignPlayerTeam(ctx context.Context, playerID, teamID int64, newBackNumber int) error {

	player := model.Player{}

	g := t.db.Preload("Team").
		Where("id = ? ", playerID).
		First(&player)

	if g.Error != nil {
		return g.Error
	}

	team, err := t.Find(ctx, teamID)
	if err != nil {
		return err
	}

	if !team.CanAssignPlayer(player) {
		logger.Info("can't assign player_id: %d to team_id: %d", playerID, teamID)
		return nil
	}

	// # add transaction
	// 1. assign player to new team
	// 2. add snapshot player memberships

	return t.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		currentTeamID := player.Team.ID

		rowsAffected, err := gorm.G[model.Player](tx).
			Where("id = ? ", playerID).
			Update(ctx, "team_id", teamID)

		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return errors.New("no row affected")
		}

		_, err = gorm.G[model.PlayerMember](tx).
			Where("player_id = ?", playerID).
			Where("team_id = ?", currentTeamID).
			Update(ctx, "left_at", now)
		if err != nil {
			return err
		}

		err = gorm.G[model.PlayerMember](tx).Create(ctx, &model.PlayerMember{
			PlayerID:         playerID,
			TeamID:           teamID,
			JoinedAt:         now,
			PlayerBackNumber: newBackNumber,
		})
		if err != nil {
			return err
		}

		return nil
	})

}

func (t *TeamManagement) FindPlayer(ctx context.Context, playerID int64) (player model.Player, err error) {

	g := t.db.
		Preload("Team").
		Preload("PlayerMember").
		Preload("PlayerMember.Team").
		Where("id = ? ", playerID).
		First(&player)

	if g.Error != nil {
		err = g.Error
		return
	}

	return
}

func (t *TeamManagement) RemovePlayer(ctx context.Context, playerID int64) error {

	_, err := gorm.G[model.Player](t.db).Where("id = ?", playerID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
