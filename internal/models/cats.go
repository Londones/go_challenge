package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Cats struct {
	gorm.Model
	Name            string `gorm:"type:varchar(100);not null"`
	BirthDate       *time.Time
	Sexe            string `gorm:"type:varchar(7)"`
	LastVaccine     *time.Time
	LastVaccineName string `gorm:"type:varchar(100)"`
	Color           string `gorm:"type:varchar(100)"`
	Behavior        string `gorm:"type:varchar(100)"`
	Sterilized      bool
	RaceID          string  `gorm:"type:varchar(100);not null"`
	Description     *string `gorm:"type:varchar(250)"`
	Reserved        bool
	AnnonceID       string         `gorm:"type:varchar(100)"`
	PicturesURL     pq.StringArray `gorm:"type:varchar(500)[]"`
}
