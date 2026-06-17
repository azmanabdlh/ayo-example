package service

import (
	"context"
	"errors"

	"github.com/azmanabdlh/ayo-example/internal/team-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"gorm.io/gorm"
)

type TeamManagement struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TeamManagement {
	teamManagement := new(TeamManagement)
	teamManagement.db = db

	model.New(db)

	return teamManagement
}

func (t *TeamManagement) Create(ctx context.Context, req dtos.TeamParam) error {
	team := new(model.Team)

	team.Name = req.Name
	team.FoundedYear = req.FoundedYear
	team.LogoURL = req.LogoURL
	team.City = req.City
	team.Address = req.Address

	return gorm.G[model.Team](t.db).Create(ctx, team)
}

func (t *TeamManagement) Find(ctx context.Context, teamID int64) (team model.Team, err error) {
	// hit SQL
	/**
		SELECT * FROM teams t
			WHERE t.deleted_at IS NULL AND t.id = x

	**/
	team, err = gorm.G[model.Team](t.db).
		Where("id = ?", teamID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return team, nil
	}

	return
}

func (t *TeamManagement) FindAll(ctx context.Context, page, limit int) ([]model.Team, error) {

	// hit SQL
	/**
		SELECT * FROM teams t
			LIMIT x
			OFFSET x
	**/
	teams, err := gorm.G[model.Team](t.db).
		Offset(page).
		Limit(limit).
		Find(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Team{}, nil
	}

	return teams, err
}

func (t *TeamManagement) Modify(ctx context.Context, teamID int64, req dtos.TeamParam) error {

	data := make(map[string]interface{})

	if req.Name != "" {
		data["name"] = req.Name
	}

	if req.LogoURL != "" {
		data["logo_url"] = req.LogoURL
	}

	if req.FoundedYear != 0 {
		data["founded_year"] = req.FoundedYear
	}

	if req.City != "" {
		data["city"] = req.City
	}

	if req.Address != "" {
		data["address"] = req.Address
	}

	rowsAffected, err := gorm.G[map[string]interface{}](t.db).
		Table("teams").
		Where("id = ?", teamID).
		Updates(ctx, data)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (t *TeamManagement) Remove(ctx context.Context, teamID int64) error {
	_, err := gorm.G[model.Team](t.db).Where("id = ?", teamID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
