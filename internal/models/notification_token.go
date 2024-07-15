package models

import (
	"github.com/jinzhu/gorm"
)

type NotificationToken struct {
	gorm.Model
	UserID  string `gorm:"not null"`
	Token  string `gorm:"not null"`
}