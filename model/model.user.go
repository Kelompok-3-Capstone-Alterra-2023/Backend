package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string   `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Username       string   `json:"username" form:"username" gorm:"type:varchar(255)unique;not null"`
	Password       string   `json:"password" form:"password" gorm:"not null"`
	Telp           string   `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat         string   `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender         string   `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	BirthDate      string   `json:"birthdate" form:"birthdate" gorm:"type:date"`
	Status_Online  bool     `json:"status_online" form:"status_online" gorm:"type:boolean"`
	Doctors        []Doctor `gorm:"many2many:user_doctors;"`
	ChatwithDoctor []Doctor `gorm:"many2many:Chatrooms"`
}
