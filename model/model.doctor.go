package model

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Email           string `gorm:"unique;not null" json:"email" form:"email"`
	Password        string `gorm:"not null" json:"password" form:"password"`
	Fullname        string `gorm:"not null" json:"fullname" form:"fullname"`
	Displayname     string `json:"displayname" form:"displayname"`
	Alumnus         string `json:"alumnus" form:"alumnus"`
	Workplace       string `json:"workplace" form:"workplace"`
	PracticeAddress string `json:"practice_address" form:"practice_address"`
	Price           float64
	Balance         float64
	// Users           []*User ` gorm:"many2many:user_doctors;"`
}
