package models

import (
	"github.com/jinzhu/gorm"
)

type Chat struct {
	gorm.Model
	User1ID   string `gorm:"not null"`
	User2ID   string `gorm:"not null"`
	AnnonceID uint   `gorm:"not null"`
}
