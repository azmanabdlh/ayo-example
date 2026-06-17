package service

import (
	"context"
	"errors"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
	"gorm.io/gorm"
)

func (l *LeagueManagement) CreateVenue(ctx context.Context, req dtos.VenueParam) error {

	err := gorm.G[model.Venue](l.db).
		Create(ctx, &model.Venue{
			Name:          req.Name,
			Address:       req.Address,
			City:          req.City,
			GoogleMapsURL: req.GoogleMapsURL,
			Capacity:      req.Capacity,
		})
	if err != nil {
		return err
	}

	return nil
}

func (l *LeagueManagement) ModifyVenue(ctx context.Context, venueID int64, req dtos.VenueParam) error {
	data := make(map[string]interface{})

	if req.Name != "" {
		data["name"] = req.Name
	}

	if req.GoogleMapsURL != "" {
		data["google_maps_url"] = req.GoogleMapsURL
	}

	if req.Capacity != 0 {
		data["address"] = req.Capacity
	}

	if req.City != "" {
		data["city"] = req.City
	}

	if req.Address != "" {
		data["address"] = req.Address
	}

	rowsAffected, err := gorm.G[map[string]interface{}](l.db).
		Table("venues").
		Where("id = ?", venueID).
		Updates(ctx, data)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no row affected")
	}

	return nil
}
