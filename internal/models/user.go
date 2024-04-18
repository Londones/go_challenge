package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);not null"`
	Email         string `gorm:"type:varchar(100);unique_index;not null"`
	Password      string `gorm:"type:varchar(100);not null"`
	AddressRue    string `gorm:"type:varchar(250)"`
	Cp            string `gorm:"type:char(5)"`
	Ville         string `gorm:"type:varchar(100)"`
	AssociationID uint
	Annonce       []Annonce
	Favorite      []Annonce
	Rating        []Rating
	Roles         []Roles `gorm:"many2many:user_role;"`
}
