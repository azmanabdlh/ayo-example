package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	leagueManagementSvc "github.com/azmanabdlh/ayo-example/internal/league-management/service"
	"github.com/azmanabdlh/ayo-example/internal/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLeagueManagement_Finish(t *testing.T) {

	now := time.Now()

	tests := []struct {
		name string

		matchID int64

		mockFn func(sqlmock.Sqlmock)

		wantErr bool
	}{
		{
			name:    "ongoing match become finished",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"phase",
						"match_date",
					},
				).AddRow(
					1,
					1,
					now.Add(-2*time.Hour),
				)

				mock.ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnRows(rows)

				mock.ExpectBegin()
				mock.
					ExpectExec(`UPDATE .*matches.*`).
					WithArgs(3, sqlmock.AnyArg(), 1).
					WillReturnResult(
						sqlmock.NewResult(
							1,
							1,
						),
					)
				mock.ExpectCommit()
			},
		},
		{
			name:    "upcoming match",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"phase",
						"match_date",
					},
				).AddRow(
					1,
					1,
					now.Add(5*time.Hour),
				)

				mock.
					ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnRows(rows)
			},
			wantErr: true,
		},

		{
			name:    "cancelled match no-op",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"phase",
						"match_date",
					},
				).AddRow(
					1,
					2,
					now.Add(-5*time.Hour),
				)

				mock.
					ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:    "already finished",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"phase",
						"match_date",
					},
				).AddRow(
					1,
					3,
					now.Add(5*time.Hour),
				)

				mock.
					ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:    "match not found",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				mock.
					ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnError(
						gorm.ErrRecordNotFound,
					)
			},
			wantErr: true,
		},
		{
			name:    "ongoing update error",
			matchID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"phase",
						"match_date",
					},
				).AddRow(
					1,
					1,
					now.Add(-2*time.Hour),
				)

				mock.
					ExpectQuery(`SELECT .* FROM "matches"`).
					WillReturnRows(rows)

				mock.ExpectBegin()

				mock.
					ExpectExec(`UPDATE "matches"`).
					WillReturnError(
						errors.New("database error"),
					)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := mocks.NewMockDB(t)

			tt.mockFn(mock)

			svc := leagueManagementSvc.New(db)

			err := svc.Finish(
				context.Background(),
				tt.matchID,
			)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(
				t,
				mock.ExpectationsWereMet(),
			)
		})
	}
}
