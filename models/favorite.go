package models

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint
	AnnonceID uint
}
