package model

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Email           string `gorm:"unique;not null" json:"email" form:"email"`
	Password        string `gorm:"not null" json:"password" form:"password"`
	Fullname        string `gorm:"not null" json:"fullname" form:"fullname"`
	Displayname     string `json:"displayname" form:"displayname" gorm:"type:varchar(255)not null"`
	Alumnus         string `json:"alumnus" form:"alumnus" gorm:"type:varchar(255)not null"`
	Work     string `json:"work" form:"work" gorm:"type:varchar(255)not null"`
	PracticeAddress string `json:"practice_address" form:"practice_address" gorm:"type:text not null"`
	Price           float64
	Balance         float64
	Photo 		 string `json:"photo" form:"photo"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	Status 	   string `json:"status" form:"status" gorm:"type:varchar(11)"`
	Specialist string `json:"specialist" form:"specialist" gorm:"type:varchar(255)"`
	Description string `json:"description" form:"description" gorm:"type:text"`
}