package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);not null"`
	Email         string `gorm:"type:varchar(100);unique_index;not null"`
	Password      string `gorm:"type:varchar(100);not null"`
	Address       *string
	AssociationID uint
	Annonce       []Annonce
	Favorite      []Annonce
	Rating        []Rating
	Roles         []Roles `gorm:"many2many:user_role;"`
}
