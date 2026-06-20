package service

import (
	"context"
	"errors"
	"maps"
	"slices"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
	teamManagement "github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LeagueManagement struct {
	db *gorm.DB
}

func New(db *gorm.DB) *LeagueManagement {
	l := new(LeagueManagement)
	l.db = db

	// auto migration
	model.New(db)

	return l
}

func (l *LeagueManagement) AssignMatchPlayerLineup(ctx context.Context, matchID int64, req dtos.MatchPlayerParam) error {
	const maxPlayerLineup = 20

	if len(req.Lineup) > maxPlayerLineup {
		return errors.New("match player lineup to much")
	}

	match, err := gorm.G[model.Match](l.db).
		Where("id = ? ", matchID).
		First(ctx)
	if err != nil {
		return err
	}

	// validate invariant case ===
	// TODO:
	// 1. assign match lineups only in the upcoming match phase
	// 2. valid teamID is playing on the match
	// 3. player is member of teamID

	if !match.IsUpcoming() {
		return errors.New("invalid match phase")
	}

	if !match.ValidTeam(req.TeamID) {
		return errors.New("invalid teamID")
	}

	playerLineup := make(map[int64]dtos.PlayerLineup)
	for _, lineup := range req.Lineup {
		playerLineup[lineup.PlayerID] = lineup
	}

	playerIDs := slices.Collect(
		maps.Keys(playerLineup),
	)

	players, err := gorm.G[teamManagement.Player](l.db).
		Where("id IN ?", playerIDs).
		Where("team_id = ?", req.TeamID).
		Find(ctx)
	if err != nil {
		return err
	}

	if len(players) != len(playerLineup) {
		return errors.New("invalid player team")
	}

	// assign match players
	// TODO: sync player lineups using upsert (insert new or update if already exists)

	return l.db.Transaction(func(tx *gorm.DB) error {
		for _, player := range players {
			lineup, _ := playerLineup[player.ID]

			err = l.db.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "match_id"},
					{Name: "team_id"},
					{Name: "player_id"},
				},

				DoUpdates: clause.AssignmentColumns([]string{
					"team_id",
					"position_slot",
					"is_starter",
				}),
			}).Create(&model.MatchPlayerLineup{
				MatchID: matchID,
				TeamID:  req.TeamID,

				PositionSlot: lineup.PositionSlot,
				IsStarter:    lineup.IsStarter,
				PlayerID:     player.ID,
			}).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (l *LeagueManagement) Finish(ctx context.Context, matchID int64) error {
	match, err := gorm.G[model.Match](l.db).
		Where("id = ?", matchID).
		First(ctx)
	if err != nil {
		return err
	}

	if match.IsCancelled() || match.IsUpcoming() || match.IsFinished() {
		return errors.New("invalid match phase")
	}

	// IsOngoing match to set finish
	return l.db.Model(&match).
		Update("phase", 3).
		Error
}

func (l *LeagueManagement) CreateMatch(ctx context.Context, matchParam dtos.MatchParam) error {

	err := gorm.G[model.Match](l.db).
		Create(ctx, &model.Match{
			Title:      matchParam.Title,
			MatchDate:  matchParam.MatchDate,
			HomeTeamID: matchParam.AwayTeamID,
			AwayTeamID: matchParam.AwayTeamID,
			VenueID:    matchParam.VenueID,
			Phase:      1, // active
		})
	if err != nil {
		return err
	}

	return nil
}

func (l *LeagueManagement) FindMatchHighlight(ctx context.Context, matchID int64) (model.MatchHighlight, error) {
	var (
		match     model.Match
		highlight model.MatchHighlight
	)

	err := l.db.
		Preload("PlayerLineup").
		Preload("Venue").
		Preload("HomeTeam").
		Preload("HomeTeam.Player").
		Preload("AwayTeam").
		Preload("AwayTeam.Player").
		Preload("Goal").
		First(&match, matchID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return highlight, nil
	}

	if err != nil {
		return highlight, err
	}

	res := match.ToMatchHighlight()
	if res == nil {

		return highlight, errors.New("invalid data")
	}

	highlight = *res
	return highlight, nil
}

func (l *LeagueManagement) AddRecordGoal(ctx context.Context, matchID int64, req dtos.GoalParam) error {
	var (
		match model.Match
	)

	err := l.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("PlayerLineup").
		First(&match, matchID).
		Error

	if err != nil {
		return err
	}

	// need to validate invariant match player
	if !match.CanRecordGoal(req.PlayerID, req.TeamID) {
		return errors.New("invalid player to this match")
	}

	// add record goal player
	err = gorm.G[model.Goal](l.db).
		Create(ctx, &model.Goal{
			PlayerID: req.PlayerID,
			MatchID:  matchID,
			TeamID:   req.TeamID,

			ScoredAtMinute: req.ScoredAtMinute,
		})
	if err != nil {
		return err
	}

	return nil
}
