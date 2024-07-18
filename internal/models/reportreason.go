package models

import "github.com/jinzhu/gorm"

type ReportReason struct {
	gorm.Model
	Reason string `gorm:"type:varchar(100);not null" json:"reason"`
}
