package models

import "github.com/jinzhu/gorm"

type ReportedAnnonce struct {
	gorm.Model
	AnnonceID      uint
	ReporterUserID string
	ReportedUserID string
	ReasonID       uint
	IsHandled      bool `gorm:"default:false"`
}
