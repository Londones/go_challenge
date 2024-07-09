package models

import (
	"time"
)

// User represents a user in the system.
// swagger:model User
type User struct {
	ID            string `gorm:"type:uuid;primary_key;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Name          string     `gorm:"type:varchar(100);not null"`
	Email         string     `gorm:"type:varchar(100);unique_index;not null"`
	Password      string     `gorm:"type:varchar(100);not null"`
	AddressRue    string     `gorm:"type:varchar(250)"`
	Cp            string     `gorm:"type:char(5)"`
	Ville         string     `gorm:"type:varchar(100)"`
	Associations  []Association `gorm:"many2many:association_members;"`
	Roles         []Roles `gorm:"many2many:user_roles;"`
	GoogleID      string
	ProfilePicURL string `gorm:"type:varchar(500)"`
}
