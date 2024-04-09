package models

import (
	"github.com/jinzhu/gorm"
)

type Annonce struct {
	gorm.Model
	UserID      uint
	Description *string `gorm:"type:varchar(250)"`
	Cats        []Cats  `gorm:"foreignKey:CustomReferer"`
}
