package models

import (
	"github.com/jinzhu/gorm"
)

type Annonce struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(250)"`
	Description *string `gorm:"type:varchar(250)"`
	UserID      string
	CatID       string
}
