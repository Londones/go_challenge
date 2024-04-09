package models

import (
	"github.com/jinzhu/gorm"
)

type Annonce struct {
	gorm.Model
	UserID uint
	Cats   []Cats `gorm:"foreignKey:CustomReferer"`
}
