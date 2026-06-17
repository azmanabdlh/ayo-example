package model

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model

	ID         int64 `gorm:"primaryKey"`
	Name       string
	Height     int
	Weight     int
	Position   string
	BackNumber int

	TeamID int64

	// relations
	Team         Team           `gorm:"foreignKey:TeamID"`
	PlayerMember []PlayerMember `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Player) IsMemberOf(teamID int64) bool {
	return p.Team.ID == teamID
}

type PlayerMember struct {
	gorm.Model

	ID               int64 `gorm:"primaryKey"`
	PlayerID         int64
	PlayerBackNumber int
	TeamID           int64
	JoinedAt         time.Time
	LeftAt           time.Time

	Player Player `gorm:"foreignKey:PlayerID"`
	Team   Team   `gorm:"foreignKey:TeamID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
