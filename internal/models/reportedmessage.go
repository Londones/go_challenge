package models

import "github.com/jinzhu/gorm"

type ReportedMessage struct {
	gorm.Model
	MessageID      uint
	ReporterUserID string
	ReportedUserID string
	ReasonID       uint
	IsHandled      bool `gorm:"default:false"`
}
