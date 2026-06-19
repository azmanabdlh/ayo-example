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
)

func TestLeagueManagement_CreateVenue(t *testing.T) {
	tests := []struct {
		name    string
		req     dtos.VenueParam
		mockFn  func(sqlmock.Sqlmock, dtos.VenueParam)
		wantErr bool
	}{
		{
			name: "success",
			req: dtos.VenueParam{
				Name:          "Pantai Gading",
				City:          "Jakarta",
				Address:       "123",
				Capacity:      123,
				GoogleMapsURL: "ok",
			},
			mockFn: func(mock sqlmock.Sqlmock, req dtos.VenueParam) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO .*venues`).
					WithArgs(
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
						sqlmock.AnyArg(), // deleted_at
						req.Name,
						req.Address,
						req.City,
						req.GoogleMapsURL,
						req.Capacity,
					).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "db error",
			mockFn: func(mock sqlmock.Sqlmock, req dtos.VenueParam) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "venues"`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						req.Name,
						req.Address,
						req.City,
						req.GoogleMapsURL,
						req.Capacity,
					).
					WillReturnError(errors.New("insert fail"))

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

			err := svc.CreateVenue(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLeagueManagement_ModifyVenue(t *testing.T) {
	tests := []struct {
		name    string
		venueID int64
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:    "success",
			venueID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*venues`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "no row affected",
			venueID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*venues`).
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			wantErr: true,
		},
		{
			name:    "db error",
			venueID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE .*venues`).
					WillReturnError(errors.New("update fail"))
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

			err := svc.ModifyVenue(context.Background(), tt.venueID, dtos.VenueParam{
				Name:     "Updated",
				Address:  "Addr",
				City:     "C",
				Capacity: 100,
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
