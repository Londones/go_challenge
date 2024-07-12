package models

import "github.com/jinzhu/gorm"

type Rating struct {
	gorm.Model
	Mark     int8
	Comment  string
	UserID   string
	AuthorID string
}
