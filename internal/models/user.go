package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique_index;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Adress   *string
	Roles    []Roles `gorm:"many2many:user_languages;"`
	Annonce  []Annonce
	Favorite []Annonce
}

type Annonce struct {
	gorm.Model
	UserID uint
	Cats   []Cats `gorm:"foreignKey:CustomReferer"`
}

type Cats struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null"`
	Color     string `gorm:"type:varchar(100);not null"`
	BirthDate *time.Time
}

type Roles struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
}
