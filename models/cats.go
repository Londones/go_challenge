package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Cats struct {
	gorm.Model
	Name            string `gorm:"type:varchar(100);not null"`
	BirthDate       *time.Time
	Sex             string `gorm:"type:varchar(7); not null"`
	LastVaccine     *time.Time
	LastVaccineName string `gorm:"type:varchar(100);not null"`
	Color           string `gorm:"type:varchar(100);not null"`
	Behavior        string `gorm:"type:varchar(100);not null"`
	Sterilized      bool
	Race            string  `gorm:"type:varchar(100);not null"`
	Description     *string `gorm:"type:varchar(250)"`
}
