package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Cats struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null"`
	Color     string `gorm:"type:varchar(100);not null"`
	BirthDate *time.Time
}
