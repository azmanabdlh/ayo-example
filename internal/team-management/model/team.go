package model

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          int64  `json:"id,omitempty" gorm:"primaryKey"`
	Name        string `json:"name,omitempty" gorm:"not null"`
	LogoURL     string `json:"logo_url,omitempty"`
	FoundedYear int    `json:"founded_year,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`

	Player []Player `json:"players,omitempty" gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (t *Team) HasPlayer() bool {
	return len(t.Player) > 0
}

func (t *Team) HasPlayerID(playerID int64) bool {
	for i := 0; i < len(t.Player); i++ {
		if t.Player[i].ID == playerID {
			return true
		}
	}

	return false
}

func (t *Team) CanAssignPlayer(player Player) bool {
	if t.ID == 0 || player.ID == 0 || player.IsMemberOf(t.ID) {
		return false
	}

	return true
}

func New(db *gorm.DB) {
	db.AutoMigrate(
		&Team{},
		&Player{},
		&PlayerMember{},
	)
}
