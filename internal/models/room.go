package models

import (
	"github.com/jinzhu/gorm"
)

type Room struct {
	gorm.Model
	User1ID   string `gorm:"not null"`
	User2ID   string `gorm:"not null"`
	AnnonceID string `gorm:"not null"`
}
