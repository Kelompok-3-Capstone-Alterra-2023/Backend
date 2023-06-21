package model

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	DoctorIDnoFK int    `gorm:"type:int"`
	UserIDnoFK   int    `gorm:"type:int"`
	Content      string `gorm:"type:text"`
	Sender string `gorm:"type:text"`
}
