package service

import (
	"context"

	"github.com/azmanabdlh/ayo-example/internal/authentication/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Provider interface {
	GenerateToken(userID int64) (string, error)
}

type Service struct {
	provider Provider
	db       *gorm.DB
}

func New(provider Provider, db *gorm.DB) *Service {

	model.New(db)

	return &Service{
		provider: provider,
		db:       db,
	}
}

func (s *Service) Login(ctx context.Context, email, passwordPlain string) (token string, err error) {

	user, err := gorm.G[model.User](s.db).
		Where("email = ?", email).
		First(ctx)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordPlain))
	if err != nil {
		return
	}

	return s.provider.GenerateToken(user.ID)
}

func (s *Service) IsAccountValid(ctx context.Context, userID int64) bool {
	var count int64

	s.db.Model(&model.User{}).
		Select("id").
		Where("id = ?", userID).
		Count(&count)

	// TODO: for complex requirment we need to add handle validation
	// 1. account rule
	// 2. account suspend
	// 3. ....

	return count > 0
}
