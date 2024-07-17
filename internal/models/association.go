package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Association struct {
	gorm.Model
	Name       string         `gorm:"type:varchar(100)"`
	AddressRue string         `gorm:"type:varchar(250)"`
	Cp         string         `gorm:"type:char(5)"`
	Ville      string         `gorm:"type:varchar(100)"`
	Phone      string         `gorm:"type:varchar(13)"`
	Email      string         `gorm:"type:varchar(100)"`
	KbisFile   string         `gorm:"type:varchar(500)"`
	Members    pq.StringArray `gorm:"type:text[]"`
	OwnerID    string         `gorm:"type:uuid;not null"`
	Verified   *bool          `gorm:"type:boolean;default:false"`
}
