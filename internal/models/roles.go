package models

import "github.com/jinzhu/gorm"

type RoleName string

const (
	AdminRole RoleName = "ADMIN"
	UserRole  RoleName = "USER"
	AssoRole  RoleName = "ASSO"
)

type Roles struct {
	gorm.Model
	Name RoleName `gorm:"type:varchar(100);not null"`
}
