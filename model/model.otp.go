package model

import (
	"gorm.io/gorm"
)

type UserOTP struct {
	gorm.Model
	Email         string `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Username      string `json:"username" form:"username" gorm:"type:varchar(255)unique;not null"`
	Fullname 	string `json:"fullname" form:"fullname" gorm:"type:varchar(255)"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Telp          string `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat        string `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender        string `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	BirthDate    string `json:"birthdate" form:"birthdate" gorm:"type:date"`
	OTP 		   string    `json:"otp" form:"otp"`
}

