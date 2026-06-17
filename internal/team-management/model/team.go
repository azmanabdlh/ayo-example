package model

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	LogoURL     string
	FoundedYear int
	Address     string
	City        string

	Player []Player `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
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
