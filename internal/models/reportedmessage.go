package models

import "github.com/jinzhu/gorm"

type ReportedMessage struct {
	gorm.Model
	MessageID      uint   `gorm:"not null" json:"messageId"`
	ReporterUserID string `gorm:"not null" json:"reporterUserId"`
	ReportedUserID string `gorm:"not null" json:"reportedUserId"`
	ReasonID       uint   `gorm:"not null" json:"reasonId"`
	IsHandled      bool   `gorm:"default:false"`
}
