package models

import "github.com/jinzhu/gorm"

type ReportedAnnonce struct {
	gorm.Model
	AnnonceID      uint   `gorm:"not null" json:"annonceId"`
	ReporterUserID string `gorm:"not null" json:"reporterUserId"`
	ReportedUserID string `gorm:"not null" json:"reportedUserId"`
	ReasonID       uint   `gorm:"not null" json:"reasonId"`
	IsHandled      bool   `gorm:"default:false"`
}
