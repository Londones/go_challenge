package models

import "github.com/jinzhu/gorm"

type FeatureFlag struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
	IsEnabled bool `gorm:"type:boolean;not null;default:true"`
}