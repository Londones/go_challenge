package models

import "github.com/jinzhu/gorm"

type Rating struct {
	gorm.Model
	ID        string `gorm:"type:uuid;primary_key;"`
	Mark      int8
	Comment   string
	UserID    uint
	AnnonceID string
}
