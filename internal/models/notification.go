package models

import "github.com/jinzhu/gorm"

type Notification struct {
	gorm.Model
	Token  string
	Title  string
	Text   string
	RoomID string
}
