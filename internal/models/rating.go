package models

import "github.com/jinzhu/gorm"

type Rating struct {
	gorm.Model
	Mark      int8
	UserID    string
	AnnonceID uint
}
