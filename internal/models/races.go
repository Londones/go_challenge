package models

import (
	"github.com/jinzhu/gorm"
)

type Races struct {
	gorm.Model
	RaceName string `gorm:"type:varchar(100);not null"`
	Cats     []Cats
}
