package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID           int64  `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func New(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
