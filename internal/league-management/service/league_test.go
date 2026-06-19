package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	leagueManagementSvc "github.com/azmanabdlh/ayo-example/internal/league-management/service"
	"github.com/azmanabdlh/ayo-example/internal/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLeagueManagement_SubstitutePlayer(t *testing.T) {
	tests := []struct {
		name    string
		matchID int64
		req     dtos.SubstitutePlayerParam
		mockFn  func(sqlmock.Sqlmock, dtos.SubstitutePlayerParam)
		wantErr bool
	}{
		{
			name:    "success",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				playersRows := sqlmock.NewRows([]string{"id", "team_id"}).
					AddRow(req.PlayerID, req.TeamID).
					AddRow(req.SubstitutedForPlayerID, req.TeamID)

				mock.ExpectQuery(`SELECT .*players.*`).WillReturnRows(playersRows)

				teamRows := sqlmock.NewRows([]string{"id"}).AddRow(req.TeamID)
				mock.ExpectQuery(`SELECT .*FROM "teams"`).WillReturnRows(teamRows)

				lineupRows := sqlmock.NewRows([]string{"id", "match_id", "team_id", "player_id", "position_slot"}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK")
				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				insertRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO "match_players"`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						int64(1),
						req.TeamID,
						"GK",
						true,
						req.Minute,
						int(0),
						req.SubstitutedForPlayerID,
						int64(0),
						req.Reason,
					).
					WillReturnRows(insertRows)

				mock.ExpectExec(`UPDATE .*match_players`).
					WithArgs(
						false,
						req.Minute,
						req.Reason,
						req.PlayerID,
						sqlmock.AnyArg(),
						int64(100),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "invalid player data",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				playersRows := sqlmock.NewRows([]string{"id", "team_id"}).
					AddRow(req.PlayerID, req.TeamID).
					AddRow(req.SubstitutedForPlayerID, req.TeamID+1)

				mock.ExpectQuery(`SELECT .*players.*`).WillReturnRows(playersRows)
			},
			wantErr: true,
		},
		{
			name:    "players query error",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				mock.ExpectQuery(`SELECT .*players.*`).WillReturnError(errors.New("db fail"))
			},
			wantErr: true,
		},
		{
			name:    "match player not found",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				playersRows := sqlmock.NewRows([]string{"id", "team_id"}).
					AddRow(req.PlayerID, req.TeamID).
					AddRow(req.SubstitutedForPlayerID, req.TeamID)
				mock.ExpectQuery(`SELECT .*players.*`).WillReturnRows(playersRows)

				teamRows := sqlmock.NewRows([]string{"id"}).AddRow(req.TeamID)
				mock.ExpectQuery(`SELECT .*FROM "teams"`).WillReturnRows(teamRows)

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:    "insert fail",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				playersRows := sqlmock.NewRows([]string{"id", "team_id"}).
					AddRow(req.PlayerID, req.TeamID).
					AddRow(req.SubstitutedForPlayerID, req.TeamID)
				mock.ExpectQuery(`SELECT .*players.*`).WillReturnRows(playersRows)

				teamRows := sqlmock.NewRows([]string{"id"}).AddRow(req.TeamID)
				mock.ExpectQuery(`SELECT .*FROM "teams"`).WillReturnRows(teamRows)

				lineupRows := sqlmock.NewRows([]string{"id", "match_id", "team_id", "player_id", "position_slot"}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK")
				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "match_players"`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						int64(1),
						req.TeamID,
						"GK",
						true,
						req.Minute,
						int(0),
						req.SubstitutedForPlayerID,
						int64(0),
						req.Reason,
					).
					WillReturnError(errors.New("insert fail"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "update fail",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				playersRows := sqlmock.NewRows([]string{"id", "team_id"}).
					AddRow(req.PlayerID, req.TeamID).
					AddRow(req.SubstitutedForPlayerID, req.TeamID)
				mock.ExpectQuery(`SELECT .*players.*`).WillReturnRows(playersRows)

				teamRows := sqlmock.NewRows([]string{"id"}).AddRow(req.TeamID)
				mock.ExpectQuery(`SELECT .*FROM "teams"`).WillReturnRows(teamRows)

				lineupRows := sqlmock.NewRows([]string{"id", "match_id", "team_id", "player_id", "position_slot"}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK")
				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				insertRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO "match_players"`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						int64(1),
						req.TeamID,
						"GK",
						true,
						req.Minute,
						int(0),
						req.SubstitutedForPlayerID,
						int64(0),
						req.Reason,
					).
					WillReturnRows(insertRows)

				mock.ExpectExec(`UPDATE .*match_players`).
					WithArgs(
						false,
						req.Minute,
						req.Reason,
						req.PlayerID,
						sqlmock.AnyArg(),
						int64(100),
					).
					WillReturnError(errors.New("update fail"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mocks.NewMockDB(t)
			tt.mockFn(mock, tt.req)

			svc := leagueManagementSvc.New(db)

			err := svc.SubstitutePlayer(context.Background(), tt.matchID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
