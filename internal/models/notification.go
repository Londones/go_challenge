package models

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	UserID  string `gorm:"not null"`
	Content string `gorm:"not null"`
	IsRead  bool   `gorm:"not null"`
}
