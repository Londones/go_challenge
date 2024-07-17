package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	RoomID   uint   `gorm:"not null"`
	SenderID string `gorm:"not null"`
	Content  string `gorm:"not null"`
	IsRead   bool   `gorm:"default:false"`
}
