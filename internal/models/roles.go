package models

import "github.com/jinzhu/gorm"

type RoleName string

const (
	Admin    RoleName = "admin"
	UserRole RoleName = "user"
)

type Roles struct {
	gorm.Model
	Name RoleName `gorm:"type:varchar(100);not null"`
}
