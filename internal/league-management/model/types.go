package model

import (
	"time"

	"gorm.io/gorm"
)

type Venue struct {
	gorm.Model

	ID uint `gorm:"primaryKey"`

	Name          string
	Address       string
	City          string
	GoogleMapsURL string
	Capacity      int

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func New(db *gorm.DB) {
	db.AutoMigrate(
		&Venue{},
		&Match{},
		&Goal{},
		&MatchPlayerLineup{},
	)
}

var slotCoordinates = map[string]struct {
	X int
	Y int
}{
	"GK": {0, 50},

	"LB":   {20, 10},
	"CB-L": {20, 35},
	"CB-R": {20, 65},
	"RB":   {20, 90},

	"CM-L": {50, 35},
	"CM-R": {50, 65},

	"LW": {75, 20},
	"RW": {75, 80},

	"ST-L": {90, 40},
	"ST-R": {90, 60},
}
