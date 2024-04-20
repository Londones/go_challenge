package models

import "github.com/jinzhu/gorm"

type Association struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100)"`
	AddressRue string `gorm:"type:varchar(250)"`
	Cp         string `gorm:"type:char(5)"`
	Ville      string `gorm:"type:varchar(100)"`
	Phone      string `gorm:"type:varchar(13)"`
	Email      string `gorm:"type:varchar(100)"`
	//ValidationDocument1 byte Je ne sais pas encore comment les mettre dans les models.
	//ValidationDocument1 byte
	//ValidationDocument1 byte
	Members  []User
	Verified bool
}
