package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	ChatID   uint   `gorm:"not null"`
	SenderID string `gorm:"not null"`
	Content  string `gorm:"not null"`
}
