package models

import "github.com/jinzhu/gorm"

type Roles struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
}
