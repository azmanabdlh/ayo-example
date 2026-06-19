package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	authService "github.com/azmanabdlh/ayo-example/internal/authentication/service"
	"github.com/azmanabdlh/ayo-example/internal/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		email     string
		password  string
		mockFn    func(sqlmock.Sqlmock, *authService.MockProvider, int64, string)
		wantToken string
		wantErr   bool
	}{
		{
			name:     "success",
			email:    "admin@test.com",
			password: "secret",
			mockFn: func(mock sqlmock.Sqlmock, provider *authService.MockProvider, userID int64, passwordHash string) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).
					AddRow(userID, "admin@test.com", passwordHash, now, now, nil)

				mock.ExpectQuery(`SELECT .*FROM "users"`).
					WillReturnRows(rows)

				provider.On("GenerateToken", userID).
					Return("jwt-token", nil).
					Once()
			},
			wantToken: "jwt-token",
		},
		{
			name:     "user not found",
			email:    "missing@test.com",
			password: "secret",
			mockFn: func(mock sqlmock.Sqlmock, provider *authService.MockProvider, userID int64, passwordHash string) {
				mock.ExpectQuery(`SELECT .*FROM "users"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:     "wrong password",
			email:    "admin@test.com",
			password: "wrongpass",
			mockFn: func(mock sqlmock.Sqlmock, provider *authService.MockProvider, userID int64, passwordHash string) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).
					AddRow(userID, "admin@test.com", passwordHash, now, now, nil)

				mock.ExpectQuery(`SELECT .*FROM "users"`).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:     "token generation failure",
			email:    "admin@test.com",
			password: "secret",
			mockFn: func(mock sqlmock.Sqlmock, provider *authService.MockProvider, userID int64, passwordHash string) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).
					AddRow(userID, "admin@test.com", passwordHash, now, now, nil)

				mock.ExpectQuery(`SELECT .*FROM "users"`).
					WillReturnRows(rows)

				provider.On("GenerateToken", userID).
					Return("", errors.New("token error")).
					Once()
			},
			wantErr: true,
		},
		{
			name:     "query failure",
			email:    "admin@test.com",
			password: "secret",
			mockFn: func(mock sqlmock.Sqlmock, provider *authService.MockProvider, userID int64, passwordHash string) {
				mock.ExpectQuery(`SELECT .*FROM "users"`).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if hashErr != nil {
		t.Fatal(hashErr)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mocks.NewMockDB(t)
			provider := authService.NewMockProvider(t)
			tt.mockFn(mock, provider, 1, string(passwordHash))

			svc := authService.New(provider, db)

			token, err := svc.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
