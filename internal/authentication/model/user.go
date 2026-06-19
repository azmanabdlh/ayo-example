package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID           int64  `json:"id,omitempty" gorm:"primaryKey"`
	Email        string `json:"email,omitempty" gorm:"uniqueIndex"`
	PasswordHash string `json:"password_hash,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func New(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
