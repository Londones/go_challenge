package models

import "github.com/jinzhu/gorm"

type Rating struct {
	gorm.Model
	UserID    uint
	AnnonceID uint
}
