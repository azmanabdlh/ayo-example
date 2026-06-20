package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	leagueManagementSvc "github.com/azmanabdlh/ayo-example/internal/league-management/service"
	"github.com/azmanabdlh/ayo-example/internal/mocks"
	"github.com/stretchr/testify/assert"
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
				lineupRows := sqlmock.NewRows([]string{
					"id",
					"match_id",
					"team_id",
					"player_id",
					"position_slot",
					"is_starter",
					"minute_in",
					"minute_out",
					"substituted_for_player_id",
					"reason",
					"created_at",
					"updated_at",
					"deleted_at",
				}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK", true, 0, 0, nil, "", time.Now(), time.Now(), nil).
					AddRow(int64(101), int64(1), req.TeamID, req.SubstitutedForPlayerID, "GK", false, 0, 0, nil, "", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*match_players`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE .*match_players`).WillReturnResult(sqlmock.NewResult(1, 1))
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
				lineupRows := sqlmock.NewRows([]string{
					"id",
					"match_id",
					"team_id",
					"player_id",
					"position_slot",
					"is_starter",
					"minute_in",
					"minute_out",
					"substituted_for_player_id",
					"reason",
					"created_at",
					"updated_at",
					"deleted_at",
				}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK", true, 0, 0, nil, "", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)
			},
			wantErr: true,
		},
		{
			name:    "match players query error",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnError(errors.New("db fail"))
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
				lineupRows := sqlmock.NewRows([]string{
					"id",
					"match_id",
					"team_id",
					"player_id",
					"position_slot",
					"is_starter",
					"minute_in",
					"minute_out",
					"substituted_for_player_id",
					"reason",
					"created_at",
					"updated_at",
					"deleted_at",
				})

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)
			},
			wantErr: true,
		},
		{
			name:    "out update fail",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				lineupRows := sqlmock.NewRows([]string{
					"id",
					"match_id",
					"team_id",
					"player_id",
					"position_slot",
					"is_starter",
					"minute_in",
					"minute_out",
					"substituted_for_player_id",
					"reason",
					"created_at",
					"updated_at",
					"deleted_at",
				}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK", true, 0, 0, nil, "", time.Now(), time.Now(), nil).
					AddRow(int64(101), int64(1), req.TeamID, req.SubstitutedForPlayerID, "GK", false, 0, 0, nil, "", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*match_players`).WillReturnError(errors.New("update fail"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "in update fail",
			matchID: 1,
			req: dtos.SubstitutePlayerParam{
				TeamID:                 10,
				PlayerID:               1,
				SubstitutedForPlayerID: 2,
				Minute:                 60,
				Reason:                 "injury",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.SubstitutePlayerParam) {
				lineupRows := sqlmock.NewRows([]string{
					"id",
					"match_id",
					"team_id",
					"player_id",
					"position_slot",
					"is_starter",
					"minute_in",
					"minute_out",
					"substituted_for_player_id",
					"reason",
					"created_at",
					"updated_at",
					"deleted_at",
				}).
					AddRow(int64(100), int64(1), req.TeamID, req.PlayerID, "GK", true, 0, 0, nil, "", time.Now(), time.Now(), nil).
					AddRow(int64(101), int64(1), req.TeamID, req.SubstitutedForPlayerID, "GK", false, 0, 0, nil, "", time.Now(), time.Now(), nil)

				mock.ExpectQuery(`SELECT .*match_players.*`).WillReturnRows(lineupRows)

				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*match_players`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE .*match_players`).WillReturnError(errors.New("update fail"))
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
