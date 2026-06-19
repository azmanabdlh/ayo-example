package model

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID         int64  `json:"id,omitempty" gorm:"primaryKey"`
	Name       string `json:"name,omitempty"`
	Height     int    `json:"height,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	Position   string `json:"position,omitempty"`
	BackNumber int    `json:"back_number,omitempty"`

	TeamID int64 `json:"team_id,omitempty"`

	// relations
	Team         Team           `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	PlayerMember []PlayerMember `json:"player_memberships,omitempty" gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p *Player) IsMemberOf(teamID int64) bool {
	return p.Team.ID == teamID
}

type PlayerMember struct {
	ID               int64     `gorm:"primaryKey"`
	PlayerID         int64     `json:"player_id"`
	PlayerBackNumber int       `json:"player_back_number"`
	TeamID           int64     `json:"team_id"`
	JoinedAt         time.Time `json:"joined_at"`
	LeftAt           time.Time `json:"left_at"`

	Player Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Team   Team   `json:"team,omitempty" gorm:"foreignKey:TeamID"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
